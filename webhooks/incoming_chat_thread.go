package webhooks

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/livechat/lc-sdk-go/customer"
)

type WebhookDetails struct {
	WebhookID string `json:"webhook_id"`
	SecretKey string `json:"secret_key"`
	Action    string `json:"action"`
}

func ParseIncomingChatThreadPayload(body []byte) (*IncomingChatThreadPayload, error) {
	var p IncomingChatThreadPayload
	return &p, json.Unmarshal(body, &p)
}

type IncomingChatThreadPayload struct {
	WebhookDetails
	Payload struct {
		Chat customer.Chat `json:"chat"`
	} `json:"payload"`
}

func (p *IncomingChatThreadPayload) UnmarshalJSON(data []byte) error {
	type PayloadAlias IncomingChatThreadPayload
	type SingleThread struct {
		Payload struct {
			Chat struct {
				Thread customer.Thread `json:"thread"`
			} `json:"chat"`
		} `json:"payload"`
	}
	var pa PayloadAlias
	if err := json.Unmarshal(data, &pa); err != nil {
		return err
	}
	*p = IncomingChatThreadPayload(pa)

	var st SingleThread
	if err := json.Unmarshal(data, &st); err != nil {
		return err
	}
	p.Payload.Chat.Threads = append(p.Payload.Chat.Threads, st.Payload.Chat.Thread)
	return nil
}

type IncomingChatThreadHandler func(*IncomingChatThreadPayload) error

func NewIncomingChatThreadHandler(h IncomingChatThreadHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		payload, err := ParseIncomingChatThreadPayload(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h(payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
