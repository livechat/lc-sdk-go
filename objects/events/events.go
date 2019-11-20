package events

import (
	"encoding/json"
	"time"
)

type Properties map[string]map[string]interface{}

type eventSpecific struct {
	Text   json.RawMessage `json:"text"`
	Fields json.RawMessage `json:"fields"`
}

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

type FilledForm struct {
	Fields []struct {
		Label string `json:"label"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"fields"`
	*Event
}

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

type Message struct {
	*Event
	Text string `json:"text,omitempty"`
}

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

type SystemMessage struct {
	Event
	Type string `json:"system_message_type,omitempty"`
	Text string `json:"text,omitempty"`
}

type File struct {
	Event
	ContentType    string `json:"content_type"`
	URL            string `json:"url"`
	Width          int    `json:"width,omitempty"`
	Height         int    `json:"height,omitempty"`
	ThumbnailURL   string `json:"thumbnail_url,omitempty"`
	Thumbnail2xURL string `json:"thumbnail2x_url,omitempty"`
}
