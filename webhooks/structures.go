// Package webhooks implements handlers and definitions of LiveChat webhooks.
//
// General LiveChat webhooks documentation is available here:
// https://developers.livechatinc.com/docs/management/configuration-api/#webhooks
package webhooks

import (
	"encoding/json"

	"github.com/livechat/lc-sdk-go/objects"
)

// Webhook represents general webhook format.
type Webhook struct {
	WebhookID      string          `json:"webhook_id"`
	SecretKey      string          `json:"secret_key"`
	Action         string          `json:"action"`
	LicenseID      int             `json:"license_id"`
	AdditionalData json.RawMessage `json:"additional_data"`
	RawPayload     json.RawMessage `json:"payload"`
	Payload        interface{}
}

// IncomingChat represents payload of incoming_chat webhook.
type IncomingChat struct {
	Chat objects.Chat `json:"chat"`
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

// ChatDeactivated represents payload of chat_deactivated webhook.
type ChatDeactivated struct {
	ChatID   string `json:"chat_id"`
	ThreadID string `json:"thread_id"`
	UserID   string `json:"user_id"`
}

// ChatPropertiesUpdated represents payload of chat_properties_updated webhook.
type ChatPropertiesUpdated struct {
	ChatID     string             `json:"chat_id"`
	Properties objects.Properties `json:"properties"`
}

// ThreadPropertiesUpdated represents payload of thread_properties_updated webhook.
type ThreadPropertiesUpdated struct {
	ChatID     string             `json:"chat_id"`
	ThreadID   string             `json:"thread_id"`
	Properties objects.Properties `json:"properties"`
}

// ChatPropertiesDeleted represents payload of chat_properties_deleted webhook.
type ChatPropertiesDeleted struct {
	ChatID     string              `json:"chat_id"`
	Properties map[string][]string `json:"properties"`
}

// ThreadPropertiesDeleted represents payload of thread_properties_deleted webhook.
type ThreadPropertiesDeleted struct {
	ChatID     string              `json:"chat_id"`
	ThreadID   string              `json:"thread_id"`
	Properties map[string][]string `json:"properties"`
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

// ThreadTagged represents payload of thread_tagged webhook.
type ThreadTagged struct {
	ChatID   string `json:"chat_id"`
	ThreadID string `json:"thread_id"`
	Tag      string `json:"tag"`
}

// ThreadUntagged represents payload of thread_untagged webhook.
type ThreadUntagged struct {
	ChatID   string `json:"chat_id"`
	ThreadID string `json:"thread_id"`
	Tag      string `json:"tag"`
}

// AgentDeleted represents payload of agent_deleted webhook.
type AgentDeleted struct {
	AgentID string `json:"agent_id"`
}

// EventsMarkedAsSeen represents payload of events_marked_as_seen webhook.
type EventsMarkedAsSeen struct {
	UserID   string `json:"user_id"`
	ChatID   string `json:"chat_id"`
	SeenUpTo string `json:"seen_up_to"`
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

// CustomerCreated represents payload of customer_created webhook.
type CustomerCreated objects.Customer

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

// RoutingStatusSet represents payload of routing_status_set webhook.
type RoutingStatusSet struct {
	AgentID string `json:"agent_id"`
	Status  string `json:"status"`
}

// UnmarshalJSON implements json.Unmarshaler interface for IncomingChat.
func (p *IncomingChat) UnmarshalJSON(data []byte) error {
	type PayloadAlias IncomingChat
	type SingleThread struct {
		Chat struct {
			Thread objects.Thread `json:"thread"`
		} `json:"chat"`
	}
	var pa PayloadAlias
	if err := json.Unmarshal(data, &pa); err != nil {
		return err
	}
	*p = IncomingChat(pa)

	var st SingleThread
	if err := json.Unmarshal(data, &st); err != nil {
		return err
	}
	p.Chat.Threads = append(p.Chat.Threads, st.Chat.Thread)
	return nil
}
