package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"octopus"
	"time"
)

const (
	resourceKeyEnvironmentName       = "name"
	resourceCreateTimeoutEnvironment = 30 * time.Minute
	resourceUpdateTimeoutEnvironment = 10 * time.Minute
	resourceDeleteTimeoutEnvironment = 15 * time.Minute
)

const computedPropertyDescription = "<computed>"

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnvironmentCreate,
		Read:   resourceEnvironmentRead,
		Update: resourceEnvironmentUpdate,
		Delete: resourceEnvironmentDelete,

		Schema: map[string]*schema.Schema{
			resourceKeyEnvironmentName: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The environment name.",
			},
		},
	}
}

// Create an environment resource.
func resourceEnvironmentCreate(data *schema.ResourceData, provider interface{}) error {
	name := data.Get(resourceKeyEnvironmentName).(string)

	log.Printf("Create environment named '%s'.", name)

	providerClient := provider.(*octopus.Client)
	providerClient.Reset() // TODO: Replace call to Reset with appropriate API call(s).

	data.SetId(name) // TODO: Use environment Id when we actually create one.

	return nil
}

// Read an environment resource.
func resourceEnvironmentRead(data *schema.ResourceData, provider interface{}) error {
	id := data.Id()
	name := data.Get(resourceKeyEnvironmentName).(string)

	log.Printf("Read environment '%s' (name = '%s'.", id, name)

	providerClient := provider.(*octopus.Client)
	providerClient.Reset() // TODO: Replace call to Reset with appropriate API call(s).

	return nil
}

// Update an environment resource.
func resourceEnvironmentUpdate(data *schema.ResourceData, provider interface{}) error {
	id := data.Id()

	log.Printf("Update environment '%s'.", id)

	providerClient := provider.(*octopus.Client)

	if data.HasChange(resourceKeyEnvironmentName) {
		old, new := data.GetChange(resourceKeyEnvironmentName)
		oldName := old.(string)
		newName := new.(string)

		log.Printf("Rename environment '%s' from '%s' to '%s'.", id, oldName, newName)
	}

	providerClient.Reset() // TODO: Replace call to Reset with appropriate API call(s).

	return nil
}

// Delete an environment resource.
func resourceEnvironmentDelete(data *schema.ResourceData, provider interface{}) error {
	id := data.Id()
	name := data.Get(resourceKeyEnvironmentName).(string)

	log.Printf("Delete Environment '%s' (name = '%s').", id, name)

	providerClient := provider.(*octopus.Client)
	providerClient.Reset() // TODO: Replace call to Reset with appropriate API call(s).

	return nil
}
