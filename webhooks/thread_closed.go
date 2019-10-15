package webhooks

import (
	"encoding/json"
	"net/http"
)

func ParseThreadClosedPayload(body []byte) (*ThreadClosedPayload, error) {
	var p ThreadClosedPayload
	return &p, json.Unmarshal(body, &p)
}

type ThreadClosedPayload struct {
	WebhookDetails
	Payload struct {
		ChatID   string `json:"chat_id"`
		ThreadID string `json:"thread_id"`
		UserID   string `json:"user_id"`
	}
}

type ThreadClosedHandler func(*ThreadClosedPayload) error

func NewThreadClosedHandler(h ThreadClosedHandler) http.HandlerFunc {
	return webhookHandler(
		func(payload interface{}) error {
			return h(payload.(*ThreadClosedPayload))
		},
		func(body []byte) (interface{}, error) {
			return ParseThreadClosedPayload(body)
		},
	)
}
