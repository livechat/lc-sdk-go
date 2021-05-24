package customer

import "github.com/livechat/lc-sdk-go/v3/objects"

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

type DynamicConfiguration struct {
	GroupID             int    `json:"group_id"`
	ClientLimitExceeded bool   `json:"client_limit_exceeded"`
	DomainAllowed       bool   `json:"domain_allowed"`
	ConfigVersion       string `json:"config_version"`
	LocalizationVersion string `json:"localization_version"`
	Language            string `json:"language"`
}

type ConfigButton struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	OnlineValue  string `json:"online_value"`
	OfflineValue string `json:"offline_value"`
}

type Configuration struct {
	Buttons        []ConfigButton               `json:"buttons"`
	TicketForm     *Form                        `json:"ticket_form,omitempty"`
	PrechatForm    *Form                        `json:"prechat_form,omitempty"`
	AllowedDomains []string                     `json:"allowed_domains,omitempty"` // CHECK: why not in docs?
	Integrations   map[string]map[string]string `json:"integrations"`
	Properties     struct {
		Group   objects.Properties `json:"group"`
		License objects.Properties `json:"license"`
	} `json:"properties"`
}
