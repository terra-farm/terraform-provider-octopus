package octopus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Projects represents a page of Project results.
type Projects struct {
	Items []Project `json:"Items"`

	PagedResults
}

// Project represents an Octopus project.
type Project struct {
	ID                              string                    `json:"Id"`
	Name                            string                    `json:"Name"`
	Description                     string                    `json:"Description"`
	Slug                            string                    `json:"Slug"`
	VersioningStrategy              ProjectVersioningStrategy `json:"VersioningStrategy"`
	VariableSetID                   string                    `json:"VariableSetId"`
	IncludedLibraryVariableSetIDs   []string                  `json:"IncludedLibraryVariableSetIds"`
	ProjectGroupID                  string                    `json:"ProjectGroupId"`
	LifeCycleID                     string                    `json:"LifeCycleId"`
	IsDisabled                      bool                      `json:"IsDisabled"`
	AutoCreateRelease               bool                      `json:"AutoCreateRelease"`
	DefaultToSkipIfAlreadyInstalled bool                      `json:"DefaultToSkipIfAlreadyInstalled"`
	Links                           map[string]string         `json:"Links"`
}

// ProjectVersioningStrategy represents the versioning strategy for an Octopus project.
type ProjectVersioningStrategy struct {
	DonorPackageStepID string `json:"DonorPackageStepId"`
	Template           string `json:"Template"`
}

// GetProjects retrieves a page of Octopus projects.
//
// skip indicates the number of results to skip over.
// Call Projects.GetSkipForNextPage() / Projects.GetSkipForPreviousPage() to get the number of items to skip for the next / previous page of results.
func (client *Client) GetProjects(skip int) (projects *Projects, err error) {
	var (
		request       *http.Request
		statusCode    int
		responseBody  []byte
		errorResponse *APIErrorResponse
	)

	requestURI := fmt.Sprintf("projects?skip=%d", skip)
	request, err = client.newRequest(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	responseBody, statusCode, err = client.executeRequest(request)
	if err != nil {
		err = fmt.Errorf("Error invoking request to read projects: %s", err.Error())

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

		return nil, errorResponse.ToError("Request to retrieve projects failed with status code %d.", statusCode)
	}

	projects = &Projects{}
	err = json.Unmarshal(responseBody, projects)
	if err != nil {
		err = fmt.Errorf("Invalid response detected when reading projects: %s", err.Error())
	}

	return
}

// GetProject retrieves an Octopus project by Id or slug.
func (client *Client) GetProject(idOrSlug string) (project *Project, err error) {
	var (
		request       *http.Request
		statusCode    int
		responseBody  []byte
		errorResponse *APIErrorResponse
	)

	requestURI := fmt.Sprintf("projects/%s", idOrSlug)
	request, err = client.newRequest(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	responseBody, statusCode, err = client.executeRequest(request)
	if err != nil {
		err = fmt.Errorf("Error invoking request to read project '%s': %s", idOrSlug, err.Error())

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

		return nil, errorResponse.ToError("Request to retrieve project '%s' failed with status code %d.", idOrSlug, statusCode)
	}

	project = &Project{}
	err = json.Unmarshal(responseBody, project)
	if err != nil {
		err = fmt.Errorf("Invalid response detected when reading project '%s': %s", idOrSlug, err.Error())
	}

	return
}

// UpdateProject updates an Octopus project.
func (client *Client) UpdateProject(project *Project) (updatedProject *Project, err error) {
	if project == nil {
		return nil, fmt.Errorf("Must supply the project to update.")
	}

	var (
		request       *http.Request
		statusCode    int
		responseBody  []byte
		errorResponse *APIErrorResponse
	)

	requestURI := fmt.Sprintf("projects/%s", project.ID)
	request, err = client.newRequest(requestURI, http.MethodPost, project)
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

		return nil, errorResponse.ToError("Request to update project '%s' failed with status code %d", project.ID, statusCode)
	}

	updatedProject = &Project{}
	err = json.Unmarshal(responseBody, updatedProject)

	return
}
