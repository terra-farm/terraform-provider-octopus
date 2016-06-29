package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"octopus"
)

const (
	resourceKeyEnvironmentName          = "name"
	resourceKeyEnvironmentDescription   = "description"
	resourceKeyEnvironmentProjectGroups = "project_groups"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnvironmentCreate,
		Read:   resourceEnvironmentRead,
		Update: resourceEnvironmentUpdate,
		Delete: resourceEnvironmentDelete,
		Exists: resourceEnvironmentExists,

		Schema: map[string]*schema.Schema{
			resourceKeyEnvironmentName: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The environment name.",
			},
			resourceKeyEnvironmentDescription: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The environment description.",
			},
			resourceKeyEnvironmentProjectGroups: &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Optional:    true,
				Default:     nil,
				Description: "The Ids of project groups associated with the environment.",
			},
		},
	}
}

// Create an environment resource.
func resourceEnvironmentCreate(data *schema.ResourceData, provider interface{}) error {
	name := data.Get(resourceKeyEnvironmentName).(string)
	description := data.Get(resourceKeyEnvironmentDescription).(string)

	log.Printf("Create environment named '%s'.", name)

	client := provider.(*octopus.Client)

	environment, err := client.CreateEnvironment(name, description, 0)
	if err != nil {
		return err
	}

	data.SetId(environment.ID)

	return nil
}

// Read an environment resource.
func resourceEnvironmentRead(data *schema.ResourceData, provider interface{}) error {
	id := data.Id()
	name := data.Get(resourceKeyEnvironmentName).(string)

	log.Printf("Read environment '%s' (name = '%s').", id, name)

	client := provider.(*octopus.Client)
	environment, err := client.GetEnvironment(id)
	if err != nil {
		return err
	}

	if environment == nil {
		// Environment has been deleted.
		data.SetId("")

		return nil
	}

	data.Set(resourceKeyEnvironmentName, environment.Name)
	data.Set(resourceKeyEnvironmentDescription, environment.Description)

	return nil
}

// Update an environment resource.
func resourceEnvironmentUpdate(data *schema.ResourceData, provider interface{}) error {
	id := data.Id()

	log.Printf("Update environment '%s'.", id)

	if !(data.HasChange(resourceKeyEnvironmentName) || data.HasChange(resourceKeyEnvironmentDescription)) {
		return nil // Nothing to do.
	}

	client := provider.(*octopus.Client)
	environment, err := client.GetEnvironment(id)
	if err != nil {
		return err
	}
	if environment != nil {
		// Environment has been deleted.
		data.SetId("")

		return nil
	}

	if data.HasChange(resourceKeyEnvironmentName) {
		environment.Name = data.Get(resourceKeyEnvironmentName).(string)
	}

	if data.HasChange(resourceKeyEnvironmentDescription) {
		environment.Description = data.Get(resourceKeyEnvironmentDescription).(string)
	}

	_, err = client.UpdateEnvironment(environment)

	return err
}

// Delete an environment resource.
func resourceEnvironmentDelete(data *schema.ResourceData, provider interface{}) error {
	id := data.Id()
	name := data.Get(resourceKeyEnvironmentName).(string)

	log.Printf("Delete Environment '%s' (name = '%s').", id, name)

	client := provider.(*octopus.Client)

	return client.DeleteEnvironment(id)
}

// Determine whether an environment resource exists.
func resourceEnvironmentExists(data *schema.ResourceData, provider interface{}) (exists bool, err error) {
	id := data.Id()

	log.Printf("Check if environment '%s' exists.", id)

	client := provider.(*octopus.Client)

	var environment *octopus.Environment
	environment, err = client.GetEnvironment(id)
	exists = environment != nil

	return
}
