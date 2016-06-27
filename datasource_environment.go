package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"octopus"
)

const (
	datasourceKeyEnvironmentSlug        = "slug"
	datasourceKeyEnvironmentName        = "name"
	datasourceKeyEnvironmentDescription = "description"
)

func datasourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Read:   datasourceEnvironmentRead,
		Exists: datasourceEnvironmentExists,

		Schema: map[string]*schema.Schema{
			datasourceKeyEnvironmentSlug: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The environment slug (last segment of the environment URL).",
			},
			datasourceKeyEnvironmentName: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The environment name.",
			},
		},
	}
}

// Read a environment datasource.
func datasourceEnvironmentRead(data *schema.ResourceData, provider interface{}) error {
	slug := data.Get(datasourceKeyEnvironmentSlug).(string)

	log.Printf("Read environment '%s'.", slug)

	client := provider.(*octopus.Client)
	environment, err := client.GetEnvironment(slug)
	if err != nil {
		return err
	}

	if environment == nil {
		// Environment has been deleted.
		data.SetId("")

		return nil
	}

	data.SetId(environment.ID)
	data.Set(datasourceKeyEnvironmentName, environment.Name)
	data.Set(datasourceKeyEnvironmentDescription, environment.Description)

	return nil
}

// Determine whether a environment datasource exists.
func datasourceEnvironmentExists(data *schema.ResourceData, provider interface{}) (exists bool, err error) {
	id := data.Id()

	log.Printf("Check if environment '%s' exists.", id)

	client := provider.(*octopus.Client)

	var environment *octopus.Environment
	environment, err = client.GetEnvironment(id)
	exists = environment != nil

	return
}
