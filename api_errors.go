package zest

import "fmt"

// APIError is the struct defining the format of Zest API errors.
type APIError struct {
	Status      int    `json:"status"`
	Description string `json:"description"`
	Raw         string `json:"raw"`
	ErrorCode   string `json:"errorCode"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s : %s", e.ErrorCode, e.Description)
}
