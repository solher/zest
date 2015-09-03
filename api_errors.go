package zest

import "fmt"

type APIError struct {
	Status      int    `json:"status"`
	Description string `json:"description"`
	Raw         string `json:"raw"`
	ErrorCode   string `json:"errorCode"`
}

func (apiError *APIError) Error() string {
	return fmt.Sprintf("%s : %s", apiError.ErrorCode, apiError.Description)
}
