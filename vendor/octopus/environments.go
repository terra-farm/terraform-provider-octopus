package octopus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Environment represents an Octopus environment.
type Environment struct {
	ID               string            `json:"Id"`
	Name             string            `json:"Name"`
	Description      string            `json:"Description"`
	SortOrder        int               `json:"SortOrder"`
	UseGuidedFailure bool              `json:"UseGuidedFailure"`
	Links            map[string]string `json:"Links"`
}

// Request body when creating a new environment.
type createEnvironment struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	SortOrder   int    `json:"SortOrder"`
}

// GetEnvironments retrieves a page of Octopus environments.
//
// skip indicates the number of results to skip over.
// Call Environments.GetSkipForNextPage() / Environments.GetSkipForPreviousPage() to get the number of items to skip for the next / previous page of results.
func (client *Client) GetEnvironments(skip int) (environments []Environment, err error) {
	var (
		request       *http.Request
		statusCode    int
		responseBody  []byte
		errorResponse *APIErrorResponse
	)

	requestURI := fmt.Sprintf("environments/all?skip=%d", skip)
	request, err = client.newRequest(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	responseBody, statusCode, err = client.executeRequest(request)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		errorResponse, err = readAPIErrorResponseAsJSON(responseBody, statusCode)
		if err != nil {
			return nil, err
		}

		return nil, errorResponse.ToError("Request to retrieve all environments failed with status code %d.", statusCode)
	}

	err = json.Unmarshal(responseBody, &environments)
	if err != nil {
		return
	}

	return
}

// GetEnvironment retrieves a specific environment by Id.
func (client *Client) GetEnvironment(id string) (environment *Environment, err error) {
	var (
		request       *http.Request
		statusCode    int
		responseBody  []byte
		errorResponse *APIErrorResponse
	)

	requestURI := fmt.Sprintf("environments/%s", id)
	request, err = client.newRequest(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	responseBody, statusCode, err = client.executeRequest(request)
	if err != nil {
		return nil, err
	}

	if statusCode == http.StatusNotFound {
		// Environment not found.
		return nil, nil
	}

	if statusCode != http.StatusOK {
		errorResponse, err = readAPIErrorResponseAsJSON(responseBody, statusCode)
		if err != nil {
			return nil, err
		}

		return nil, errorResponse.ToError("Request to retrieve environment '%s' failed with status code %d.", id, statusCode)
	}

	environment = &Environment{}
	err = json.Unmarshal(responseBody, environment)

	return
}

// CreateEnvironment creates a new environment by Id.
func (client *Client) CreateEnvironment(name string, description string, sortOrder int) (environment *Environment, err error) {
	var (
		request       *http.Request
		statusCode    int
		responseBody  []byte
		errorResponse *APIErrorResponse
	)

	request, err = client.newRequest("environments", http.MethodPost, &createEnvironment{
		Name:        name,
		Description: description,
		SortOrder:   sortOrder,
	})
	if err != nil {
		return nil, err
	}

	responseBody, statusCode, err = client.executeRequest(request)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusCreated {
		errorResponse, err = readAPIErrorResponseAsJSON(responseBody, statusCode)
		if err != nil {
			return nil, err
		}

		return nil, errorResponse.ToError("Request to create environment '%s' failed with status code %d", name, statusCode)
	}

	environment = &Environment{}
	err = json.Unmarshal(responseBody, environment)

	return
}

// UpdateEnvironment updates a specific environment by Id.
func (client *Client) UpdateEnvironment(environment *Environment) (updatedEnvironment *Environment, err error) {
	var (
		request       *http.Request
		statusCode    int
		responseBody  []byte
		errorResponse *APIErrorResponse
	)

	if environment == nil {
		return nil, fmt.Errorf("Must supply an environment to update.")
	}

	requestURI := fmt.Sprintf("environments/%s", environment.ID)
	request, err = client.newRequest(requestURI, http.MethodPut, environment)
	if err != nil {
		return nil, err
	}
	responseBody, statusCode, err = client.executeRequest(request)
	if err != nil {
		return nil, err
	}

	if statusCode == http.StatusNotFound {
		// Environment not found.
		return nil, nil
	}

	if statusCode != http.StatusOK {
		errorResponse, err = readAPIErrorResponseAsJSON(responseBody, statusCode)
		if err != nil {
			return nil, err
		}

		return nil, errorResponse.ToError("Request to update environment '%s' failed with status code %d", environment.ID, statusCode)
	}

	updatedEnvironment = &Environment{}
	err = json.Unmarshal(responseBody, updatedEnvironment)

	return
}

// DeleteEnvironment deletes an environment.
func (client *Client) DeleteEnvironment(id string) (err error) {
	var (
		request       *http.Request
		statusCode    int
		responseBody  []byte
		errorResponse *APIErrorResponse
	)

	requestURI := fmt.Sprintf("environments/%s", id)
	request, err = client.newRequest(requestURI, http.MethodDelete, nil)
	if err != nil {
		return err
	}
	responseBody, statusCode, err = client.executeRequest(request)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		errorResponse, err = readAPIErrorResponseAsJSON(responseBody, statusCode)
		if err != nil {
			return err
		}

		return errorResponse.ToError("Request to delete environment '%s' failed with status code %d", id, statusCode)
	}

	return nil
}
