package octopus

import (
	"testing"
)

/*
 * Integration tests
 */

// Get machine by Id (successful).
func Test_Client_GetMachine_Success(test *testing.T) {
	testClientRequest(test, &ClientTest{
		APIKey: "my-test-api-key",
		Invoke: func(test *testing.T, client *Client) {
			machine, err := client.GetMachine("Machines-1")
			if err != nil {
				test.Fatal(err)
			}

			verifyGetMachineTestResponse(test, machine)
		},
		Handle: testRespondOK(getMachineTestResponse),
	})
}

/*
 * Test responses.
 */

const getMachineTestResponse = `
	{
		"Id": "Machines-1",
		"Name": "my-server.lab.au.test.cloud",
		"Thumbprint": "B3092FEA722388E326CFF4D2F6E124B5727AEEDC",
		"Uri": "https://10.110.21.15:10933/",
		"IsDisabled": false,
		"EnvironmentIds": [
			"Environments-1",
			"Environments-2"
		],
		"Roles": [
			"auditing-db",
			"identity-db",
			"provisioning-db",
			"semantic-logger"
		],
		"Status": "NeedsUpgrade",
		"HasLatestCalamari": true,
		"StatusSummary": "This machine is running an old version of Tentacle (3.2.19).",
		"Endpoint": {
			"CommunicationStyle": "TentaclePassive",
			"Uri": "https://10.110.21.15:10933/",
			"Thumbprint": "B3092FEA722388E326CFF4D2F6E124B5727AEEDC",
			"TentacleVersionDetails": {
				"UpgradeLocked": false,
				"Version": "3.2.19",
				"UpgradeSuggested": true,
				"UpgradeRequired": false
			},
			"Id": null,
			"LastModifiedOn": null,
			"LastModifiedBy": null,
			"Links": {}
		},
		"Links": {
			"Self": "/api/machines/Machines-1",
			"Connection": "/api/machines/Machines-1/connection"
		}
	}
`

func verifyGetMachineTestResponse(test *testing.T, machine *Machine) {
	expect := expect(test)

	expect.NotNil("Machine", machine)
	expect.EqualsString("Machine.ID", "Machines-1", machine.ID)
	expect.EqualsString("Machine.URI", "https://10.110.21.15:10933/", machine.URI)
	expect.EqualsString("Machine.Thumbprint", "B3092FEA722388E326CFF4D2F6E124B5727AEEDC", machine.Thumbprint)
}
