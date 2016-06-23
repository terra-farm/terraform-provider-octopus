package octopus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

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

		return nil, errorResponse.ToError("Request to retrieve variable set '%s' failed with status code %d", id, statusCode)
	}

	variableSet = &VariableSet{}
	err = json.Unmarshal(responseBody, variableSet)

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
	request, err = client.newRequest(requestURI, http.MethodPost, variableSet)
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

// VariableSet represents a set of variables associated with an Octopus project.
type VariableSet struct {
	ID        string            `json:"Id"`
	OwnerID   string            `json:"OwnerId"`
	Version   int               `json:"Version"`
	Variables []Variable        `json:"Variables"`
	Links     map[string]string `json:"Links"`
}

// Variable represents a variable in an Octopus project.
type Variable struct {
	ID          string         `json:"Id"`
	Name        string         `json:"Name"`
	Value       string         `json:"Value"`
	Scope       VariableScopes `json:"ScopeValues"`
	IsSensitive bool           `json:"IsSensitive"`
	IsEditable  bool           `json:"IsEditable"`
}

// VariableScopes represents the scope(s) to which a variable applies.
type VariableScopes struct {
	Channels     []string `json:"Channel,omitempty"`
	Environments []string `json:"Environment,omitempty"`
	Roles        []string `json:"Role,omitempty"`
	Machines     []string `json:"Machine,omitempty"`
	Actions      []string `json:"Action,omitempty"`
	Projects     []string `json:"Project,omitempty"`
}

// GetVariablesByName retrieves all instances of a variable by name (regardless of scope).
func (variableSet *VariableSet) GetVariablesByName(name string) []Variable {
	matchingVariables := []Variable{}
	for _, variable := range variableSet.Variables {
		if variable.HasName(name) {
			matchingVariables = append(matchingVariables, variable)
		}
	}

	return matchingVariables
}

// GetVariablesByNameAndScope retrieves all instances of a variable by name and scope.
// Pass nil for any scope to match all values.
func (variableSet *VariableSet) GetVariablesByNameAndScope(name string, environment *string, role *string, machine *string, action *string, project *string) []Variable {
	matchingVariables := []Variable{}
	for _, variable := range variableSet.Variables {
		if variable.HasName(name) && variable.MatchesScope(environment, role, machine, action, project) {
			matchingVariables = append(matchingVariables, variable)
		}
	}

	return matchingVariables
}

// HasName determines whether a variable has the specified name (case-insensitive).
func (variable Variable) HasName(name string) bool {
	return strings.ToLower(variable.Name) == strings.ToLower(name)
}

// MatchesScope determines whether a variable matches the specified scope(s).
// Passing nil for a scope will ignore that scope.
func (variable Variable) MatchesScope(environment *string, role *string, machine *string, action *string, project *string) bool {
	if !scopesContain(variable.Scope.Environments, environment) {
		return false
	}

	if !scopesContain(variable.Scope.Roles, role) {
		return false
	}

	if !scopesContain(variable.Scope.Machines, machine) {
		return false
	}

	if !scopesContain(variable.Scope.Actions, action) {
		return false
	}

	if !scopesContain(variable.Scope.Projects, project) {
		return false
	}

	return true
}

func scopesContain(scopes []string, value *string) bool {
	// Nil means match any value.
	if value == nil {
		return true
	}

	// Empty scope list means match any value.
	if len(scopes) == 0 {
		return true
	}

	matchValue := strings.ToLower(*value)

	for _, scope := range scopes {
		if strings.ToLower(scope) == matchValue {
			return true
		}
	}

	return false
}
