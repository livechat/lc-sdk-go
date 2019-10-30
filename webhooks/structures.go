package webhooks

import (
	"encoding/json"

	"github.com/livechat/lc-sdk-go/customer"
	"github.com/livechat/lc-sdk-go/objects/events"
)

type WebhookBase struct {
	WebhookID      string          `json:"webhook_id"`
	SecretKey      string          `json:"secret_key"`
	Action         string          `json:"action"`
	LicenseID      int             `json:"license_id"`
	Payload        json.RawMessage `json:"payload"`
	AdditionalData json.RawMessage `json:"additional_data"`
}

type IncomingChatThread struct {
	Chat customer.Chat `json:"chat"`
}

type ThreadClosed struct {
	ChatID   string `json:"chat_id"`
	ThreadID string `json:"thread_id"`
	UserID   string `json:"user_id"`
}

type AccessSet struct {
	Resource string          `json:"resource"`
	ID       string          `json:"id"`
	Access   customer.Access `json:"access"`
}

type ChatUserAdded struct {
	ChatID   string        `json:"chat_id"`
	ThreadID string        `json:"thread_id"`
	User     customer.User `json:"user"`
	UserType string        `json:"user_type"` // TODO shall this be enum?
}

type ChatUserRemoved struct {
	ChatID   string `json:"chat_id"`
	ThreadID string `json:"thread_id"`
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"` // TODO shall this be enum?
}

type IncomingEvent struct {
	ChatID   string       `json:"chat_id"`
	ThreadID string       `json:"thread_id"`
	Event    events.Event `json:"event"`
}

type EventUpdated struct {
	ChatID   string       `json:"chat_id"`
	ThreadID string       `json:"thread_id"`
	Event    events.Event `json:"event"`
}

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

type ChatPropertiesUpdated struct {
	ChatID     string              `json:"chat_id"`
	Properties customer.Properties `json:"properties"`
}

type ChatPropertiesDeleted struct {
	ChatID     string              `json:"chat_id"`
	Properties map[string][]string `json:"properties"`
}

type ChatThreadPropertiesUpdated struct {
	ChatID     string              `json:"chat_id"`
	ThreadID   string              `json:"thread_id"`
	Properties customer.Properties `json:"properties"`
}

type ChatThreadPropertiesDeleted struct {
	ChatID     string              `json:"chat_id"`
	ThreadID   string              `json:"thread_id"`
	Properties map[string][]string `json:"properties"`
}

type EventPropertiesUpdated struct {
	ChatID     string              `json:"chat_id"`
	ThreadID   string              `json:"thread_id"`
	EventID    string              `json:"event_id"`
	Properties customer.Properties `json:"properties"`
}

type EventPropertiesDeleted struct {
	ChatID     string              `json:"chat_id"`
	ThreadID   string              `json:"thread_id"`
	EventID    string              `json:"event_id"`
	Properties map[string][]string `json:"properties"`
}

type FollowUpRequested struct {
	ChatID     string `json:"chat_id"`
	ThreadID   string `json:"thread_id"`
	CustomerID string `json:"customer_id"`
}

type ChatThreadTagged struct {
	ChatID   string `json:"chat_id"`
	ThreadID string `json:"thread_id"`
	Tag      string `json:"tag"`
}

type ChatThreadUntagged struct {
	ChatID   string `json:"chat_id"`
	ThreadID string `json:"thread_id"`
	Tag      string `json:"tag"`
}

type AgentStatusChanged struct {
	AgentID string `json:"agent_id"`
	Status  string `json:"status"` // TODO - should this be enum?
}

type AgentDeleted struct {
	AgentID string `json:"agent_id"`
}

type EventsMarkedAsSeen struct {
	UserID   string `json:"user_id"`
	ChatID   string `json:"chat_id"`
	SeenUpTo string `json:"seen_up_to"` // TODO - should we parse this into time type?
}

func (p *IncomingChatThread) UnmarshalJSON(data []byte) error {
	type PayloadAlias IncomingChatThread
	type SingleThread struct {
		Chat struct {
			Thread customer.Thread `json:"thread"`
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
