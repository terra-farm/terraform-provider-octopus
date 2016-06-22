package main

import (
	"fmt"
	"github.com/hashicorp/terraform/communicator/winrm"
	"github.com/hashicorp/terraform/helper/config"
	"github.com/hashicorp/terraform/terraform"
	"log"
)

// Provisioner creates the Octopus Deploy resource provisioner.
func Provisioner() terraform.ResourceProvisioner {
	return &OctopusProvisioner{}
}

// OctopusProvisioner represents the Octopus Deploy provisioner.
type OctopusProvisioner struct {
	// TODO: Decide what state we need to hold.
}

// Apply executes provisioner
func (provisioner *OctopusProvisioner) Apply(output terraform.UIOutput, state *terraform.InstanceState, cfg *terraform.ResourceConfig) error {
	log.Print("Executing Octopus Deploy provisioner.")

	err := ensureConnectionIsWinRM(state)
	if err != nil {
		return err
	}

	communicator, err := winrm.New(state)
	if err != nil {
		return err
	}

	communicator.Connect(output)
	defer communicator.Disconnect()

	// TODO: Use communicator to upload and execute provisioning script.

	return nil
}

// Validate checks if the required arguments are configured for the provisioner.
func (provisioner *OctopusProvisioner) Validate(cfg *terraform.ResourceConfig) (warnings []string, errors []error) {
	validator := config.Validator{
		Required: []string{
		// TODO: List required fields.
		},
		Optional: []string{
		// TODO: List optional fields.
		},
	}

	return validator.Validate(cfg)
}

func ensureConnectionIsWinRM(state *terraform.InstanceState) error {
	connectionType := state.Ephemeral.ConnInfo["type"]

	switch connectionType {
	case "winrm":
		return nil
	default:
		return fmt.Errorf("Connection type '%s' not supported for provisioning Octopus Deploy", connectionType)
	}
}
