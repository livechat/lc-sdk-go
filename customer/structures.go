package customer

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
