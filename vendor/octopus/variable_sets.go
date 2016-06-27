package octopus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// VariableSet represents a set of variables associated with an Octopus project.
type VariableSet struct {
	ID          string                 `json:"Id"`
	OwnerID     string                 `json:"OwnerId"`
	Version     int                    `json:"Version"`
	Variables   []Variable             `json:"Variables"`
	ScopeValues VariableSetScopeValues `json:"ScopeValues"`
	Links       map[string]string      `json:"Links"`
}

// VariableSetScopeValues represents summary information for the entities referenced by a VariableSet's scope values.
type VariableSetScopeValues struct {
	Channels     []EntitySummary `json:"Channels,omitempty"`
	Environments []EntitySummary `json:"Environments,omitempty"`
	Roles        []EntitySummary `json:"Roles,omitempty"`
	Machines     []EntitySummary `json:"Machines,omitempty"`
	Actions      []EntitySummary `json:"Actions,omitempty"`
	Projects     []EntitySummary `json:"Projects,omitempty"`
}

// GetVariableByID retrieves a specific instance of a variable by Id.
func (variableSet *VariableSet) GetVariableByID(id string) *Variable {
	index := variableSet.GetVariableIndexByID(id)
	if index != -1 {
		return &variableSet.Variables[index]
	}

	return nil
}

// GetVariableIndexByID retrieves the index of a specific instance of a variable by Id.
func (variableSet *VariableSet) GetVariableIndexByID(id string) int {
	for index, variable := range variableSet.Variables {
		if variable.ID == id {
			return index
		}
	}

	return -1
}

// UpdateVariableFunc represents a function that updates a variable
type UpdateVariableFunc func(variable *Variable)

// UpdateVariable updates a specific instance of a variable by Id.
func (variableSet *VariableSet) UpdateVariable(id string, update UpdateVariableFunc) bool {
	index := variableSet.GetVariableIndexByID(id)
	if index == -1 {
		return false
	}

	variable := &variableSet.Variables[index]
	update(variable)
	variableSet.Variables[index] = *variable

	return true
}

// GetVariablesByName retrieves all instances of a variable by name (regardless of scope).
func (variableSet *VariableSet) GetVariablesByName(name string) []Variable {
	return filterVariablesByName(variableSet.Variables, name)
}

// GetVariablesByNameAndScopes retrieves all instances of a variable by name and scope.
func (variableSet *VariableSet) GetVariablesByNameAndScopes(name string, scopes VariableScopes) []Variable {
	return filterVariablesByNameAndScopes(variableSet.Variables, name, scopes)
}

// GetVariableSet retrieves an Octopus variable set by Id.
func (client *Client) GetVariableSet(id string) (variableSet *VariableSet, err error) {
	var (
		request       *http.Request
		statusCode    int
		responseBody  []byte
		errorResponse *APIErrorResponse
	)

	requestURI := fmt.Sprintf("variables/%s", id)
	request, err = client.newRequest(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	responseBody, statusCode, err = client.executeRequest(request)
	if err != nil {
		err = fmt.Errorf("Error invoking request to read variable set '%s': %s", id, err.Error())

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

		return nil, errorResponse.ToError("Request to retrieve variable set '%s' failed with status code %d", id, statusCode)
	}

	variableSet = &VariableSet{}
	err = json.Unmarshal(responseBody, variableSet)
	if err != nil {
		err = fmt.Errorf("Invalid response detected when reading variable set '%s': %s", id, err.Error())
	}

	return
}

// UpdateVariableSet updates an Octopus variable set.
func (client *Client) UpdateVariableSet(variableSet *VariableSet) (updatedVariableSet *VariableSet, err error) {
	if variableSet == nil {
		return nil, fmt.Errorf("Must supply the variable set to update.")
	}

	var (
		request       *http.Request
		statusCode    int
		responseBody  []byte
		errorResponse *APIErrorResponse
	)

	requestURI := fmt.Sprintf("variables/%s", variableSet.ID)
	request, err = client.newRequest(requestURI, http.MethodPut, variableSet)
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

		return nil, errorResponse.ToError("Request to update variable set '%s' failed with status code %d", variableSet.ID, statusCode)
	}

	updatedVariableSet = &VariableSet{}
	err = json.Unmarshal(responseBody, updatedVariableSet)

	return
}

// GetProjectVariableSet retrieves the variable set associated with an Octopus project.
func (client *Client) GetProjectVariableSet(projectID string) (variableSet *VariableSet, err error) {
	var project *Project

	project, err = client.GetProject(projectID)
	if err != nil {
		return
	}

	if project == nil {
		return nil, fmt.Errorf("Project '%s' not found.", projectID)
	}

	return client.GetVariableSet(project.VariableSetID)
}
