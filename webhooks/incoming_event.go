package webhooks

import (
	"encoding/json"
	"net/http"

	"github.com/livechat/lc-sdk-go/objects/events"
)

func ParseIncomingEventPayload(body []byte) (*IncomingEventPayload, error) {
	var p IncomingEventPayload
	return &p, json.Unmarshal(body, &p)
}

type IncomingEventPayload struct {
	WebhookDetails
	Payload struct {
		Event    events.Event `json:"event"`
		ChatID   string       `json:"chat_id"`
		ThreadID string       `json:"thread_id"`
	} `json:"payload"`
}

func NewIncomingEventHandler(h func(*IncomingEventPayload) error) http.HandlerFunc {
	return webhookHandler(
		func(payload interface{}) error {
			return h(payload.(*IncomingEventPayload))
		},
		func(body []byte) (interface{}, error) {
			return ParseIncomingEventPayload(body)
		},
	)
}
