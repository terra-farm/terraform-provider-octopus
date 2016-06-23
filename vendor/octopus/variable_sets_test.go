package octopus

import (
	"testing"
)

// Unit test
//
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

// Unit test
//
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
			Variable{
				Name: "Var1",
				Scope: VariableScopes{
					Environments: []string{
						"Env2",
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

	variablesMatchingEnvironment := variableSet.GetVariablesByNameAndScope("var1", stringToPtr("env1"), nil, nil, nil, nil)
	expect.EqualsInt("len(variablesMatchingEnvironment)", 2, len(variablesMatchingEnvironment))
}
