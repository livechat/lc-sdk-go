package configuration

import (
	"errors"
	"net/http"
	"time"

	"github.com/livechat/lc-sdk-go/internal"

	"github.com/livechat/lc-sdk-go/configuration/action"

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

type Webhook struct {
	ID             string          `json:"webhook_id"`
	Action         action.Webhook  `json:"action"`
	SecretKey      string          `json:"secret_key"`
	URL            string          `json:"url"`
	AdditionalData []string        `json:"additional_data,omitempty"`
	Description    string          `json:"description,omitempty"`
	Filters        *WebhookFilters `json:"filters,omitempty"`
}

type WebhookFilters struct {
	AuthorType    string `json:"author_type"`
	OnlyMyChats   bool   `json:"only_my_chats"`
	ChatMemberIds struct {
		AgentsAny     []string `json:"agents_any,omitempty"`
		AgentsExclude []string `json:"agents_exclude,omitempty"`
	} `json:"chat_member_ids"`
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
