package octopus

import (
	"strings"
)

// Variable represents a variable in an Octopus project.
type Variable struct {
	ID          string         `json:"Id"`
	Name        string         `json:"Name"`
	Value       string         `json:"Value"`
	Scope       VariableScopes `json:"Scope"`
	IsSensitive bool           `json:"IsSensitive"`
	IsEditable  bool           `json:"IsEditable"`
}

// HasName determines whether a variable has the specified name (case-insensitive).
func (variable Variable) HasName(name string) bool {
	return strings.ToLower(variable.Name) == strings.ToLower(name)
}

// MatchesScopes determines whether a variable matches the specified scope(s).
func (variable Variable) MatchesScopes(matchScopes VariableScopes) bool {
	return variable.Scope.MatchesScopes(matchScopes)
}

// Return all instances of a variable that have the specified name (regardless of scope).
//
// Comparisons are case-insensitive.
func filterVariablesByName(variables []Variable, name string) []Variable {
	matchingVariables := []Variable{}
	for _, variable := range variables {
		if variable.HasName(name) {
			matchingVariables = append(matchingVariables, variable)
		}
	}

	return matchingVariables
}

// Return all instances of a variable that match the specified by name and scopes.
//
// Comparisons are case-insensitive.
func filterVariablesByNameAndScopes(variables []Variable, name string, scopes VariableScopes) []Variable {
	matchingVariables := []Variable{}
	for _, variable := range variables {
		if variable.HasName(name) && variable.MatchesScopes(scopes) {
			matchingVariables = append(matchingVariables, variable)
		}
	}

	return matchingVariables
}
