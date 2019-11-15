// Package objects defines common LiveChat structures.
//
// General common LiveChat structures documentation is available here:
// https://developers.livechatinc.com/docs/messaging/customer-chat-api/#other-common-structures
package objects

import (
	"encoding/json"
	"strconv"
	"time"
)

// Properties represents LiveChat properties in form of property_namespace -> property -> value.
type Properties map[string]map[string]interface{}

// User represents base of both Customer and Agent objects.
//
// To get speficic user type's structure, call Agent() or Customer() (based on Type value).
type User struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Present  bool   `json:"present"`
	LastSeen Time   `json:"last_seen_timestamp"`
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

// The Time type is a helper type to convert string time into golang time representation.
type Time struct {
	time.Time
}

// MarshalJSON implements json.Marshaler interface for Time.
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(t.Unix(), 10)), nil
}

// UnmarshalJSON implements json.Unmarshaler interface for Time.
func (t *Time) UnmarshalJSON(s []byte) (err error) {
	q, err := strconv.ParseInt(string(s), 10, 64)
	if err != nil {
		return err
	}
	*t = Time{Time: time.Unix(q, 0)}
	return
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
	Threads    []Thread   `json:"threads,omitempty"`
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
	Timestamp        Time       `json:"timestamp"`
	Active           bool       `json:"active"`
	UserIDs          []string   `json:"user_ids"`
	RestrictedAccess bool       `json:"restricted_access"`
	Order            int        `json:"order"`
	Properties       Properties `json:"properties"`
	Access           Access     `json:"access"`
	Events           []*Event   `json:"events"`
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
	AgentLastEventCreatedAt    time.Time         `json:"agent_last_event_created_at"`
	CustomerLastEventCreatedAt time.Time         `json:"customer_last_event_created_at"`
	CreatedAt                  time.Time         `json:"created_at"`
	Fields                     map[string]string `json:"fields"`
}

type eventSpecific struct {
	Text   json.RawMessage `json:"text"`
	Fields json.RawMessage `json:"fields"`
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
	*Event
}

// FilledForm function converts Event object to FilledForm object if Event's Type is "filled_form".
// If Type is different or Event is malformed, then it returns nil.
func (e *Event) FilledForm() *FilledForm {
	if e.Type != "filled_form" {
		return nil
	}
	var f FilledForm

	f.Event = e
	if err := json.Unmarshal(e.Fields, &f.Fields); err != nil {
		return nil
	}
	return &f
}

// Message represents LiveChat message event.
type Message struct {
	*Event
	Text string `json:"text,omitempty"`
}

// Message function converts Event object to Message object if Event's Type is "message".
// If Type is different or Event is malformed, then it returns nil.
func (e *Event) Message() *Message {
	if e.Type != "message" {
		return nil
	}
	var m Message

	m.Event = e
	if err := json.Unmarshal(e.Text, &m.Text); err != nil {
		return nil
	}
	return &m
}

// SystemMessage represents LiveChat system message event.
type SystemMessage struct {
	Event
	Type string `json:"system_message_type,omitempty"`
	Text string `json:"text,omitempty"`
}
