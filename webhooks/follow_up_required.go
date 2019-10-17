package webhooks

import (
	"encoding/json"
	"net/http"
)

func ParseFollowUpRequiredPayload(body []byte) (*FollowUpRequiredPayload, error) {
	var p FollowUpRequiredPayload
	return &p, json.Unmarshal(body, &p)
}

type FollowUpRequiredPayload struct {
	WebhookDetails
	Payload struct {
		ChatID     string `json:"chat_id"`
		ThreadID   string `json:"thread_id"`
		CustomerID string `json:"customer_id"`
	} `json:"payload"`
}

func NewFollowUpRequiredPayload(h func(*FollowUpRequiredPayload) error) http.HandlerFunc {
	return webhookHandler(
		func(payload interface{}) error {
			return h(payload.(*FollowUpRequiredPayload))
		},
		func(body []byte) (interface{}, error) {
			return ParseFollowUpRequiredPayload(body)
		},
	)
}
