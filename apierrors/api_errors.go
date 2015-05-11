package apierrors

import (
	"fmt"
	"strings"
)

type APIError struct {
	Status      int    `json:"status"`
	Description string `json:"description"`
	Raw         string `json:"raw"`
	ErrorCode   string `json:"errorCode"`
}

func (apiError *APIError) Error() string {
	return fmt.Sprintf("%s : %s", apiError.ErrorCode, apiError.Description)
}

var (
	BodyDecodingError   = &APIError{Description: "Could not decode the JSON request.", ErrorCode: "BODY_DECODING_ERROR"}
	FilterDecodingError = &APIError{Description: "Could not decode the JSON query filter.", ErrorCode: "FILTER_DECODING_ERROR"}
	InvalidPathParams   = &APIError{Description: "Invalid path parameters.", ErrorCode: "INVALID_PATH_PARAMS"}
)

var (
	InvalidCredentials  = &APIError{Description: "Invalid credentials.", ErrorCode: "INVALID_CREDENTIALS"}
	SessionNotFound     = &APIError{Description: "No active sessions were found.", ErrorCode: "SESSION_NOT_FOUND"}
	DatabaseError       = &APIError{Description: "An error occured with the database. Please retry later.", ErrorCode: "SESSION_CREATION_ERROR"}
	InternalServerError = &APIError{Description: "An internal error occured. Please retry later.", ErrorCode: "INTERNAL_SERVER_ERROR"}
	Unauthorized        = &APIError{Description: "Authorization Required.", ErrorCode: "AUTHORIZATION_REQUIRED"}
)

func BlankParam(paramName string) *APIError {
	titleName := strings.Title(strings.ToLower(paramName))
	upperName := strings.ToUpper(paramName)
	return &APIError{Description: titleName + " can't be blank.", ErrorCode: "BLANK_" + upperName}
}

func Make(apiError APIError, status int, err error) *APIError {
	if err != nil {
		apiError.Raw = err.Error()
	}

	apiError.Status = status

	return &apiError
}
