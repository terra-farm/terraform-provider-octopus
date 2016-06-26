package octopus

import (
	"strings"
)

// VariableScopes represents the scope(s) to which a variable applies.
type VariableScopes struct {
	Channels     []string `json:"Channel,omitempty"`
	Environments []string `json:"Environment,omitempty"`
	Roles        []string `json:"Role,omitempty"`
	Machines     []string `json:"Machine,omitempty"`
	Actions      []string `json:"Action,omitempty"`
	Projects     []string `json:"Project,omitempty"`
}

// MatchesScopes determines whether a variable matches the specified scope(s).
// Passing nil for a scope will ignore that scope.
func (scopes VariableScopes) MatchesScopes(matchScopes VariableScopes) bool {
	if !areScopeValuesEquivalent(scopes.Environments, matchScopes.Environments) {
		return false
	}

	if !areScopeValuesEquivalent(scopes.Roles, matchScopes.Roles) {
		return false
	}

	if !areScopeValuesEquivalent(scopes.Machines, matchScopes.Machines) {
		return false
	}

	if !areScopeValuesEquivalent(scopes.Actions, matchScopes.Actions) {
		return false
	}

	if !areScopeValuesEquivalent(scopes.Projects, matchScopes.Projects) {
		return false
	}

	return true
}

// Are the specified sets of scope values considered equivalent (ignoring case and order)?
func areScopeValuesEquivalent(values1 []string, values2 []string) bool {
	if len(values1) != len(values2) {
		return false
	}

	combinedValues := make(map[string]bool, len(values1))
	for _, value1 := range values1 {
		combinedValues[strings.ToLower(value1)] = true
	}
	for _, value2 := range values2 {
		combinedValues[strings.ToLower(value2)] = true
	}

	return len(combinedValues) == len(values1)
}
