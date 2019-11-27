package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	api_errors "github.com/livechat/lc-sdk-go/errors"
	"github.com/livechat/lc-sdk-go/objects/authorization"
)

const apiVersion = "3.1"

// API provides the base client for making requests to Livechat Web APIs.
type API struct {
	httpClient *http.Client
	// APIURL defines base of API endpoint (without version, component and action).
	APIURL      string
	clientID    string
	version     string
	name        string
	tokenGetter func() *authorization.Token
}

// NewAPI returns ready to use raw API client. This is a base that is used internally
// by specialized clients for each API, you should use those instead
//
// If provided client is nil, then default http client with 20s timeout is used.
func NewAPI(t authorization.TokenGetter, client *http.Client, clientID, name string) (*API, error) {
	if t == nil {
		return nil, errors.New("cannot initialize api without TokenGetter")
	}

	if client == nil {
		client = &http.Client{
			Timeout: 20 * time.Second,
		}
	}

	return &API{
		tokenGetter: t,
		APIURL:      "https://api.livechatinc.com",
		clientID:    clientID,
		version:     apiVersion,
		name:        name,
		httpClient:  client,
	}, nil
}

// Call sends request to API with given action
func (a *API) Call(action string, reqPayload interface{}, respPayload interface{}) error {
	rawBody, err := json.Marshal(reqPayload)
	if err != nil {
		return err
	}
	token := a.tokenGetter()
	if token == nil {
		return fmt.Errorf("couldn't get token")
	}

	url := fmt.Sprintf("%s/v%s/%s/action/%s?license_id=%v", a.APIURL, a.version, a.name, action, token.LicenseID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(rawBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Set("User-agent", fmt.Sprintf("GO SDK Application %s", a.clientID))
	req.Header.Set("X-Region", token.Region)

	return a.send(req, respPayload)
}

// UploadFile uploads a file to LiveChat CDN.
// Returned URL shall be used in call to SendFile or SendEvent or it'll become invalid
// in about 24 hours.
// UploadFile is supported only in Agent and Customer API
func (a *API) UploadFile(filename string, file []byte) (string, error) {
	switch a.name {
	case "customer":
	case "agent":
	default:
		return "", fmt.Errorf("UploadFile method is not supported in %v API", a.name)
	}

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
	token := a.tokenGetter()
	if token == nil {
		return "", fmt.Errorf("couldn't get token")
	}
	url := fmt.Sprintf("%s/v%s/%s/action/upload_file?license_id=%v", a.APIURL, a.version, a.name, token.LicenseID)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("couldn't create new POST request to '%v': %v", url, err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Set("User-agent", fmt.Sprintf("GO SDK Application %s", a.clientID))
	req.Header.Set("X-Region", token.Region)

	var resp struct {
		URL string `json:"url"`
	}
	err = a.send(req, &resp)
	return resp.URL, err
}

func (a *API) send(req *http.Request, respPayload interface{}) error {
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		apiErr := &api_errors.ErrAPI{}
		if err := json.Unmarshal(bodyBytes, apiErr); err != nil {
			return fmt.Errorf("couldn't unmarshal error response: %s (code: %d, raw body: %s)", err.Error(), resp.StatusCode, string(bodyBytes))
		}
		if apiErr.Error() == "" {
			return fmt.Errorf("couldn't unmarshal error response (code: %d, raw body: %s)", resp.StatusCode, string(bodyBytes))
		}
		return apiErr
	}

	if err != nil {
		return err
	}

	return json.Unmarshal(bodyBytes, respPayload)
}
