// Package objects defines common LiveChat structures.
//
// General common LiveChat structures documentation is available here:
// https://developers.livechatinc.com/docs/messaging/customer-chat-api/#other-common-structures
package objects

import (
	"encoding/json"
	"fmt"
	"time"
)

func unmarshalOptionalRawField(source json.RawMessage, target interface{}) error {
	if source != nil {
		return json.Unmarshal(source, target)
	}
	return nil
}

// Properties represents LiveChat properties in form of property_namespace -> property -> value.
type Properties map[string]map[string]interface{}

// User represents base of both Customer and Agent
//
// To get speficic user type's structure, call Agent() or Customer() (based on Type value).
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
}

// Agent function converts User object to Agent object if User's Type is "agent".
// If Type is different or User is malformed, then it returns nil.
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
	return &c
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
	LastPages []struct {
		OpenedAt time.Time `json:"opened_at"`
		URL      string    `json:"url"`
		Title    string    `json:"title"`
	} `json:"last_pages"`
}

// Chat represents LiveChat chat.
type Chat struct {
	ID         string     `json:"id,omitempty"`
	Properties Properties `json:"properties,omitempty"`
	Access     Access     `json:"access,omitempty"`
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

// Thread represents LiveChat chat thread
type Thread struct {
	ID               string     `json:"id"`
	Active           bool       `json:"active"`
	UserIDs          []string   `json:"user_ids"`
	RestrictedAccess bool       `json:"restricted_access"`
	Properties       Properties `json:"properties"`
	Access           Access     `json:"access"`
	Events           []*Event   `json:"events"`
	PreviousThreadID string     `json:"previous_thread_id"`
	NextThreadID     string     `json:"next_thread_id"`
	CreatedAt        time.Time  `json:"created_at"`
}

// Access represents LiveChat chat and thread access
type Access struct {
	GroupIDs []int `json:"group_ids"`
}

// Agent represents LiveChat agent.
type Agent struct {
	*User
	RoutingStatus string `json:"routing_status"`
}

// Customer represents LiveChat customer.
type Customer struct {
	*User
	LastVisit  Visit `json:"last_visit"`
	Statistics struct {
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
}

// ThreadSummary represents a short summary of a thread
type ThreadSummary struct {
	ID         string     `json:"id"`
	UserIDs    []string   `json:"user_ids"`
	Properties Properties `json:"properties,omitempty"`
	Active     bool       `json:"active"`
	Access     Access     `json:"access,omitempty"`
	Tags       []string   `json:"tags,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// ChatSummary represents a short summary of a chat
type ChatSummary struct {
	ID               string `json:"id"`
	LastEventPerType map[string]struct {
		ThreadID        string    `json:"thread_id"`
		ThreadCreatedAt time.Time `json:"thread_created_at"`
		Event           Event     `json:"event"`
	} `json:"last_event_per_type,omitempty"`
	Users             []*User        `json:"users"`
	LastThreadSummary *ThreadSummary `json:"last_thread_summary,omitempty"`
	Properties        Properties     `json:"properties,omitempty"`
	Access            Access         `json:"access,omitempty"`
	Order             uint64         `json:"order,omitempty"`
	IsFollowed        bool           `json:"is_followed"`
}

// InitialChat represents initial chat used in StartChat or ActivateChat.
type InitialChat struct {
	ID         string         `json:"id"`
	Access     *Access        `json:"access,omitempty"`
	Properties Properties     `json:"properties,omitempty"`
	Thread     *InitialThread `json:"thread,omitempty"`
}

// Validate checks if there are no unsupported event types in InitialChat Thread
func (chat *InitialChat) Validate() error {
	if chat.Thread != nil {
		for _, e := range chat.Thread.Events {
			if err := ValidateEvent(e); err != nil {
				return err
			}
		}
	}
	return nil
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

// InitialThread represents initial chat thread used in StartChat or ActivateChat.
type InitialThread struct {
	Events     []interface{} `json:"events,omitempty"`
	Properties Properties    `json:"properties,omitempty"`
}

type eventSpecific struct {
	Text        json.RawMessage `json:"text"`
	Fields      json.RawMessage `json:"fields"`
	ContentType json.RawMessage `json:"content_type"`
	URL         json.RawMessage `json:"url"`
	Width       json.RawMessage `json:"width"`
	Height      json.RawMessage `json:"height"`
	Name        json.RawMessage `json:"name"`
	TemplateID  json.RawMessage `json:"template_id"`
	Elements    json.RawMessage `json:"elements"`
	Postback    json.RawMessage `json:"postback"`
}

// Event represents base of all LiveChat chat events.
//
// To get speficic event type's structure, call appropriate function based on Event's Type.
type Event struct {
	ID         string     `json:"id,omitempty"`
	CustomID   string     `json:"custom_id,omitempty"`
	CreatedAt  time.Time  `json:"created_at,omitempty"`
	AuthorID   string     `json:"author_id"`
	Properties Properties `json:"properties,omitempty"`
	Recipients string     `json:"recipients,omitempty"`
	Type       string     `json:"type,omitempty"`
	eventSpecific
}

// FilledForm represents LiveChat filled form event.
type FilledForm struct {
	Fields []struct {
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
	if err := unmarshalOptionalRawField(e.Postback, &m.Postback); err != nil {
		return nil
	}
	return &m
}

// SystemMessage represents LiveChat system message event.
type SystemMessage struct {
	Event
	Type     string            `json:"system_message_type,omitempty"`
	Text     string            `json:"text,omitempty"`
	TextVars map[string]string `json:"text_vars,omitempty"`
}

// File represents LiveChat file event
type File struct {
	Event
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Name        string `json:"name"`
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
	if err := json.Unmarshal(e.URL, &f.URL); err != nil {
		return nil
	}
	if err := unmarshalOptionalRawField(e.Width, &f.Width); err != nil {
		return nil
	}
	if err := unmarshalOptionalRawField(e.Height, &f.Height); err != nil {
		return nil
	}
	if err := json.Unmarshal(e.Name, &f.Name); err != nil {
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
	UserIds    []string `json:"user_ids"`
	Value      string   `json:"value"`
	PostbackID string   `json:"postback_id"`
	// Allowed values: compact, full, tall
	WebviewHeight string `json:"webview_height"`
	// Allowed values: new, current
	Target string `json:"target,omitempty"`
}

// RichMessageImage represents image in LiveChat rich message
type RichMessageImage struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
	Size        int    `json:"size"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
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
