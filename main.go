package main

import (
	"github.com/hashicorp/terraform/plugin"
)

// The main program entry-point.
func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc:    Provider,
		ProvisionerFunc: Provisioner,
	})
}
