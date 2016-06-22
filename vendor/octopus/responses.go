package octopus

import "fmt"

// APIErrorResponse represents the basic error response most commonly received when making API calls.
type APIErrorResponse struct {
	Message string   `json:"ErrorMessage"`
	Errors  []string `json:"Errors"`
}

// ToError creates an error representing the API response.
func (response *APIErrorResponse) ToError(errorMessageOrFormat string, formatArgs ...interface{}) error {
	return &APIError{
		Message:  fmt.Sprintf(errorMessageOrFormat, formatArgs...),
		Response: *response,
	}
}

// APIError is an error representing an error response from an API.
type APIError struct {
	Message  string
	Response APIErrorResponse
}

// Error returns the error message associated with the APIError.
func (apiError *APIError) Error() string {
	return fmt.Sprintf("%s: %s", apiError.Message, apiError.Response.Message)
}

var _ error = &APIError{}
