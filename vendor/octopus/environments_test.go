package octopus

import (
	"testing"
)

/*
 * Integration tests
 */

// Create environment (successful).
func Test_Client_CreateEnvironment_Success(test *testing.T) {
	testClientRequest(test, &ClientTest{
		APIKey: "my-test-api-key",
		Invoke: func(test *testing.T, client *Client) {
			environment, err := client.CreateEnvironment("TerraformTest", "Terraform test environment", 0)
			if err != nil {
				test.Fatal(err)
			}

			verifyCreateEnvironmentTestResponse(test, environment)
		},
		Handle: testRespondCreated(createEnvironmentTestResponse),
	})
}

/*
 * Test responses.
 */

const createEnvironmentTestResponse = `
	{
		"Id": "Environments-444",
		"Name": "TerraformTest",
		"Description": "Terraform test environment",
		"SortOrder": 91,
		"UseGuidedFailure": false,
		"LastModifiedOn": "2016-06-29T06:04:30.232+00:00",
		"LastModifiedBy": "foo@bar.com.cloud",
		"Links": {
			"Self": "/api/environments/Environments-444",
			"Machines": "/api/environments/Environments-444/machines"
		}
	}
`

func verifyCreateEnvironmentTestResponse(test *testing.T, environment *Environment) {
	expect := expect(test)

	expect.NotNil("Environment", environment)
	expect.EqualsString("Environment.ID", "Environments-444", environment.ID)
	expect.EqualsString("Environment.Name", "TerraformTest", environment.Name)
}
