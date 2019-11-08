package configuration

import (
	"errors"
	"net/http"
	"time"

	"github.com/livechat/lc-sdk-go/internal"
	"github.com/livechat/lc-sdk-go/objects/authorization"
)

const apiVersion = "v3.1"

type API struct {
	base internal.APIBase
}

func NewAPI(t authorization.TokenGetter, client *http.Client, clientID string) (*API, error) {
	if t == nil {
		return nil, errors.New("cannot initialize api without TokenGetter")
	}

	if client == nil {
		client = &http.Client{
			Timeout: 20 * time.Second,
		}
	}

	return &API{
		internal.APIBase{
			ApiVersion:  apiVersion,
			ApiName:     "configuration",
			TokenGetter: t,
			ApiURL:      "https://api.livechatinc.com",
			ClientID:    clientID,
			HttpClient:  client,
		},
	}, nil
}

func (a *API) RegisterWebhook(w *Webhook) error {
	var resp struct {
		WebhookID string `json:"webhook_id"`
	}
	err := a.base.Call("register_webhook", w, &resp)
	if err != nil {
		w.ID = resp.WebhookID
	}
	return err
}

func (a *API) WebhooksConfig() ([]Webhook, error) {
	req := map[string]string{}
	ws := make([]Webhook, 0)
	err := a.base.Call("get_webhooks_config", req, ws)

	return ws, err
}

func (a *API) UnregisterWebhook(id string) error {
	req := map[string]string{"webhook_id": id}
	return a.base.Call("unregister_webhook", req, nil)
}

func (a *API) ChangeURL(url string) {
	a.base.ApiURL = url
}
