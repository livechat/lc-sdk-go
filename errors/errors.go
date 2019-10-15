package errors

import "fmt"

type ErrAPI struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	StatusCode int
}

func (e *ErrAPI) Error() string {
	return fmt.Sprintf("API error: %s - %s", e.Type, e.Message)
}
