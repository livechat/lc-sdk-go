package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	api_errors "github.com/livechat/lc-sdk-go/errors"
	"github.com/livechat/lc-sdk-go/objects/authorization"
)

type APIBase struct {
	ApiVersion  string
	ApiName     string
	HttpClient  *http.Client
	ApiURL      string
	ClientID    string
	TokenGetter func() *authorization.Token
}

func (a *APIBase) Send(req *http.Request, respPayload interface{}) error {
	resp, err := a.HttpClient.Do(req)
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

func (a *APIBase) Call(action string, reqPayload interface{}, respPayload interface{}) error {
	rawBody, err := json.Marshal(reqPayload)
	if err != nil {
		return err
	}
	token := a.TokenGetter()
	if token == nil {
		return fmt.Errorf("couldn't get token")
	}

	url := fmt.Sprintf("%s/%s/%s/action/%s?license_id=%v", a.ApiURL, a.ApiVersion, a.ApiName, action, token.LicenseID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(rawBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Set("User-agent", fmt.Sprintf("GO SDK Application %s", a.ClientID))
	req.Header.Set("X-Region", token.Region)

	return a.Send(req, respPayload)
}
