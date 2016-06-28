package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"octopus"
)

const (
	datasourceKeyVariableProjectID = "project"
	datasourceKeyVariableName      = "name"
	datasourceKeyVariableValue     = "value"
	datasourceKeyVariableVariables = "variables"
	datasourceKeyVariableRoles     = "roles"
	datasourceKeyVariableMachines  = "machines"
	datasourceKeyVariableActions   = "actions"
)

func datasourceVariable() *schema.Resource {
	return &schema.Resource{
		Read:   datasourceVariableRead,
		Exists: datasourceVariableExists,

		Schema: map[string]*schema.Schema{
			datasourceKeyVariableProjectID: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The variable name.",
			},
			datasourceKeyVariableName: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			datasourceKeyVariableValue: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			datasourceKeyVariableVariables: &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			datasourceKeyVariableRoles: &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			datasourceKeyVariableMachines: &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			datasourceKeyVariableActions: &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

// Read a variable data-source.
func datasourceVariableRead(data *schema.ResourceData, provider interface{}) error {
	projectID := data.Get(resourceKeyVariableProjectID).(string)
	name := data.Get(resourceKeyVariableName).(string)

	propertyHelper := propertyHelper(data)
	targetScope := octopus.VariableScopes{
		Environments: propertyHelper.GetStringList(resourceKeyVariableEnvironments),
		Roles:        propertyHelper.GetStringList(resourceKeyVariableRoles),
		Machines:     propertyHelper.GetStringList(resourceKeyVariableMachines),
		Actions:      propertyHelper.GetStringList(resourceKeyVariableActions),
	}

	log.Printf("Read variable '%s' in project '%s', targeting scope %+v", name, projectID, targetScope)

	client := provider.(*octopus.Client)

	variableSet, err := client.GetProjectVariableSet(projectID)
	if err != nil {
		return fmt.Errorf("Error retrieving variable set for project '%s': %s", projectID, err.Error())
	}
	if variableSet == nil {
		return fmt.Errorf("Cannot find variable set for project '%s'", projectID)
	}

	matchingVariables := variableSet.GetVariablesByNameAndScopes(name, targetScope)
	if len(matchingVariables) == 0 {
		return fmt.Errorf("Cannot find variable '%s' in project '%s' with scope %+v", name, projectID, targetScope)
	} else if len(matchingVariables) != 1 {
		return fmt.Errorf("Multiple variables exactly match name '%s' in project '%s' with scope %+v", name, projectID, targetScope)
	}

	variable := matchingVariables[0]
	data.SetId(variable.ID)
	data.Set(datasourceKeyVariableName, variable.Name)

	return nil
}

// Determine whether a variable data-source exists.
func datasourceVariableExists(data *schema.ResourceData, provider interface{}) (exists bool, err error) {
	id := data.Id()
	projectID := data.Get(resourceKeyVariableProjectID).(string)
	name := data.Get(resourceKeyVariableName).(string)

	log.Printf("Check if variable '%s' (name = '%s') exists in project '%s'", id, name, projectID)

	client := provider.(*octopus.Client)

	variableSet, err := client.GetProjectVariableSet(projectID)
	if err != nil {
		err = fmt.Errorf("Error retrieving variable set for project '%s': %s", projectID, err.Error())

		return
	}
	if variableSet == nil {
		err = fmt.Errorf("Cannot find variable set for project '%s'", projectID)

		return
	}

	variable := variableSet.GetVariableByID(id)
	if variable == nil {
		log.Printf("Variable '%s' (name = '%s') not found", id, name)

		return false, nil
	}

	exists = true

	return
}
