package configuration

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	api_errors "github.com/livechat/lc-sdk-go/errors"
)

const apiVersion = "v3.2"

type API struct {
	httpClient  *http.Client
	ApiURL      string
	clientID    string
	tokenGetter func() *Token
}

type Token struct {
	LicenseID   int
	AccessToken string
	Region      string
}

type TokenGetter func() *Token

func NewAPI(t TokenGetter, client *http.Client, clientID string) (*API, error) {
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
		ApiURL:      "https://api.livechatinc.com/",
		clientID:    clientID,
		httpClient:  client,
	}, nil
}

func (a *API) RegisterWebhook(wh *Webhook) (string, error) {
	var resp registerWebhookResponse
	err := a.call("register_webhook", wh, &resp)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (a *API) UnregisterWebhook(id string) error {
	req := &unregisterWebhookRequest{
		ID: id,
	}
	return a.call("unregister_webhook", req, nil)
}

func (a *API) GetWebhooksConfig() ([]*RegisteredWebhook, error) {
	resp := []*RegisteredWebhook{}
	err := a.call("get_webhooks_config", nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *API) call(action string, reqPayload interface{}, respPayload interface{}) error {
	rawBody, err := json.Marshal(reqPayload)
	if err != nil {
		return err
	}
	token := a.tokenGetter()
	if token == nil {
		return fmt.Errorf("couldn't get token")
	}

	url := fmt.Sprintf("%s/%s/configuration/action/%s", a.ApiURL, apiVersion, action)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(rawBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Set("User-agent", fmt.Sprintf("GO SDK Application %s", a.clientID))
	req.Header.Set("X-Region", token.Region)

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
