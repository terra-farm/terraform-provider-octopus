package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"octopus"
)

const (
	datasourceKeyMachineSlug       = "slug"
	datasourceKeyMachineName       = "name"
	datasourceKeyMachineURI        = "uri"
	datasourceKeyMachineThumbprint = "thumbprint"
)

func datasourceMachine() *schema.Resource {
	return &schema.Resource{
		Read:   datasourceMachineRead,
		Exists: datasourceMachineExists,

		Schema: map[string]*schema.Schema{
			datasourceKeyMachineSlug: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The machine slug (last segment of the machine URL in Octopus UI).",
			},
			datasourceKeyMachineName: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The machine name.",
			},
			datasourceKeyMachineURI: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The machine URI.",
			},
			datasourceKeyMachineThumbprint: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The machine thumbprint.",
			},
		},
	}
}

// Read a machine data-source.
func datasourceMachineRead(data *schema.ResourceData, provider interface{}) error {
	slug := data.Get(datasourceKeyMachineSlug).(string)

	log.Printf("Read machine '%s'.", slug)

	client := provider.(*octopus.Client)
	machine, err := client.GetMachine(slug)
	if err != nil {
		return err
	}

	if machine == nil {
		// Machine has been deleted.
		data.SetId("")

		return nil
	}

	data.SetId(machine.ID)
	data.Set(datasourceKeyMachineName, machine.Name)
	data.Set(datasourceKeyMachineURI, machine.URI)
	data.Set(datasourceKeyMachineThumbprint, machine.Thumbprint)

	return nil
}

// Determine whether a machine datasource exists.
func datasourceMachineExists(data *schema.ResourceData, provider interface{}) (exists bool, err error) {
	slug := data.Get(datasourceKeyMachineSlug).(string)

	log.Printf("Check if machine '%s' exists.", slug)

	client := provider.(*octopus.Client)

	var machine *octopus.Machine
	machine, err = client.GetMachine(slug)
	exists = machine != nil

	return
}
