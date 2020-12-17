package main

import (
	"context"
	"errors"

	"github.com/livechat/lc-sdk-go/v3/agent"
	"github.com/livechat/lc-sdk-go/v3/authorization"
	"github.com/livechat/lc-sdk-go/v3/objects"
	"github.com/livechat/lc-sdk-go/v3/webhooks"
)

type IncomingEventHandler struct {
	cfg *Configuration
	tr  tokensRepository
}

func NewIncomingEventHandler(cfg *Configuration, tr tokensRepository) *IncomingEventHandler {
	return &IncomingEventHandler{cfg, tr}
}

func (h *IncomingEventHandler) Handle(ctx context.Context, wh *webhooks.Webhook) error {
	payload, ok := wh.Payload.(*webhooks.IncomingEvent)
	if !ok {
		return errors.New("type assertion failed")
	}
	if payload.Event.Type != "message" {
		return nil
	}

	t := h.tr.Get(wh.WebhookID)
	if t == nil {
		return errors.New("retrieving webhook token failed")
	}

	tg := func() *authorization.Token {
		return &authorization.Token{
			AccessToken: t.AccessToken,
			Region:      t.Region,
		}
	}
	api, err := agent.NewAPI(tg, nil, t.ClientID)
	if err != nil {
		return errors.New("agent-api initilization failed")
	}

	msg := &objects.Message{
		Event: objects.Event{
			Type:       "message",
			Recipients: "all",
		},
		Text: "You said: " + payload.Event.Message().Text,
	}
	api.SendEvent(payload.ChatID, msg, true)

	return nil
}
