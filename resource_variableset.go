package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"octopus"
)

const (
	resourceKeyVariableSetProjectID = "project"
)

func resourceVariableSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceVariableSetCreate,
		Read:   resourceVariableSetRead,
		Update: resourceVariableSetUpdate,
		Delete: resourceVariableSetDelete,

		Schema: map[string]*schema.Schema{
			resourceKeyVariableSetProjectID: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The variable set name.",
			},
		},
	}
}

// Create a variable set resource.
func resourceVariableSetCreate(data *schema.ResourceData, provider interface{}) error {
	projectID := data.Get(resourceKeyVariableSetProjectID).(string)

	log.Printf("Create variable set for project '%s'.", projectID)

	providerClient := provider.(*octopus.Client)
	providerClient.Reset() // TODO: Replace call to Reset with appropriate API call(s).

	data.SetId(projectID) // TODO: Use variable set Id when we actually create one.

	return nil
}

// Read a variable set resource.
func resourceVariableSetRead(data *schema.ResourceData, provider interface{}) error {
	id := data.Id()
	projectID := data.Get(resourceKeyVariableSetProjectID).(string)

	log.Printf("Read variable set '%s' (for project '%s').", id, projectID)

	providerClient := provider.(*octopus.Client)
	providerClient.Reset() // TODO: Replace call to Reset with appropriate API call(s).

	return nil
}

// Update a variable set resource.
func resourceVariableSetUpdate(data *schema.ResourceData, provider interface{}) error {
	id := data.Id()
	projectID := data.Get(resourceKeyVariableSetProjectID).(string)

	log.Printf("Update variable set '%s' (for project '%s').", id, projectID)

	providerClient := provider.(*octopus.Client)
	providerClient.Reset() // TODO: Replace call to Reset with appropriate API call(s).

	return nil
}

// Delete a variable set resource.
func resourceVariableSetDelete(data *schema.ResourceData, provider interface{}) error {
	id := data.Id()
	projectID := data.Get(resourceKeyVariableSetProjectID).(string)

	log.Printf("Delete variable set '%s' (for project '%s').", id, projectID)

	providerClient := provider.(*octopus.Client)
	providerClient.Reset() // TODO: Replace call to Reset with appropriate API call(s).

	return nil
}
