package webhooks

import (
	"encoding/json"
	"time"
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

type User struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Present  bool   `json:"present"`
	LastSeen customer.Time   `json:"last_seen_timestamp"`
	userSpecific
}

type userSpecific struct {
	RoutingStatus json.RawMessage `json:"routing_status"`
	LastVisit  json.RawMessage `json:"last_visit"`
	Statistics json.RawMessage `json:"statistics"`
	AgentLastEventCreatedAt    json.RawMessage `json:"agent_last_event_created_at"`
	CustomerLastEventCreatedAt json.RawMessage `json:"customer_last_event_created_at"`
}

type Agent struct {
	RoutingStatus string `json:"routing_status"`
	*User
}

type Customer struct {
	LastVisit  customer.Visit `json:"last_visit"`
	Statistics struct {
		VisitsCount            int `json:"visits_count"`
		ThreadsCount           int `json:"threads_count"`
		ChatsCount             int `json:"chats_count"`
		PageViewsCount         int `json:"page_views_count"`
		GreetingsShownCount    int `json:"greetings_shown_count"`
		GreetingsAcceptedCount int `json:"greetings_accepted_count"`
	}
	AgentLastEventCreatedAt    time.Time `json:"agent_last_event_created_at"`
	CustomerLastEventCreatedAt time.Time `json:"customer_last_event_created_at"`
	*User
}

func (u *User) Agent() *Agent {
	if u.Type != "agent" {
		return nil
	}
	var a Agent

	a.User = u
	if err := json.Unmarshal(u.RoutingStatus, &a.RoutingStatus); err != nil {
		return nil
	}
	return &a
}

func (u *User) Customer() *Customer {
	if u.Type != "customer" {
		return nil
	}
	var c Customer

	c.User = u
	if err := json.Unmarshal(u.LastVisit, &c.LastVisit); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.Statistics, &c.Statistics); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.AgentLastEventCreatedAt, &c.AgentLastEventCreatedAt); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.CustomerLastEventCreatedAt, &c.CustomerLastEventCreatedAt); err != nil {
		return nil
	}
	return &c
}

type ChatUserAdded struct {
	ChatID   string  `json:"chat_id"`
	ThreadID string  `json:"thread_id"`
	User     User    `json:"user"`
	UserType string  `json:"user_type"` // TODO shall this be enum?
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
