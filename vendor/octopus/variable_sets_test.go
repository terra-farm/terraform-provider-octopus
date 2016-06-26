package octopus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
 * Integration tests
 */

// Get VLAN by Id (successful).
func TestClient_GetVariableSet_Success(test *testing.T) {
	expect := expect(test)

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		apiKeyHeader, ok := request.Header["X-Octopus-Apikey"] // Note casing - header name gets normalised
		expect.IsTrue("header 'X-Octopus-ApiKey' is present", ok)
		expect.EqualsInt("Header 'X-Octopus-ApiKey' value count", 1, len(apiKeyHeader))
		expect.EqualsString("Header 'X-Octopus-ApiKey'", "my-test-api-key", apiKeyHeader[0])

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		fmt.Fprintln(writer, getVariableSetTestResponse)
	}))
	defer testServer.Close()

	client, err := NewClientWithAPIKey(testServer.URL, "my-test-api-key")
	if err != nil {
		test.Fatal(err)
	}

	variableSet, err := client.GetVariableSet("my-variable-set")
	if err != nil {
		test.Fatal(err)
	}

	verifyGetVatiableSetTestResponse(test, variableSet)
}

/*
 * Unit tests
 */

// Retrieve variables from a variable set by name.
func Test_VariableSet_GetVariablesByName(test *testing.T) {
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
func Test_VariableSet_GetVariablesByNameAndEnvironment(test *testing.T) {
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

/*
 * Test responses.
 */

const getVariableSetTestResponse = `
	{
		"Id": "variableset-Projects-105",
		"OwnerId": "Projects-105",
		"Version": 173,
		"Variables": [
			{
				"Id": "e775471f-cc48-731c-9a91-e8099581ad93",
				"Name": "SqlServerInstanceName",
				"Value": "my-server.lab.au.my-net.cloud",
				"Scope": {
					"Environment": [
						"Environments-130",
						"Environments-131"
					]
				},
				"IsSensitive": false,
				"IsEditable": true,
				"Prompt": null
			},
			{
				"Id": "56876c2b-016b-54b4-0499-4497df7ffb3e",
				"Name": "AuditingDatabase",
				"Value": "Data Source=#{SqlServerInstanceName};Initial Catalog=#{AuditingDatabaseName};UID=#{SqlServerLogin};PWD=#{SqlServerPassword}",
				"Scope": {
					"Environment": [
						"Environments-130",
						"Environments-131"
					]
				},
				"IsSensitive": false,
				"IsEditable": true,
				"Prompt": null
			}
		],
		"ScopeValues": {
			"Environments": [
				{
					"Id": "Environments-130",
					"Name": "Platform R2.0 Development AU"
				},
				{
					"Id": "Environments-131",
					"Name": "Platform R2.0 Automation AU"
				}
			],
			"Machines": [],
			"Actions": [],
			"Roles": [],
			"Channels": []
		},
		"Links": {
			"Self": "/api/variables/variableset-Projects-105"
		}
	}
`

func verifyGetVatiableSetTestResponse(test *testing.T, variableSet *VariableSet) {
	expect := expect(test)

	expect.NotNil("VariableSet", variableSet)
	expect.EqualsString("VariableSet.ID", "variableset-Projects-105", variableSet.ID)
	expect.EqualsString("VariableSet.OwnerID", "Projects-105", variableSet.OwnerID)
	expect.EqualsInt("VariableSet.Version", 173, variableSet.Version)
	expect.EqualsInt("VariableSet.Variables.Length", 2, len(variableSet.Variables))

	variable1 := variableSet.Variables[0]
	expect.EqualsString("VariableSet.Variables[0].ID", "e775471f-cc48-731c-9a91-e8099581ad93", variable1.ID)
	expect.EqualsString("VariableSet.Variables[0].Name", "SqlServerInstanceName", variable1.Name)

	variable2 := variableSet.Variables[1]
	expect.EqualsString("VariableSet.Variables[1].ID", "56876c2b-016b-54b4-0499-4497df7ffb3e", variable2.ID)
	expect.EqualsString("VariableSet.Variables[1].Name", "AuditingDatabase", variable2.Name)
}
