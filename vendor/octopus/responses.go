package octopus

import "fmt"

// APIResponse represents the basic response most commonly received when making API calls.
type APIResponse struct {
	// The API status message (if any).
	Message string `json:"message"`
}

// ToError creates an error representing the API response.
func (response *APIResponse) ToError(errorMessageOrFormat string, formatArgs ...interface{}) error {
	return &APIError{
		Message:  fmt.Sprintf(errorMessageOrFormat, formatArgs...),
		Response: *response,
	}
}

// APIError is an error representing an error response from an API.
type APIError struct {
	Message  string
	Response APIResponse
}

// Error returns the error message associated with the APIError.
func (apiError *APIError) Error() string {
	return apiError.Message
}

var _ error = &APIError{}

// TODO: Well-known API response codes
