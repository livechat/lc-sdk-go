package webhooks

import (
	"encoding/json"
	"net/http"
)

func ParseChatUserRemovedPayload(body []byte) (*ChatUserRemovedPayload, error) {
	var p ChatUserRemovedPayload
	return &p, json.Unmarshal(body, &p)
}

type ChatUserRemovedPayload struct {
	WebhookDetails
	Payload struct {
		ChatID   string `json:"chat_id"`
		UserType string `json:"user_type"`
		UserID   string `json:"user_id"`
	}
}

type ChatUserRemovedHandler func(*ChatUserRemovedPayload) error

func NewChatUserRemovedHandler(h ChatUserRemovedHandler) http.HandlerFunc {
	return webhookHandler(
		func(payload interface{}) error {
			return h(payload.(*ChatUserRemovedPayload))
		},
		func(body []byte) (interface{}, error) {
			return ParseChatUserRemovedPayload(body)
		},
	)
}
