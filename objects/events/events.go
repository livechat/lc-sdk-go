package events

import (
	"encoding/json"
	"time"
)

type Properties map[string]map[string]interface{}

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
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
	Name        string `json:"name"`
}

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
	if err := json.Unmarshal(e.Width, &f.Width); err != nil {
		return nil
	}
	if err := json.Unmarshal(e.Height, &f.Height); err != nil {
		return nil
	}
	if err := json.Unmarshal(e.Name, &f.Name); err != nil {
		return nil
	}

	return &f
}

type RichMessage struct {
	Event
	TemplateID string               `json:"template_id"`
	Elements   []RichMessageElement `json:"elements"`
}

type RichMessageElement struct {
	Buttons  []RichMessageButton `json:"buttons"`
	Title    string              `json:"title"`
	Subtitle string              `json:"subtitle"`
	Image    RichMessageImage    `json:"image"`
}

type RichMessageButton struct {
	Text    string   `json:"text"`
	Type    string   `json:"type"`
	UserIds []string `json:"user_ids"`
	Value   string   `json:"value"`
	// Allowed values: compact, full, tall
	WebviewHeight string `json:"webview_height"`
}

type RichMessageImage struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
	Size        int    `json:"size"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
}

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
