package errors

import "fmt"

// ErrAPI represents structure of errors returned by all LiveChat APIs (configuration, agent chat and customer chat APIs).
type ErrAPI struct {
	Details    *ErrDetails `json:"error"`
	StatusCode int
}

type ErrDetails struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (e *ErrAPI) Error() string {
	if e.Details == nil {
		return ""
	}
	return fmt.Sprintf("API responded with status code %d: error: %s - %s", e.StatusCode, e.Details.Type, e.Details.Message)
}
