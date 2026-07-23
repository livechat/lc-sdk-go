package customer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/livechat/lc-sdk-go/v6/internal"
)

// Properties represents LiveChat properties in form of property_namespace -> property -> value.
type Properties map[string]map[string]interface{}

// Form struct describes schema of custom form (e-mail, prechat or postchat survey).
type Form struct {
	ID     string `json:"id"`
	Fields []struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		Label    string `json:"label"`
		Required bool   `json:"required"`
		Options  []struct {
			ID    string `json:"id"`
			Type  int    `json:"group_id"`
			Label string `json:"label"`
		} `json:"options"`
	} `json:"fields"`
}

// PredictedAgent is an agent returned by GetPredictedAgent method.
type PredictedAgent struct {
	Agent struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar"`
		IsBot     bool   `json:"is_bot"`
		JobTitle  string `json:"job_title"`
		Type      string `json:"type"`
	} `json:"agent"`
	Queue bool `json:"queue"`
}

// PredictedAgentV2 is an agent returned by RequestWelcomeMessage method.
type PredictedAgentV2 struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	IsBot    bool   `json:"is_bot"`
	BotType  string `json:"bot_type,omitempty"`
	JobTitle string `json:"job_title"`
	Type     string `json:"type"`
}

// URLInfo contains some OpenGraph info of the URL.
type URLInfo struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	URL              string `json:"url"`
	ImageURL         string `json:"image_url"`
	ImageOriginalURL string `json:"image_original_url"`
	ImageWidth       int    `json:"image_width"`
	ImageHeight      int    `json:"image_height"`
}

type ConfigButton struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	OnlineValue  string `json:"online_value"`
	OfflineValue string `json:"offline_value"`
}

// User represents base of both Customer and Agent
//
// To get specific user type's structure, call Agent() or Customer() (based on Type value).
type User struct {
	ID             string    `json:"id"`
	Type           string    `json:"type"`
	Name           string    `json:"name"`
	Avatar         string    `json:"avatar"`
	Present        bool      `json:"present"`
	EventsSeenUpTo time.Time `json:"events_seen_up_to"`
	userSpecific
}

type userSpecific struct {
	SessionFields json.RawMessage `json:"session_fields,omitempty"`
	Email         json.RawMessage `json:"email,omitempty"`
	EmailVerified json.RawMessage `json:"email_verified,omitempty"`
	PhoneNumber   json.RawMessage `json:"phone_number,omitempty"`
	Address       json.RawMessage `json:"address,omitempty"`
}

// Agent function converts User object to Agent object if User's Type is "agent".
// If Type is different or User is malformed, then it returns nil.
func (u *User) Agent() *Agent {
	if u.Type != "agent" {
		return nil
	}

	return &Agent{
		User: u,
	}
}

// Customer function converts User object to Customer object if User's Type is "customer".
// If Type is different or User is malformed, then it returns nil.
func (u *User) Customer() *Customer {
	if u.Type != "customer" {
		return nil
	}
	var c Customer

	c.User = u
	if err := json.Unmarshal(u.Email, &c.Email); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.EmailVerified, &c.EmailVerified); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.SessionFields, &c.SessionFields); err != nil {
		return nil
	}
	if err := internal.UnmarshalOptionalRawField(u.PhoneNumber, &c.PhoneNumber); err != nil {
		return nil
	}
	if err := internal.UnmarshalOptionalRawField(u.Address, &c.Address); err != nil {
		return nil
	}
	return &c
}

// Chat represents LiveChat chat.
type Chat struct {
	ID         string     `json:"id,omitempty"`
	Properties Properties `json:"properties,omitempty"`
	Access     *Access    `json:"access,omitempty"`
	Thread     *Thread    `json:"thread,omitempty"`
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

// Queue represents position of a thread in a queue
type Queue struct {
	Position int       `json:"position"`
	WaitTime int       `json:"wait_time"`
	QueuedAt time.Time `json:"queued_at"`
}

// Thread represents LiveChat chat thread
type Thread struct {
	ID               string     `json:"id"`
	Active           bool       `json:"active"`
	UserIDs          []string   `json:"user_ids"`
	Properties       Properties `json:"properties"`
	Access           *Access    `json:"access"`
	Events           []*Event   `json:"events"`
	PreviousThreadID string     `json:"previous_thread_id"`
	NextThreadID     string     `json:"next_thread_id"`
	CreatedAt        time.Time  `json:"created_at"`
	Queue            *Queue     `json:"queue,omitempty"`
}

// Access represents LiveChat chat and thread access
type Access struct {
	GroupIDs []int `json:"group_ids"`
}

// Agent represents LiveChat agent.
type Agent struct {
	*User
}

// Customer represents LiveChat customer.
type Customer struct {
	*User
	NameIsDefault bool                `json:"name_is_default"`
	Email         string              `json:"email,omitempty"`
	PhoneNumber   string              `json:"phone_number,omitempty"`
	EmailVerified bool                `json:"email_verified,omitempty"`
	SessionFields []map[string]string `json:"session_fields,omitempty"`
	Address       *Address            `json:"address,omitempty"`
}

type Address struct {
	Address    string `json:"address,omitempty"`
	City       string `json:"city,omitempty"`
	Country    string `json:"country,omitempty"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
}

// ChatInfo represents a short description of a chat
type ChatInfo struct {
	ID                  string    `json:"id"`
	LastThreadCreatedAt time.Time `json:"last_thread_created_at"`
	LastThreadID        string    `json:"last_thread_id,omitempty"`
	LastEventPerType    map[string]struct {
		ThreadID        string    `json:"thread_id"`
		ThreadCreatedAt time.Time `json:"thread_created_at"`
		Event           Event     `json:"event"`
	} `json:"last_event_per_type,omitempty"`
	Users      []*User    `json:"users"`
	Access     *Access    `json:"access,omitempty"`
	Properties Properties `json:"properties,omitempty"`
	Active     bool       `json:"active"`
}

// InitialThread represents initial chat thread used in StartChat or ResumeChat.
type InitialThread struct {
	Events     []interface{} `json:"events,omitempty"`
	Properties Properties    `json:"properties,omitempty"`
}

// InitialChat represents initial chat used in StartChat or ResumeChat.
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
	case *System:
	case Event:
	case File:
	case Message:
	case RichMessage:
	case SystemMessage:
	case System:
	default:
		return fmt.Errorf("event type %T not supported", v)
	}

	return nil
}

type eventSpecific struct {
	Text              json.RawMessage `json:"text"`
	TextVars          json.RawMessage `json:"text_vars"`
	Fields            json.RawMessage `json:"fields"`
	FormType          json.RawMessage `json:"form_type"`
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
	Source            json.RawMessage `json:"source"`
	Subtype           json.RawMessage `json:"subtype"`
	Details           json.RawMessage `json:"details"`
	Version           json.RawMessage `json:"version"`
}

// Event represents base of all LiveChat chat events.
//
// To get specific event type's structure, call appropriate function based on Event's Type.
type Event struct {
	ID         string     `json:"id,omitempty"`
	CustomID   string     `json:"custom_id,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	AuthorID   string     `json:"author_id"`
	Properties Properties `json:"properties,omitempty"`
	Recipients string     `json:"recipients,omitempty"`
	Type       string     `json:"type,omitempty"`
	eventSpecific
}

// FilledForm represents LiveChat filled form event.
type FilledForm struct {
	Fields   []FilledFormField `json:"fields"`
	FormType string            `json:"form_type"`
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
	if err := internal.UnmarshalOptionalRawField(e.FormType, &f.FormType); err != nil {
		return nil
	}
	return &f
}

// FilledFormField represents a field in LiveChat filled form event.
type FilledFormField struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Label string `json:"label"`
	filledFormFieldSpecific
}

type filledFormFieldSpecific struct {
	Answer  json.RawMessage `json:"answer,omitempty"`
	Answers json.RawMessage `json:"answers,omitempty"`
}

// FilledFormFieldSingle represents a field in LiveChat filled form event with a single answer (e.g. text field).
type FilledFormFieldSingle struct {
	FilledFormField
	Answer string `json:"answer"`
}

// FilledFormFieldSingleChoice represents a field in LiveChat filled form event with a single choice answer (e.g. radio button).
type FilledFormFieldSingleChoice struct {
	FilledFormField
	Answer struct {
		ID    string `json:"id"`
		Label string `json:"label"`
	} `json:"answer"`
}

// FilledFormFieldMultiChoice represents a field in LiveChat filled form event with a multiple choice answer (e.g. checkbox).
type FilledFormFieldMultiChoice struct {
	FilledFormField
	Answers []struct {
		ID    string `json:"id"`
		Label string `json:"label"`
	} `json:"answer"`
}

// FilledFormEmailCheckbox represents a field in LiveChat filled form event with a checkbox for email answer.
type FilledFormEmailCheckbox struct {
	FilledFormField
	Answer bool `json:"answer"`
}

// FilledFormFieldGroupChooser represents a field in LiveChat filled form event with a group_chooser answer.
type FilledFormFieldGroupChooser struct {
	FilledFormField
	Answer struct {
		ID      string `json:"id"`
		Label   string `json:"label"`
		GroupID int    `json:"group_id"`
	} `json:"answer"`
}

// Single method converts FilledFormField object to FilledFormFieldSingle object
// if FilledFormField's Type is one of "name", "email", "question", "textarea", "subject".
// If Type is different or FilledFormField is malformed, then it returns nil.
func (f *FilledFormField) Single() *FilledFormFieldSingle {
	supportedTypes := map[string]struct{}{
		"name":     {},
		"email":    {},
		"question": {},
		"textarea": {},
		"subject":  {},
	}
	if _, ok := supportedTypes[f.Type]; !ok {
		return nil
	}
	var s FilledFormFieldSingle
	s.ID, s.Label = f.ID, f.Label
	if err := json.Unmarshal(f.Answer, &s.Answer); err != nil {
		return nil
	}
	return &s
}

// SingleChoice method converts FilledFormField object to FilledFormFieldSingleChoice object
// if FilledFormField's Type is "radio" or "select".
// If Type is different or FilledFormField is malformed, then it returns nil.
func (f *FilledFormField) SingleChoice() *FilledFormFieldSingleChoice {
	supportedTypes := map[string]struct{}{
		"radio":  {},
		"select": {},
	}
	if _, ok := supportedTypes[f.Type]; !ok {
		return nil
	}
	var sc FilledFormFieldSingleChoice
	sc.ID, sc.Label = f.ID, f.Label
	if err := json.Unmarshal(f.Answer, &sc.Answer); err != nil {
		return nil
	}
	return &sc
}

// MultiChoice method converts FilledFormField object to FilledFormFieldMultiChoice object
// if FilledFormField's Type is "checkbox".
// If Type is different or FilledFormField is malformed, then it returns nil.
func (f *FilledFormField) MultiChoice() *FilledFormFieldMultiChoice {
	if f.Type != "checkbox" {
		return nil
	}
	var mc FilledFormFieldMultiChoice
	mc.ID, mc.Label = f.ID, f.Label
	if err := json.Unmarshal(f.Answers, &mc.Answers); err != nil {
		return nil
	}
	return &mc
}

// EmailCheckbox method converts FilledFormField object to FilledFormEmailCheckbox object
// if FilledFormField's Type is "checkbox_for_email".
// If Type is different or FilledFormField is malformed, then it returns nil.
func (f *FilledFormField) EmailCheckbox() *FilledFormEmailCheckbox {
	if f.Type != "checkbox_for_email" {
		return nil
	}
	var ec FilledFormEmailCheckbox
	ec.ID, ec.Label = f.ID, f.Label
	if err := json.Unmarshal(f.Answer, &ec.Answer); err != nil {
		return nil
	}
	return &ec
}

// GroupChooser method converts FilledFormField object to FilledFormFieldGroupChooser object
// if FilledFormField's Type is "group_chooser".
// If Type is different or FilledFormField is malformed, then it returns nil.
func (f *FilledFormField) GroupChooser() *FilledFormFieldGroupChooser {
	if f.Type != "group_chooser" {
		return nil
	}
	var gc FilledFormFieldGroupChooser
	gc.ID, gc.Label = f.ID, f.Label
	if err := json.Unmarshal(f.Answer, &gc.Answer); err != nil {
		return nil
	}
	return &gc
}

// Postback represents postback data in LiveChat message event.
type Postback struct {
	ID        string             `json:"id"`
	ThreadID  string             `json:"thread_id"`
	EventID   string             `json:"event_id"`
	Type      string             `json:"type,omitempty"`
	Value     string             `json:"value,omitempty"`
	Ecommerce *PostbackEcommerce `json:"ecommerce,omitempty"`
}

type PostbackEcommerce struct {
	ProductID string `json:"product_id"`
	OptionID  string `json:"option_id,omitempty"`
	Quantity  int    `json:"quantity,omitempty"`
}

// Message represents LiveChat message event.
type Message struct {
	Event
	Text     string    `json:"text"`
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
	Buttons   []RichMessageButton   `json:"buttons"`
	Title     string                `json:"title"`
	Subtitle  string                `json:"subtitle"`
	Image     *RichMessageImage     `json:"image,omitempty"`
	Ecommerce *RichMessageEcommerce `json:"ecommerce,omitempty"`
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

// RichMessageEcommerce represents ecommerce element in LiveChat rich message
type RichMessageEcommerce struct {
	ProductID string                       `json:"product_id"`
	Label     string                       `json:"label"`
	ViewType  string                       `json:"view_type"`
	Options   []RichMessageEcommerceOption `json:"options,omitempty"`
	Addons    []RichMessageEcommerceAddon  `json:"addons,omitempty"`
}

type RichMessageEcommerceOption struct {
	OptionID          string `json:"option_id"`
	Label             string `json:"label"`
	Price             string `json:"price,omitempty"`
	RegularPrice      string `json:"regular_price,omitempty"`
	Currency          string `json:"currency,omitempty"`
	Color             string `json:"color,omitempty"`
	ImageURL          string `json:"image_url,omitempty"`
	ImageThumbnailURL string `json:"image_thumbnail_url,omitempty"`
	Available         *bool  `json:"available,omitempty"`
	Selected          *bool  `json:"selected,omitempty"`
}

type RichMessageEcommerceAddon struct {
	AddonType string `json:"addon_type"`
	RangeFrom string `json:"range_from,omitempty"`
	RangeTo   string `json:"range_to,omitempty"`
	Currency  string `json:"currency,omitempty"`
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

// System represents LiveChat system event (replacement for the system_message)
type System struct {
	Event
	Source  string `json:"source"`
	Subtype string `json:"subtype"`
	Details string `json:"details"`
	Version int32  `json:"version"`
}

// System function converts Event object to System object if Event's Type is "system".
// If Type is different or Event is malformed, then it returns nil.
func (e *Event) System() *System {
	if e.Type != "system" {
		return nil
	}
	var s System

	s.Event = *e
	if err := json.Unmarshal(e.Source, &s.Source); err != nil {
		return nil
	}
	if err := json.Unmarshal(e.Subtype, &s.Subtype); err != nil {
		return nil
	}
	if err := json.Unmarshal(e.Details, &s.Details); err != nil {
		return nil
	}
	if err := json.Unmarshal(e.Version, &s.Version); err != nil {
		return nil
	}
	return &s
}

type AgentStatus struct {
	AgentID string `json:"agent_id,omitempty"`
	Status  string `json:"status,omitempty"`
}
