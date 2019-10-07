// Package webhooks implements handlers and definitions of LiveChat webhooks.
//
// General LiveChat webhooks documentation is available here:
// https://developers.livechatinc.com/docs/management/configuration-api/#webhooks
package webhooks

import (
	"encoding/json"

	"github.com/livechat/lc-sdk-go/objects"
)

// WebhookBase represents general webhook format.
type WebhookBase struct {
	WebhookID      string          `json:"webhook_id"`
	SecretKey      string          `json:"secret_key"`
	Action         string          `json:"action"`
	Payload        json.RawMessage `json:"payload"`
	AdditionalData json.RawMessage `json:"additional_data"`
}

// IncomingChatThread represents payload of incoming_chat_thread webhook.
type IncomingChatThread struct {
	Chat objects.Chat `json:"chat"`
}

// ThreadClosed represents payload of thread_closed webhook.
type ThreadClosed struct {
	ChatID   string `json:"chat_id"`
	ThreadID string `json:"thread_id"`
	UserID   string `json:"user_id"`
}

// AccessGranted represents payload of access_granted webhook.
type AccessGranted struct {
	Resource string         `json:"resource"`
	ID       string         `json:"id"`
	Access   objects.Access `json:"access"`
}

// AccessRevoked represents payload of access_revoked webhook.
type AccessRevoked struct {
	Resource string         `json:"resource"`
	ID       string         `json:"id"`
	Access   objects.Access `json:"access"`
}

// AccessSet represents payload of access_set webhook.
type AccessSet struct {
	Resource string         `json:"resource"`
	ID       string         `json:"id"`
	Access   objects.Access `json:"access"`
}

// ChatUserAdded represents payload of chat_user_added webhook.
type ChatUserAdded struct {
	ChatID   string       `json:"chat_id"`
	ThreadID string       `json:"thread_id"`
	User     objects.User `json:"user"`
	UserType string       `json:"user_type"`
}

// ChatUserRemoved represents payload of chat_user_removed webhook.
type ChatUserRemoved struct {
	ChatID   string `json:"chat_id"`
	ThreadID string `json:"thread_id"`
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
}

// IncomingEvent represents payload of incoming_event webhook.
type IncomingEvent struct {
	ChatID   string        `json:"chat_id"`
	ThreadID string        `json:"thread_id"`
	Event    objects.Event `json:"event"`
}

// EventUpdated represents payload of event_updated webhook.
type EventUpdated struct {
	ChatID   string        `json:"chat_id"`
	ThreadID string        `json:"thread_id"`
	Event    objects.Event `json:"event"`
}

// IncomingRichMessagePostback represents payload of incoming_rich_message_postback webhook.
type IncomingRichMessagePostback struct {
	UserID   string `json:"user_id"`
	ChatID   string `json:"chat_id"`
	ThreadID string `json:"thread_id"`
	EventID  string `json:"event_id"`
	Postback struct {
		ID      string `json:"id"`
		Toggled bool   `json:"toggled"`
	} `json:"postback"`
}

// ChatPropertiesUpdated represents payload of chat_properties_updated webhook.
type ChatPropertiesUpdated struct {
	ChatID     string             `json:"chat_id"`
	Properties objects.Properties `json:"properties"`
}

// ChatPropertiesDeleted represents payload of chat_properties_deleted webhook.
type ChatPropertiesDeleted struct {
	ChatID     string              `json:"chat_id"`
	Properties map[string][]string `json:"properties"`
}

// ChatThreadPropertiesDeleted represents payload of chat_thread_properties_deleted webhook.
type ChatThreadPropertiesDeleted struct {
	ChatID     string              `json:"chat_id"`
	ThreadID   string              `json:"thread_id"`
	Properties map[string][]string `json:"properties"`
}

// ChatThreadPropertiesUpdated represents payload of chat_thread_properties_updated webhook.
type ChatThreadPropertiesUpdated struct {
	ChatID     string             `json:"chat_id"`
	ThreadID   string             `json:"thread_id"`
	Properties objects.Properties `json:"properties"`
}

// EventPropertiesUpdated represents payload of event_properties_updated webhook.
type EventPropertiesUpdated struct {
	ChatID     string             `json:"chat_id"`
	ThreadID   string             `json:"thread_id"`
	EventID    string             `json:"event_id"`
	Properties objects.Properties `json:"properties"`
}

// EventPropertiesDeleted represents payload of event_properties_deleted webhook.
type EventPropertiesDeleted struct {
	ChatID     string              `json:"chat_id"`
	ThreadID   string              `json:"thread_id"`
	EventID    string              `json:"event_id"`
	Properties map[string][]string `json:"properties"`
}

// ChatThreadTagged represents payload of chat_thread_tagged webhook.
type ChatThreadTagged struct {
	ChatID   string `json:"chat_id"`
	ThreadID string `json:"thread_id"`
	Tag      string `json:"tag"`
}

// ChatThreadUntagged represents payload of chat_thread_untagged webhook.
type ChatThreadUntagged struct {
	ChatID   string `json:"chat_id"`
	ThreadID string `json:"thread_id"`
	Tag      string `json:"tag"`
}

// AgentStatusChanged represents payload of agent_status_changed webhook.
type AgentStatusChanged struct {
	AgentID string `json:"agent_id"`
	Status  string `json:"status"`
}

// AgentDeleted represents payload of agent_deleted webhook.
type AgentDeleted struct {
	AgentID string `json:"agent_id"`
}

// CustomerCreated represents payload of customer_created webhook.
type CustomerCreated objects.Customer

// EventsMarkedAsSeen represents payload of events_marked_as_seen webhook.
type EventsMarkedAsSeen struct {
	UserID   string `json:"user_id"`
	ChatID   string `json:"chat_id"`
	SeenUpTo string `json:"seen_up_to"`
}

// UnmarshalJSON implements json.Unmarshaler interface for IncomingChatThread.
func (p *IncomingChatThread) UnmarshalJSON(data []byte) error {
	type PayloadAlias IncomingChatThread
	type SingleThread struct {
		Chat struct {
			Thread objects.Thread `json:"thread"`
		} `json:"chat"`
	}
	var pa PayloadAlias
	if err := json.Unmarshal(data, &pa); err != nil {
		return err
	}
	*p = IncomingChatThread(pa)

	var st SingleThread
	if err := json.Unmarshal(data, &st); err != nil {
		return err
	}
	p.Chat.Threads = append(p.Chat.Threads, st.Chat.Thread)
	return nil
}
