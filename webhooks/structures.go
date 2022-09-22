package webhooks

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/livechat/lc-sdk-go/v5/configuration"
	"github.com/livechat/lc-sdk-go/v5/internal"
)

// Properties represents LiveChat properties in form of property_namespace -> property -> value.
type Properties map[string]map[string]interface{}

// Webhook represents general webhook format.
type Webhook struct {
	WebhookID      string          `json:"webhook_id"`
	SecretKey      string          `json:"secret_key"`
	Action         string          `json:"action"`
	OrganizationID string          `json:"organization_id"`
	AdditionalData json.RawMessage `json:"additional_data"`
	RawPayload     json.RawMessage `json:"payload"`
	Payload        interface{}
}

// IncomingChat represents payload of incoming_chat webhook.
type IncomingChat struct {
	Chat Chat `json:"chat"`
}

// IncomingEvent represents payload of incoming_event webhook.
type IncomingEvent struct {
	ChatID   string `json:"chat_id"`
	ThreadID string `json:"thread_id"`
	Event    Event  `json:"event"`
}

// EventUpdated represents payload of event_updated webhook.
type EventUpdated struct {
	ChatID   string `json:"chat_id"`
	ThreadID string `json:"thread_id"`
	Event    Event  `json:"event"`
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
	ChatID     string     `json:"chat_id"`
	Properties Properties `json:"properties"`
}

// ThreadPropertiesUpdated represents payload of thread_properties_updated webhook.
type ThreadPropertiesUpdated struct {
	ChatID     string     `json:"chat_id"`
	ThreadID   string     `json:"thread_id"`
	Properties Properties `json:"properties"`
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

// UserAddedToChat represents payload of user_added_to_chat webhook.
type UserAddedToChat struct {
	ChatID      string `json:"chat_id"`
	ThreadID    string `json:"thread_id"`
	User        User   `json:"user"`
	UserType    string `json:"user_type"`
	Reason      string `json:"reason"`
	RequesterID string `json:"requester_id"`
}

// UserRemovedFromChat represents payload of user_removed_from_chat webhook.
type UserRemovedFromChat struct {
	ChatID      string `json:"chat_id"`
	ThreadID    string `json:"thread_id"`
	UserID      string `json:"user_id"`
	UserType    string `json:"user_type"`
	Reason      string `json:"reason"`
	RequesterID string `json:"requester_id"`
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

// AgentCreated represents payload of agent_created webhook.
type AgentCreated = configuration.Agent

// AgentUpdated represents payload of agent_updated webhook.
type AgentUpdated = configuration.Agent

// AgentDeleted represents payload of agent_deleted webhook.
type AgentDeleted struct {
	ID string `json:"id"`
}

// AgentSuspended represents payload of agent_suspended webhook.
type AgentSuspended struct {
	ID string `json:"id"`
}

// AgentUnsuspended represents payload of agent_unsuspended webhook.
type AgentUnsuspended struct {
	ID string `json:"id"`
}

// AgentApproved represents payload of agent_approved webhook.
type AgentApproved struct {
	ID string `json:"id"`
}

// EventsMarkedAsSeen represents payload of events_marked_as_seen webhook.
type EventsMarkedAsSeen struct {
	UserID   string `json:"user_id"`
	ChatID   string `json:"chat_id"`
	SeenUpTo string `json:"seen_up_to"`
}

// ChatAccessUpdated represents payload of chat_access_updated webhook.
type ChatAccessUpdated struct {
	ID     string `json:"id"`
	Access Access `json:"access"`
}

// IncomingCustomer represents payload of incoming_customer webhook.
type IncomingCustomer Customer

// EventPropertiesUpdated represents payload of event_properties_updated webhook.
type EventPropertiesUpdated struct {
	ChatID     string     `json:"chat_id"`
	ThreadID   string     `json:"thread_id"`
	EventID    string     `json:"event_id"`
	Properties Properties `json:"properties"`
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

// ChatTransferred represents payload of chat_transferred webhook.
type ChatTransferred struct {
	ChatID        string `json:"chat_id"`
	ThreadID      string `json:"thread_id,omitempty"`
	RequesterID   string `json:"requester_id,omitempty"`
	Reason        string `json:"reason"`
	TransferredTo struct {
		AgentIDs []string `json:"agent_ids,omitempty"`
		GroupIDs []int    `json:"group_ids,omitempty"`
	} `json:"transferred_to"`
	Queue struct {
		Position int    `json:"position"`
		WaitTime int    `json:"wait_time"`
		QueuedAt string `json:"queued_at"`
	} `json:"queue,omitempty"`
}

// CustomerSessionFieldsUpdated represents payload of customer_session_fields_updated webhook.
type CustomerSessionFieldsUpdated struct {
	ID         string `json:"id"`
	ActiveChat struct {
		ChatID   string `json:"chat_id"`
		ThreadID string `json:"thread_id"`
	} `json:"active_chat"`
	SessionFields []map[string]string `json:"session_fields"`
}

// GroupCreated represents payload of group_created webhook.
type GroupCreated struct {
	ID              int               `json:"id"`
	Name            string            `json:"name"`
	LanguageCode    string            `json:"language_code"`
	AgentPriorities map[string]string `json:"agent_priorities"`
}

// GroupUpdated represents payload of group_updated webhook.
type GroupUpdated struct {
	ID              int               `json:"id"`
	Name            string            `json:"name,omitempty"`
	LanguageCode    string            `json:"language_code,omitempty"`
	AgentPriorities map[string]string `json:"agent_priorities"`
}

// GroupDeleted represents payload of group_deleted webhook.
type GroupDeleted struct {
	ID int `json:"id"`
}

// AutoAccessAdded represents payload of auto_access_added webhook.
type AutoAccessAdded = configuration.AutoAccess

// AutoAccessUpdated represents payload of auto_access_updated webhook.
type AutoAccessUpdated = configuration.AutoAccess

// AutoAccessDeleted represents payload of auto_access_deleted webhook.
type AutoAccessDeleted struct {
	ID string `json:"id"`
}

// BotCreated represents payload of bot_created webhook.
type BotCreated struct {
	ID                   string                       `json:"id"`
	Name                 string                       `json:"name"`
	Avatar               string                       `json:"avatar,omitempty"`
	MaxChatsCount        *uint                        `json:"max_chats_count,omitempty"`
	DefaultGroupPriority configuration.GroupPriority  `json:"default_group_priority,omitempty"`
	Groups               []*configuration.GroupConfig `json:"groups,omitempty"`
	WorkScheduler        configuration.WorkScheduler  `json:"work_scheduler,omitempty"`
	Timezone             string                       `json:"timezone,omitempty"`
	OwnerClientID        string                       `json:"owner_client_id"`
	JobTitle             string                       `json:"job_title,omitempty"`
}

// BotUpdated represents payload of bot_updated webhook.
type BotUpdated struct {
	ID                   string                       `json:"id"`
	Name                 string                       `json:"name,omitempty"`
	Avatar               string                       `json:"avatar,omitempty"`
	MaxChatsCount        *uint                        `json:"max_chats_count,omitempty"`
	DefaultGroupPriority configuration.GroupPriority  `json:"default_group_priority,omitempty"`
	Groups               []*configuration.GroupConfig `json:"groups,omitempty"`
	WorkScheduler        configuration.WorkScheduler  `json:"work_scheduler,omitempty"`
	Timezone             string                       `json:"timezone,omitempty"`
	JobTitle             string                       `json:"job_title,omitempty"`
}

// BotDeleted represents payload of bot_deleted webhook.
type BotDeleted struct {
	ID string `json:"id"`
}

// UnmarshalJSON implements json.Unmarshaler interface for IncomingChat.
func (p *IncomingChat) UnmarshalJSON(data []byte) error {
	type PayloadAlias IncomingChat
	type SingleThread struct {
		Chat struct {
			Thread Thread `json:"thread"`
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

// User represents base of both Customer and Agent
//
// To get specific user type's structure, call Agent() or Customer() (based on Type value).
type User struct {
	ID             string    `json:"id"`
	Type           string    `json:"type"`
	Name           string    `json:"name"`
	Avatar         string    `json:"avatar"`
	Email          string    `json:"email"`
	Present        bool      `json:"present"`
	EventsSeenUpTo time.Time `json:"events_seen_up_to"`
	userSpecific
}

type userSpecific struct {
	RoutingStatus              json.RawMessage `json:"routing_status"`
	LastVisit                  json.RawMessage `json:"last_visit"`
	Statistics                 json.RawMessage `json:"statistics"`
	AgentLastEventCreatedAt    json.RawMessage `json:"agent_last_event_created_at"`
	CustomerLastEventCreatedAt json.RawMessage `json:"customer_last_event_created_at"`
	SessionFields              json.RawMessage `json:"session_fields"`
	Followed                   json.RawMessage `json:"followed"`
	Online                     json.RawMessage `json:"online"`
	State                      json.RawMessage `json:"state"`
	GroupIDs                   json.RawMessage `json:"group_ids"`
	EmailVerified              json.RawMessage `json:"email_verified"`
	CreatedAt                  json.RawMessage `json:"created_at"`
	Visibility                 json.RawMessage `json:"visibility"`
}

// Agent function converts User object to Agent object if User's Type is "agent".
// If Type is different or User is malformed, then it returns nil.
func (u *User) Agent() *Agent {
	if u.Type != "agent" {
		return nil
	}
	var a Agent

	a.User = u
	if err := internal.UnmarshalOptionalRawField(u.RoutingStatus, &a.RoutingStatus); err != nil {
		return nil
	}
	if err := internal.UnmarshalOptionalRawField(u.Visibility, &a.Visibility); err != nil {
		return nil
	}
	return &a
}

// Customer function converts User object to Customer object if User's Type is "customer".
// If Type is different or User is malformed, then it returns nil.
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
	if err := json.Unmarshal(u.EmailVerified, &c.EmailVerified); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.CreatedAt, &c.CreatedAt); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.Followed, &c.Followed); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.Online, &c.Online); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.State, &c.State); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.SessionFields, &c.SessionFields); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.GroupIDs, &c.GroupIDs); err != nil {
		return nil
	}
	return &c
}

// Chat represents LiveChat chat.
type Chat struct {
	ID         string     `json:"id,omitempty"`
	Properties Properties `json:"properties,omitempty"`
	Access     *Access    `json:"access,omitempty"`
	Thread     Thread     `json:"thread,omitempty"`
	Threads    []Thread   `json:"threads,omitempty"`
	IsFollowed bool       `json:"is_followed,omitempty"`
	Agents     map[string]*Agent
	Customers  map[string]*Customer
}

// Users function returns combined list of Chat's Agents and Customers.
func (c *Chat) Users() []*User {
	u := make([]*User, 0, len(c.Agents)+len(c.Customers))
	for _, a := range c.Agents {
		u = append(u, a.User)
	}
	for _, cu := range c.Customers {
		u = append(u, cu.User)
	}

	return u
}

// UnmarshalJSON implements json.Unmarshaler interface for Chat.
func (c *Chat) UnmarshalJSON(data []byte) error {
	type ChatAlias Chat
	var cs struct {
		*ChatAlias
		Users []json.RawMessage `json:"users"`
	}

	if err := json.Unmarshal(data, &cs); err != nil {
		return err
	}

	var t struct {
		Type string `json:"type"`
	}

	*c = (Chat)(*cs.ChatAlias)
	c.Agents = make(map[string]*Agent)
	c.Customers = make(map[string]*Customer)
	for _, u := range cs.Users {
		if err := json.Unmarshal(u, &t); err != nil {
			return err
		}
		switch t.Type {
		case "agent":
			var a Agent
			if err := json.Unmarshal(u, &a); err != nil {
				return err
			}
			c.Agents[a.ID] = &a
		case "customer":
			var cu Customer
			if err := json.Unmarshal(u, &cu); err != nil {
				return err
			}
			c.Customers[cu.ID] = &cu
		}
	}

	return nil
}

// Queue represents position of a thread in a queue
type Queue struct {
	Position int       `json:"position"`
	WaitTime int       `json:"wait_time"`
	QueuedAt time.Time `json:"queued_at,omitempty"`
}

// Thread represents LiveChat chat thread
type Thread struct {
	ID                        string     `json:"id"`
	Active                    bool       `json:"active"`
	UserIDs                   []string   `json:"user_ids"`
	RestrictedAccess          bool       `json:"restricted_access"`
	Properties                Properties `json:"properties"`
	Access                    *Access    `json:"access"`
	Tags                      []string   `json:"tags,omitempty"`
	Events                    []*Event   `json:"events"`
	PreviousThreadID          string     `json:"previous_thread_id"`
	NextThreadID              string     `json:"next_thread_id"`
	CreatedAt                 time.Time  `json:"created_at"`
	PreviousAccesibleThreadID string     `json:"previous_accessible_thread_id,omitempty"`
	NextAccessibleThreadID    string     `json:"next_accessible_thread_id,omitempty"`
	Queue                     *Queue     `json:"queue,omitempty"`
	QueuesDuration            *int       `json:"queues_duration,omitempty"`
}

// Access represents LiveChat chat and thread access
type Access struct {
	GroupIDs []int `json:"group_ids"`
}

// Agent represents LiveChat agent.
type Agent struct {
	*User
	RoutingStatus string `json:"routing_status,omitempty"`
	Visibility    string `json:"visibility,omitempty"`
}

// Visit contains information about particular customer's visit.
type Visit struct {
	IP          string `json:"ip"`
	UserAgent   string `json:"user_agent"`
	Geolocation struct {
		Country     string `json:"country"`
		CountryCode string `json:"country_code"`
		Region      string `json:"region"`
		City        string `json:"city"`
		Timezone    string `json:"timezone"`
	} `json:"geolocation"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
	Referrer  string    `json:"referrer"`
	LastPages []struct {
		OpenedAt time.Time `json:"opened_at"`
		URL      string    `json:"url"`
		Title    string    `json:"title"`
	} `json:"last_pages"`
}

// Customer represents LiveChat customer.
type Customer struct {
	*User
	EmailVerified bool  `json:"email_verified"`
	LastVisit     Visit `json:"last_visit"`
	Statistics    struct {
		VisitsCount            int `json:"visits_count"`
		ThreadsCount           int `json:"threads_count"`
		ChatsCount             int `json:"chats_count"`
		PageViewsCount         int `json:"page_views_count"`
		GreetingsShownCount    int `json:"greetings_shown_count"`
		GreetingsAcceptedCount int `json:"greetings_accepted_count"`
	} `json:"statistics"`
	AgentLastEventCreatedAt    time.Time           `json:"agent_last_event_created_at"`
	CustomerLastEventCreatedAt time.Time           `json:"customer_last_event_created_at"`
	CreatedAt                  time.Time           `json:"created_at"`
	SessionFields              []map[string]string `json:"session_fields"`
	Followed                   bool                `json:"followed"`
	Online                     bool                `json:"online"`
	State                      string              `json:"state"`
	GroupIDs                   []int               `json:"group_ids"`
}

// ValidateEvent checks if given interface resolves into supported event type
func ValidateEvent(e interface{}) error {
	switch v := e.(type) {
	case *Event:
	case *File:
	case *Message:
	case *RichMessage:
	case *SystemMessage:
	case Event:
	case File:
	case Message:
	case RichMessage:
	case SystemMessage:
	default:
		return fmt.Errorf("event type %T not supported", v)
	}

	return nil
}

type eventSpecific struct {
	Text              json.RawMessage `json:"text"`
	TextVars          json.RawMessage `json:"text_vars"`
	Fields            json.RawMessage `json:"fields"`
	ContentType       json.RawMessage `json:"content_type"`
	Name              json.RawMessage `json:"name"`
	URL               json.RawMessage `json:"url"`
	ThumbnailURL      json.RawMessage `json:"thumbnail_url"`
	Thumbnail2xURL    json.RawMessage `json:"thumbnail2x_url"`
	Width             json.RawMessage `json:"width"`
	Height            json.RawMessage `json:"height"`
	Size              json.RawMessage `json:"size"`
	TemplateID        json.RawMessage `json:"template_id"`
	Elements          json.RawMessage `json:"elements"`
	Postback          json.RawMessage `json:"postback"`
	AlternativeText   json.RawMessage `json:"alternative_text"`
	SystemMessageType json.RawMessage `json:"system_message_type"`
}

// Event represents base of all LiveChat chat events.
//
// To get specific event type's structure, call appropriate function based on Event's Type.
type Event struct {
	ID         string     `json:"id,omitempty"`
	CustomID   string     `json:"custom_id,omitempty"`
	CreatedAt  time.Time  `json:"created_at,omitempty"`
	AuthorID   string     `json:"author_id"`
	Properties Properties `json:"properties,omitempty"`
	Visibility string     `json:"visibility,omitempty"`
	Type       string     `json:"type,omitempty"`
	eventSpecific
}

// FilledForm represents LiveChat filled form event.
type FilledForm struct {
	Fields []struct {
		ID    string `json:"id"`
		Label string `json:"label"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"fields"`
	Event
}

// FilledForm function converts Event object to FilledForm object if Event's Type is "filled_form".
// If Type is different or Event is malformed, then it returns nil.
func (e *Event) FilledForm() *FilledForm {
	if e.Type != "filled_form" {
		return nil
	}
	var f FilledForm

	f.Event = *e
	if err := json.Unmarshal(e.Fields, &f.Fields); err != nil {
		return nil
	}
	return &f
}

// Postback represents postback data in LiveChat message event.
type Postback struct {
	ID       string `json:"id"`
	ThreadID string `json:"thread_id"`
	EventID  string `json:"event_id"`
	Type     string `json:"type,omitempty"`
	Value    string `json:"value,omitempty"`
}

// Message represents LiveChat message event.
type Message struct {
	Event
	Text     string    `json:"text,omitempty"`
	Postback *Postback `json:"postback,omitempty"`
}

// Message function converts Event object to Message object if Event's Type is "message".
// If Type is different or Event is malformed, then it returns nil.
func (e *Event) Message() *Message {
	if e.Type != "message" {
		return nil
	}
	var m Message

	m.Event = *e
	if err := json.Unmarshal(e.Text, &m.Text); err != nil {
		return nil
	}
	if err := internal.UnmarshalOptionalRawField(e.Postback, &m.Postback); err != nil {
		return nil
	}
	return &m
}

// SystemMessage represents LiveChat system message event.
type SystemMessage struct {
	Event
	SystemMessageType string            `json:"system_message_type"`
	Text              string            `json:"text,omitempty"`
	TextVars          map[string]string `json:"text_vars,omitempty"`
}

// SystemMessage function converts Event object to SystemMessage object if Event's Type is "system_message".
// If Type is different or Event is malformed, then it returns nil.
func (e *Event) SystemMessage() *SystemMessage {
	if e.Type != "system_message" {
		return nil
	}
	var sm SystemMessage

	sm.Event = *e
	if err := json.Unmarshal(e.SystemMessageType, &sm.SystemMessageType); err != nil {
		return nil
	}
	if err := internal.UnmarshalOptionalRawField(e.Text, &sm.Text); err != nil {
		return nil
	}
	if err := internal.UnmarshalOptionalRawField(e.TextVars, &sm.TextVars); err != nil {
		return nil
	}
	return &sm
}

// File represents LiveChat file event
type File struct {
	Event
	ContentType     string `json:"content_type"`
	Name            string `json:"name"`
	URL             string `json:"url"`
	ThumbnailURL    string `json:"thumbnail_url,omitempty"`
	Thumbnail2xURL  string `json:"thumbnail2x_url,omitempty"`
	Width           int    `json:"width,omitempty"`
	Height          int    `json:"height,omitempty"`
	Size            int    `json:"size,omitempty"`
	AlternativeText string `json:"alternative_text,omitempty"`
}

// File function converts Event object to File object if Event's Type is "file".
// If Type is different or Event is malformed, then it returns nil.
func (e *Event) File() *File {
	if e.Type != "file" {
		return nil
	}
	var f File

	f.Event = *e
	if err := json.Unmarshal(e.ContentType, &f.ContentType); err != nil {
		return nil
	}
	if err := json.Unmarshal(e.Name, &f.Name); err != nil {
		return nil
	}
	if err := json.Unmarshal(e.URL, &f.URL); err != nil {
		return nil
	}
	if err := internal.UnmarshalOptionalRawField(e.ThumbnailURL, &f.ThumbnailURL); err != nil {
		return nil
	}
	if err := internal.UnmarshalOptionalRawField(e.Thumbnail2xURL, &f.Thumbnail2xURL); err != nil {
		return nil
	}
	if err := internal.UnmarshalOptionalRawField(e.Width, &f.Width); err != nil {
		return nil
	}
	if err := internal.UnmarshalOptionalRawField(e.Height, &f.Height); err != nil {
		return nil
	}
	if err := internal.UnmarshalOptionalRawField(e.Size, &f.Size); err != nil {
		return nil
	}
	if err := internal.UnmarshalOptionalRawField(e.AlternativeText, &f.AlternativeText); err != nil {
		return nil
	}

	return &f
}

// RichMessage represents LiveChat rich message event
type RichMessage struct {
	Event
	TemplateID string               `json:"template_id"`
	Elements   []RichMessageElement `json:"elements"`
}

// RichMessageElement represents element of LiveChat rich message
type RichMessageElement struct {
	Buttons  []RichMessageButton `json:"buttons"`
	Title    string              `json:"title"`
	Subtitle string              `json:"subtitle"`
	Image    *RichMessageImage   `json:"image,omitempty"`
}

// RichMessageButton represents button in LiveChat rich message
type RichMessageButton struct {
	Text       string   `json:"text"`
	Type       string   `json:"type"`
	Value      string   `json:"value"`
	UserIds    []string `json:"user_ids"`
	PostbackID string   `json:"postback_id"`
	// Allowed values: compact, full, tall
	WebviewHeight string `json:"webview_height"`
	// Allowed values: new, current
	Target string `json:"target,omitempty"`
}

// RichMessageImage represents image in LiveChat rich message
type RichMessageImage struct {
	URL             string `json:"url"`
	Name            string `json:"name,omitempty"`
	ContentType     string `json:"content_type,omitempty"`
	Size            int    `json:"size,omitempty"`
	Width           int    `json:"width,omitempty"`
	Height          int    `json:"height,omitempty"`
	AlternativeText string `json:"alternative_text,omitempty"`
}

// RichMessage function converts Event object to RichMessage object if Event's Type is "rich_message".
// If Type is different or Event is malformed, then it returns nil.
func (e *Event) RichMessage() *RichMessage {
	if e.Type != "rich_message" {
		return nil
	}
	var rm RichMessage

	rm.Event = *e
	if err := json.Unmarshal(e.TemplateID, &rm.TemplateID); err != nil {
		return nil
	}
	if err := json.Unmarshal(e.Elements, &rm.Elements); err != nil {
		return nil
	}

	return &rm
}
