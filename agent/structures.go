package agent

import "github.com/livechat/lc-sdk-go/objects"

type Filters struct {
	IncludeActive bool               `json:"include_active,omitempty"`
	GroupIDs      []int32            `json:"group_ids,omitempty"`
	AgentIDs      []int32            `json:"agent_ids,omitempty"`
	ThreadIDs     []int32            `json:"thread_ids,omitempty"`
	Properties    objects.Properties `json:"properties,omitempty"`
	Query         string             `json:"query,omitempty"`
	DateFrom      string             `json:"date_from,omitempty"`
	DateTo        string             `json:"date_to,omitempty"`
}

// ThreadSummary represents a short summary of a thread.
type ThreadSummary struct {
	ID          string `json:"id"`
	Order       int32  `json:"order"`
	TotalEvents uint   `json:"total_events"`
}

// Form struct describes schema of custom form (e-mail, prechat or postchat survey).
type Form struct {
	ID     string `json:"id"`
	Fields []struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		Label    string `json:"label"`
		Required bool   `json:"required"`
	} `json:"fields"`
}

// PredictedAgent is an agent returned by GetPredictedAgent method.
type PredictedAgent struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar"`
	IsBot     bool   `json:"is_bot"`
	JobTitle  string `json:"job_title"`
	Type      string `json:"type"`
}

// URLDetails contains some OpenGraph details of the URL.
type URLDetails struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	ImageURL    string `json:"image_url"`
	ImageWidth  int    `json:"image_width"`
	ImageHeight int    `json:"image_height"`
}

// InitialChat represents initial chat used in StartChat or ActivateChat.
type InitialChat struct {
	ID         string             `json:"id"`
	Access     *objects.Access    `json:"access,omitempty"`
	Properties objects.Properties `json:"properties,omitempty"`
	Users      []*objects.User    `json:"users,omitempty"`
	Thread     *InitialThread     `json:"thread,omitempty"`
}

// InitialThread represents initial chat thread used in StartChat or ActivateChat.
type InitialThread struct {
	Events     []interface{}      `json:"events,omitempty"`
	Properties objects.Properties `json:"properties,omitempty"`
}
