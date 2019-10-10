package customer

import "time"

type Properties map[string]string

type Chat struct {
	ID         string            `json:"id,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
	Access     Access            `json:"access,omitempty"`
	Users      []User            `json:"users,omitempty"`
	Threads    []Thread          `json:"threads,omitempty"`
}

type Thread struct {
	ID               string     `json:"id"`
	Timestamp        time.Time  `json:"timestamp"`
	Active           bool       `json:"active"`
	UserIDs          []string   `json:"user_ids"`
	RestrictedAccess bool       `json:"restricted_access"`
	Order            int32      `json:"order"`
	Properties       Properties `json:"properties"`
	Access           Access     `json:"access"`
}

type Access struct {
	GroupIDs []int32 `json:"group_ids"`
}

type User struct {
	ID       string    `json:"id"`
	Type     string    `json:"type"`
	Name     string    `json:"name"`
	Avatar   string    `json:"avatar"`
	Present  bool      `json:"present"`
	LastSeen time.Time `json:"last_seen_timestamp"`
}
