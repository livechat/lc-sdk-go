package customer

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/livechat/lc-sdk-go/internal/events"
)

type Properties map[string]map[string]interface{}

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(t.Unix(), 10)), nil
}

func (t *Time) UnmarshalJSON(s []byte) (err error) {
	q, err := strconv.ParseInt(string(s), 10, 64)
	if err != nil {
		return err
	}
	*t = Time{Time: time.Unix(q, 0)}
	return
}

type Chat struct {
	ID         string     `json:"id,omitempty"`
	Properties Properties `json:"properties,omitempty"`
	Access     Access     `json:"access,omitempty"`
	Threads    []Thread   `json:"threads,omitempty"`
	Agents     map[string]*Agent
	Customers  map[string]*Customer
}

func (c *Chat) Users() []User {
	u := make([]User, len(c.Agents)+len(c.Customers))
	var i int
	for _, a := range c.Agents {
		u[i] = a.User
		i += 1
	}
	for _, cu := range c.Customers {
		u[i] = cu.User
		i += 1
	}

	return u
}

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

type Events []*events.Event

type Thread struct {
	ID               string     `json:"id"`
	Timestamp        Time       `json:"timestamp"`
	Active           bool       `json:"active"`
	UserIDs          []string   `json:"user_ids"`
	RestrictedAccess bool       `json:"restricted_access"`
	Order            int        `json:"order"`
	Properties       Properties `json:"properties"`
	Access           Access     `json:"access"`
	Events           Events     `json:"events"`
}

type Access struct {
	GroupIDs []int `json:"group_ids"`
}

type User struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Present  bool   `json:"present"`
	LastSeen Time   `json:"last_seen_timestamp"`
}

type Agent struct {
	User
	RoutingStatus string `json:"routing_status"`
}

type Customer struct {
	User
	LastVisit  Visit `json:"last_visit"`
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
}

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
