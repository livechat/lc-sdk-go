package errors

import "fmt"

type ErrAPI struct {
	ErrorDetails *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
	StatusCode int
}

func (e *ErrAPI) Error() string {
	if e.ErrorDetails == nil {
		return ""
	}
	return fmt.Sprintf("API error: %s - %s", e.ErrorDetails.Type, e.ErrorDetails.Message)
}
