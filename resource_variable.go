package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"octopus"
)

const (
	resourceKeyVariableProjectID   = "project"
	resourceKeyVariableName        = "name"
	resourceKeyVariableEnvironment = "environment"
	resourceKeyVariableRole        = "role"
	resourceKeyVariableMachine     = "machine"
	resourceKeyVariableAction      = "action"
)

func resourceVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceVariableCreate,
		Read:   resourceVariableRead,
		Update: resourceVariableUpdate,
		Delete: resourceVariableDelete,

		Schema: map[string]*schema.Schema{
			resourceKeyVariableProjectID: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The variable name.",
			},
			resourceKeyVariableName: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			resourceKeyVariableEnvironment: &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			resourceKeyVariableRole: &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			resourceKeyVariableMachine: &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			resourceKeyVariableAction: &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

// Create a variable resource.
func resourceVariableCreate(data *schema.ResourceData, provider interface{}) error {
	propertyHelper := propertyHelper(data)

	projectID := data.Get(resourceKeyVariableProjectID).(string)
	name := data.Get(resourceKeyVariableName).(string)

	targetScope := octopus.VariableScopes{
		Environments: propertyHelper.GetStringList(resourceKeyVariableEnvironment),
		Roles:        propertyHelper.GetStringList(resourceKeyVariableRole),
		Machines:     propertyHelper.GetStringList(resourceKeyVariableMachine),
		Actions:      propertyHelper.GetStringList(resourceKeyVariableAction),
	}

	log.Printf("Create variable '%s' for project '%s' (must match scopes: %#v)...", name, projectID, targetScope)

	providerClient := provider.(*octopus.Client)

	variableSet, err := providerClient.GetProjectVariableSet(projectID)
	if err != nil {
		return err
	}
	if variableSet == nil {
		return fmt.Errorf("Cannot find variable for project '%s'.", projectID)
	}

	matchingVariables := variableSet.GetVariablesByNameAndScopes(name, targetScope)

	var variable octopus.Variable
	if len(matchingVariables) == 0 {
		log.Printf("Create variable '%s' for project '%s' with scope %#v...", name, projectID, targetScope)

		variableSet.Variables = append(variableSet.Variables, octopus.Variable{
			Name:  name,
			Scope: targetScope,
		})

		variableSet, err = providerClient.UpdateVariableSet(variableSet)
		if err != nil {
			return err
		}

		matchingVariables = variableSet.GetVariablesByNameAndScopes(name, targetScope)
		if len(matchingVariables) != 1 {
			return fmt.Errorf("Found %d matching variables named '%s' for scope %#v (after attempting to create this variable for that scope).", len(matchingVariables), name, targetScope)
		}
	} else if len(matchingVariables) == 1 {
		log.Printf("Variable '%s' already exists for project '%s' with scope %#v.", name, projectID, targetScope)
	} else {
		return fmt.Errorf("Multiple variables exactly match scope %#v for variable '%s'.", name, targetScope)
	}

	variable = matchingVariables[0]
	data.SetId(variable.ID)

	return nil
}

// Read a variable resource.
func resourceVariableRead(data *schema.ResourceData, provider interface{}) error {
	id := data.Id()
	projectID := data.Get(resourceKeyVariableProjectID).(string)

	log.Printf("Read variable '%s' (for project '%s').", id, projectID)

	providerClient := provider.(*octopus.Client)
	providerClient.Reset() // TODO: Replace call to Reset with appropriate API call(s).

	return nil
}

// Update a variable resource.
func resourceVariableUpdate(data *schema.ResourceData, provider interface{}) error {
	id := data.Id()
	projectID := data.Get(resourceKeyVariableProjectID).(string)

	log.Printf("Update variable '%s' (for project '%s').", id, projectID)

	providerClient := provider.(*octopus.Client)
	providerClient.Reset() // TODO: Replace call to Reset with appropriate API call(s).

	return nil
}

// Delete a variable resource.
func resourceVariableDelete(data *schema.ResourceData, provider interface{}) error {
	id := data.Id()
	projectID := data.Get(resourceKeyVariableProjectID).(string)

	log.Printf("Delete variable '%s' (for project '%s').", id, projectID)

	providerClient := provider.(*octopus.Client)
	providerClient.Reset() // TODO: Replace call to Reset with appropriate API call(s).

	return nil
}
