package octopus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ProjectGroup represents a group of related Octopus projects.
type ProjectGroup struct {
	ID                string            `json:"Id"`
	Name              string            `json:"Name"`
	Description       string            `json:"Description"`
	EnvironmentIDs    []string          `json:"EnvironmentIds"`
	RetentionPolicyID string            `json:"RetentionPolicyId"`
	Links             map[string]string `json:"Links"`
}

// ProjectGroups represents a page of ProjectGroup results.
type ProjectGroups struct {
	Items []ProjectGroup `json:"Items"`

	PagedResults
}

// GetProjectGroups retrieves a page of Octopus project groups.
//
// skip indicates the number of results to skip over.
// Call ProjectGroups.GetSkipForNextPage() / ProjectGroups.GetSkipForPreviousPage() to get the number of items to skip for the next / previous page of results.
func (client *Client) GetProjectGroups(skip int) (projectGroups *ProjectGroups, err error) {
	var (
		request       *http.Request
		statusCode    int
		responseBody  []byte
		errorResponse *APIErrorResponse
	)

	requestURI := fmt.Sprintf("projectGroups?skip=%d", skip)
	request, err = client.newRequest(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	responseBody, statusCode, err = client.executeRequest(request)
	if err != nil {
		err = fmt.Errorf("Error invoking request to read project groups: %s", err.Error())

		return
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

		return nil, errorResponse.ToError("Request to retrieve project groups failed with status code %d.", statusCode)
	}

	projectGroups = &ProjectGroups{}
	err = json.Unmarshal(responseBody, projectGroups)
	if err != nil {
		err = fmt.Errorf("Invalid response detected when reading project groups: %s", err.Error())
	}

	return
}

// GetProjectGroup retrieves an Octopus projectGroup group by Id or slug.
func (client *Client) GetProjectGroup(id string) (projectGroup *ProjectGroup, err error) {
	var (
		request       *http.Request
		statusCode    int
		responseBody  []byte
		errorResponse *APIErrorResponse
	)

	requestURI := fmt.Sprintf("projectGroups/%s", id)
	request, err = client.newRequest(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	responseBody, statusCode, err = client.executeRequest(request)
	if err != nil {
		err = fmt.Errorf("Error invoking request to read project group '%s': %s", id, err.Error())

		return
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

		return nil, errorResponse.ToError("Request to retrieve project group '%s' failed with status code %d.", id, statusCode)
	}

	projectGroup = &ProjectGroup{}
	err = json.Unmarshal(responseBody, projectGroup)
	if err != nil {
		err = fmt.Errorf("Invalid response detected when reading project group '%s': %s", id, err.Error())
	}

	return
}
