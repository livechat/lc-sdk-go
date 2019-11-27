// Package objects defines common LiveChat structures.
//
// General common LiveChat structures documentation is available here:
// https://developers.livechatinc.com/docs/messaging/customer-chat-api/#other-common-structures
package objects

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/livechat/lc-sdk-go/objects/events"
)

// Properties represents LiveChat properties in form of property_namespace -> property -> value.
type Properties map[string]map[string]interface{}

// User represents base of both Customer and Agent
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
	ID               string          `json:"id"`
	Timestamp        Time            `json:"timestamp"`
	Active           bool            `json:"active"`
	UserIDs          []string        `json:"user_ids"`
	RestrictedAccess bool            `json:"restricted_access"`
	Order            int             `json:"order"`
	Properties       Properties      `json:"properties"`
	Access           Access          `json:"access"`
	Events           []*events.Event `json:"events"`
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

// ThreadSummary represents a short summary of a thread.
type ThreadSummary struct {
	ID          string `json:"id"`
	Order       int32  `json:"order"`
	TotalEvents uint   `json:"total_events"`
}

// ChatSummary represents a short summary of a chat
type ChatSummary struct {
	ID                string         `json:"id"`
	LastEventPerType  interface{}    `json:"last_event_per_type,omitempty"`
	Users             []interface{}  `json:"users"`
	LastThreadSummary *ThreadSummary `json:"last_thread_summary,omitempty"`
	Properties        Properties     `json:"properties,omitempty"`
	Access            interface{}    `json:"access,omitempty"`
	Order             uint64         `json:"order,omitempty"`
	IsFollowed        bool           `json:"is_followed"`
}

// InitialChat represents initial chat used in StartChat or ActivateChat.
type InitialChat struct {
	ID         string         `json:"id"`
	Access     *Access        `json:"access,omitempty"`
	Properties Properties     `json:"properties,omitempty"`
	Users      []*User        `json:"users,omitempty"`
	Thread     *InitialThread `json:"thread,omitempty"`
}

// InitialThread represents initial chat thread used in StartChat or ActivateChat.
type InitialThread struct {
	Events     []interface{} `json:"events,omitempty"`
	Properties Properties    `json:"properties,omitempty"`
}
