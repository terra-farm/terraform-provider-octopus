package octopus

import (
	"testing"
)

/*
 * Unit tests
 */

// Retrieve variables from a variable set by name.
func Test_Get_VariableSet_VariablesByName(test *testing.T) {
	expect := expect(test)

	variableSet := VariableSet{
		Variables: []Variable{
			Variable{
				Name: "Var1",
			},
			Variable{
				Name: "Var1",
			},
			Variable{
				Name: "Var2",
			},
			Variable{
				Name: "Var3",
			},
		},
	}

	variablesMatchingEnvironment := variableSet.GetVariablesByName("var1")
	expect.EqualsInt("len(variablesMatchingEnvironment)", 2, len(variablesMatchingEnvironment))
}

// Retrieve variables from a variable set by name and environment.
func Test_Get_VariableSet_VariablesByNameAndEnvironment(test *testing.T) {
	expect := expect(test)

	variableSet := VariableSet{
		Variables: []Variable{
			Variable{
				Name: "Var1",
			},
			Variable{
				Name: "Var1",
				Scope: VariableScopes{
					Environments: []string{
						"Env1",
					},
				},
			},
			Variable{ // Should be returned
				Name: "Var1",
				Scope: VariableScopes{
					Environments: []string{
						"Env1",
						"Env2",
					},
				},
			},
			Variable{
				Name: "Var1",
				Scope: VariableScopes{
					Environments: []string{
						"Env2",
					},
				},
			},
			Variable{ // Should be returned
				Name: "Var1",
				Scope: VariableScopes{
					Environments: []string{
						"Env2",
						"Env1",
					},
				},
			},
			Variable{
				Name: "Var2",
			},
			Variable{
				Name: "Var3",
			},
		},
	}

	variablesMatchingEnvironment := variableSet.GetVariablesByNameAndScopes("var1", VariableScopes{
		Environments: []string{"env2", "env1"},
	})
	expect.EqualsInt("len(variablesMatchingEnvironment)", 2, len(variablesMatchingEnvironment))
}
