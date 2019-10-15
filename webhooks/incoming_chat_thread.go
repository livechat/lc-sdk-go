package webhooks

import (
	"encoding/json"
	"net/http"

	"github.com/livechat/lc-sdk-go/customer"
)

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
	return webhookHandler(
		func(payload interface{}) error {
			return h(payload.(*IncomingChatThreadPayload))
		},
		func(body []byte) (interface{}, error) {
			return ParseIncomingChatThreadPayload(body)
		},
	)
}
