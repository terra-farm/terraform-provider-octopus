package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"octopus"
)

const (
	datasourceKeyProjectSlug = "slug"
	datasourceKeyProjectName = "name"
)

func datasourceProject() *schema.Resource {
	return &schema.Resource{
		Read:   datasourceProjectRead,
		Exists: datasourceProjectExists,

		Schema: map[string]*schema.Schema{
			datasourceKeyProjectSlug: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The project slug (last segment of the project URL).",
			},
			datasourceKeyProjectName: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The project name.",
			},
		},
	}
}

// Read a project datasource.
func datasourceProjectRead(data *schema.ResourceData, provider interface{}) error {
	slug := data.Get(datasourceKeyProjectSlug).(string)

	log.Printf("Read project '%s'.", slug)

	client := provider.(*octopus.Client)
	project, err := client.GetProject(slug)
	if err != nil {
		return err
	}

	if project == nil {
		// Project has been deleted.
		data.SetId("")

		return nil
	}

	data.SetId(project.ID)
	data.Set(datasourceKeyProjectName, project.Name)

	return nil
}

// Determine whether a project datasource exists.
func datasourceProjectExists(data *schema.ResourceData, provider interface{}) (exists bool, err error) {
	id := data.Id()

	log.Printf("Check if project '%s' exists.", id)

	client := provider.(*octopus.Client)

	var project *octopus.Project
	project, err = client.GetProject(id)
	exists = project != nil

	return
}
