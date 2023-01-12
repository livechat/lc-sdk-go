package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/livechat/lc-sdk-go/v5/authorization"
	api_errors "github.com/livechat/lc-sdk-go/v5/errors"
	"github.com/livechat/lc-sdk-go/v5/metrics"
)

const apiVersion = "3.5"

// RetryStrategyFunc is called by each API method if set to retry when handling an error.
// If not set, there will be no retry at all.
//
// It accepts two arguments: attempts - number of sent requests (starting from 0)
// and err - error as ErrAPI struct (with StatusCode and Details)
// It returns info whether to retry the request.
type RetryStrategyFunc func(attempts uint, err error) bool

// StatsSinkFunc is called after each API method with statistics of that method execution.
type StatsSinkFunc func(callStats metrics.APICallStats)

type api struct {
	httpClient            *http.Client
	clientID              string
	tokenGetter           authorization.TokenGetter
	httpEndpointGenerator HTTPEndpointGenerator
	host                  string
	customHeaders         http.Header
	retryStrategy         RetryStrategyFunc
	statsSink             StatsSinkFunc
	logger                *log.Logger
}

// HTTPRequestGenerator is called by each API method to generate api http url.
type HTTPEndpointGenerator func(*authorization.Token, string, string) string

type CallOptions struct {
	Method string
}

// NewAPI returns ready to use raw API client. This is a base that is used internally
// by specialized clients for each API, you should use those instead
//
// If provided client is nil, then default http client with 20s timeout is used.
func NewAPI(t authorization.TokenGetter, client *http.Client, clientID string, r HTTPEndpointGenerator) (*api, error) {
	if t == nil {
		return nil, errors.New("cannot initialize api without TokenGetter")
	}

	if client == nil {
		client = &http.Client{
			Timeout: 20 * time.Second,
		}
	}

	return &api{
		tokenGetter:           t,
		clientID:              clientID,
		httpClient:            client,
		host:                  "https://api.livechatinc.com",
		httpEndpointGenerator: r,
		customHeaders:         make(http.Header),
		statsSink:             func(metrics.APICallStats) {},
		logger:                log.Default(),
	}, nil
}

// Call sends request to API with given action
func (a *api) Call(action string, reqPayload interface{}, respPayload interface{}, opts ...*CallOptions) error {
	token, err := a.getToken()
	if err != nil {
		return err
	}
	start := time.Now()

	endpoint := a.httpEndpointGenerator(token, a.host, action)
	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		return fmt.Errorf("couldn't create new http request: %v", err)
	}

	if len(opts) > 0 && opts[0].Method == http.MethodGet {
		req.Method = opts[0].Method
		qs, err := query.Values(reqPayload)
		if err != nil {
			return err
		}

		encodedQuery := qs.Encode()
		if req.URL.RawQuery != "" && encodedQuery != "" {
			encodedQuery = "&" + encodedQuery
		}
		req.URL.RawQuery = req.URL.RawQuery + encodedQuery
	} else {
		rawBody, err := json.Marshal(reqPayload)
		if err != nil {
			return err
		}
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(rawBody)), nil
		}
		req.Body, _ = req.GetBody()
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.Type, token.AccessToken))
	req.Header.Set("User-agent", fmt.Sprintf("GO SDK Application %s", a.clientID))
	req.Header.Set("X-Region", token.Region)

	for key, val := range a.customHeaders {
		if len(val) == 0 {
			continue
		}
		req.Header.Set(key, val[0])
	}
	err = a.send(req, respPayload)

	executionTime := time.Since(start)
	a.statsSink(metrics.APICallStats{Method: action, ExecutionTime: executionTime, Success: err == nil})

	return err
}

// SetCustomHeader allows to set a custom header (e.g. X-Debug-Id or X-Author-Id) that will be sent in every request
func (a *api) SetCustomHeader(key, val string) {
	a.customHeaders.Set(key, val)
}

// SetRetryStrategy allows to set a retry strategy that will be performed in every failed request
func (a *api) SetRetryStrategy(f RetryStrategyFunc) {
	a.retryStrategy = f
}

// SetStatsSink allows to set a statistics sink that will send API calls metrics data to SDK consumers
func (a *api) SetStatsSink(f StatsSinkFunc) {
	a.statsSink = f
}

type fileUploadAPI struct{ *api }

// NewAPIWithFileUpload returns ready to use raw API client with file upload functionality.
func NewAPIWithFileUpload(t authorization.TokenGetter, client *http.Client, clientID string, r HTTPEndpointGenerator) (*fileUploadAPI, error) {
	api, err := NewAPI(t, client, clientID, r)
	if err != nil {
		return nil, err
	}
	return &fileUploadAPI{api}, nil
}

// UploadFile uploads a file to LiveChat CDN.
// Returned URL shall be used in call to SendFile or SendEvent or it'll become invalid
// in about 24 hours.
func (a *fileUploadAPI) UploadFile(filename string, file []byte) (string, error) {
	token := a.tokenGetter()
	if token == nil {
		return "", errors.New("couldn't get token")
	}
	start := time.Now()

	endpoint := a.httpEndpointGenerator(token, a.host, "upload_file")
	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("couldn't create new http request: %v", err)
	}
	req.Method = "POST"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	w, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("couldn't create form file: %v", err)
	}
	if _, err := w.Write(file); err != nil {
		return "", fmt.Errorf("couldn't write file to multipart writer: %v", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("couldn't close multipart writer: %v", err)
	}

	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(body.Bytes())), nil
	}
	req.Body, _ = req.GetBody()

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.Type, token.AccessToken))
	req.Header.Set("User-agent", fmt.Sprintf("GO SDK Application %s", a.clientID))
	req.Header.Set("X-Region", token.Region)

	var resp struct {
		URL string `json:"url"`
	}
	err = a.send(req, &resp)

	executionTime := time.Since(start)
	a.statsSink(metrics.APICallStats{Method: "upload_file", ExecutionTime: executionTime, Success: err == nil})

	return resp.URL, err
}

func (a *api) send(req *http.Request, respPayload interface{}) error {
	var attempts uint
	var do func() error

	do = func() error {
		resp, err := a.httpClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			apiErr := &api_errors.ErrAPI{}
			if err := json.Unmarshal(bodyBytes, apiErr); err != nil {
				return fmt.Errorf("couldn't unmarshal error response: %s (code: %d, raw body: %s)", err.Error(), resp.StatusCode, string(bodyBytes))
			}
			if apiErr.Error() == "" {
				return fmt.Errorf("couldn't unmarshal error response (code: %d, raw body: %s)", resp.StatusCode, string(bodyBytes))
			}

			if a.retryStrategy == nil || !a.retryStrategy(attempts, apiErr) {
				return apiErr
			}

			token, err := a.getToken()
			if err != nil {
				return err
			}

			req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.Type, token.AccessToken))
			if req.Body != nil {
				reqBody, err := req.GetBody()
				if err != nil {
					return fmt.Errorf("couldn't get request body: %v", err)
				}
				req.Body = reqBody
			}

			attempts++
			return do()
		}

		if err != nil {
			return err
		}

		if h := resp.Header.Get("Legacy"); h != "" {
			a.logger.Printf("[Notice] This is a legacy version. It will be deprecated after %s.", h)
		}
		if h := resp.Header.Get("Deprecation"); h != "" {
			a.logger.Printf("[Warning] This version is deprecated. It will be decommissioned after %s.", h)
		}

		return json.Unmarshal(bodyBytes, respPayload)
	}

	return do()
}

func (a *api) getToken() (*authorization.Token, error) {
	token := a.tokenGetter()
	if token == nil {
		return nil, errors.New("couldn't get token")
	}
	if token.Type != authorization.BearerToken && token.Type != authorization.BasicToken {
		return nil, errors.New("unsupported token type")
	}

	return token, nil
}

// SetCustomHost allows to change API host address. This method is mostly for LiveChat internal testing and should not be used in production environments.
func (a *api) SetCustomHost(host string) {
	a.host = host
}

// DefaultHTTPRequestGenerator generates API request for given service in stable version.
func DefaultHTTPRequestGenerator(name string) HTTPEndpointGenerator {
	return func(token *authorization.Token, host, action string) string {
		return fmt.Sprintf("%s/v%s/%s/action/%s", host, apiVersion, name, action)
	}
}

// SetLogger allows to set a custom logger. When it's not called, a default logger will be used.
func (a *api) SetLogger(logger *log.Logger) {
	a.logger = logger
}
