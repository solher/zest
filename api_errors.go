package zest

import "fmt"

// APIError defines the format of Zest API errors.
type APIError struct {
	// The status code.
	Status int `json:"status"`
	// The description of the API error.
	Description string `json:"description"`
	// A raw description of what triggered the API error.
	Raw string `json:"raw"`
	// The token uniquely identifying the API error.
	ErrorCode string `json:"errorCode"`
	// Additional infos.
	Params map[string]interface{} `json:"params,omitempty"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("%s : %s", e.ErrorCode, e.Description)
}
