package errors

import "fmt"

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
