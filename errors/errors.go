package errors

import "fmt"

// ErrAPI represents structure of errors returned by all LiveChat APIs (configuration, agent chat and customer chat APIs).
type ErrAPI struct {
	Details *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
	StatusCode int
}

func (e *ErrAPI) Error() string {
	if e.Details == nil {
		return ""
	}
	return fmt.Sprintf("API error: %s - %s", e.Details.Type, e.Details.Message)
}
