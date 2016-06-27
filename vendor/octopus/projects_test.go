package octopus

import (
	"testing"
)

/*
 * Integration tests
 */

// Get project by Id (successful).
func Test_Client_GetProject_Success(test *testing.T) {
	testClientRequest(test, &ClientTest{
		APIKey: "my-test-api-key",
		Invoke: func(test *testing.T, client *Client) {
			project, err := client.GetProject("Projects-105")
			if err != nil {
				test.Fatal(err)
			}

			verifyGetProjectTestResponse(test, project)
		},
		Handle: testRespondOK(getProjectTestResponse),
	})
}

/*
 * Test responses.
 */

const getProjectTestResponse = `
	{
		"Id": "Projects-501",
		"VariableSetId": "variableset-Projects-501",
		"DeploymentProcessId": "deploymentprocess-Projects-501",
		"IncludedLibraryVariableSetIds": [],
		"DefaultToSkipIfAlreadyInstalled": false,
		"VersioningStrategy": {
			"DonorPackageStepId": null,
			"Template": "#{Octopus.Version.LastMajor}.#{Octopus.Version.LastMinor}.#{Octopus.Version.NextPatch}"
		},
		"ReleaseCreationStrategy": {
			"ReleaseCreationPackageStepId": null,
			"ChannelId": null
		},
		"Name": "TerraformTest",
		"Slug": "terraformtest",
		"Description": "Adam's Terraform Test Project (please do not touch).",
		"IsDisabled": false,
		"ProjectGroupId": "ProjectGroups-49",
		"LifecycleId": "Lifecycles-1",
		"AutoCreateRelease": false,
		"Links": {
			"Self": "/api/projects/Projects-501",
			"Releases": "/api/projects/Projects-501/releases{/version}{?skip}",
			"Channels": "/api/projects/Projects-501/channels",
			"OrderChannels": "/api/projects/Projects-501/channels/order",
			"Variables": "/api/variables/variableset-Projects-501",
			"Progression": "/api/progression/Projects-501",
			"DeploymentProcess": "/api/deploymentprocesses/deploymentprocess-Projects-501",
			"Web": "/app#/projects/Projects-501",
			"Logo": "/api/projects/Projects-501/logo"
		}
	}
`

func verifyGetProjectTestResponse(test *testing.T, project *Project) {
	expect := expect(test)

	expect.NotNil("Project", project)
	expect.EqualsString("Project.ID", "Projects-501", project.ID)
	expect.EqualsString("Project.Name", "TerraformTest", project.Name)
	expect.EqualsString("Project.Slug", "terraformtest", project.Slug)
	expect.EqualsString("Project.Description", "Adam's Terraform Test Project (please do not touch).", project.Description)
	expect.EqualsString("Project.VariableSetID", "variableset-Projects-501", project.VariableSetID)
}
