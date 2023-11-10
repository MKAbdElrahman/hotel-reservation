package types

import "fmt"

type HTTPError struct {
	Description string `json:"description,omitempty"`
	Metadata    string `json:"metadata,omitempty"`
	StatusCode  int    `json:"statusCode"`
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("description: %s,  metadata: %s", e.Description, e.Metadata)
}

func NewHTTPError(description, metadata string, statusCode int) *HTTPError {
	return &HTTPError{
		Description: description,
		Metadata:    metadata,
		StatusCode:  statusCode,
	}
}
