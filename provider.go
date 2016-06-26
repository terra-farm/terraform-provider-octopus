package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"octopus"
	"os"
)

// Provider creates the Octopus Deploy resource provider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		// Provider settings schema
		Schema: map[string]*schema.Schema{
			"server_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The base URL of the Octopus Deploy server.",
			},
			"api_key": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Default:   "",
				ConflictsWith: []string{
					"username",
					"password",
				},
				Description: "The API key used to authenticate to the Octopus Deploy API (if not specified, then the OCTOPUS_API_KEY environment variable will be used).",
			},
		},

		// Provider resource definitions
		ResourcesMap: map[string]*schema.Resource{
			"octopus_environment": resourceEnvironment(),
			"octopus_variable":    resourceVariable(),
		},

		// Provider configuration
		ConfigureFunc: configureProvider,
	}
}

// Configure the provider.
// Returns the provider's compute API client.
func configureProvider(providerSettings *schema.ResourceData) (interface{}, error) {
	server := providerSettings.Get("server").(string)
	apiKey := providerSettings.Get("api_key").(string)

	if isEmpty(apiKey) {
		apiKey = os.Getenv("OCTOPUS_API_KEY")
	}
	if isEmpty(apiKey) {
		return nil, fmt.Errorf("The 'api_key' property was not specified for the 'octopus' provider, and the 'OCTOPUS_API_KEY' environment variable is not present. Please supply either one of these to configure the API key used to authenticate to Octopus Deploy.")
	}

	var (
		client *octopus.Client
		err    error
	)
	client, err = octopus.NewClientWithAPIKey(server, apiKey)
	if err != nil {
		return nil, err
	}

	return client, nil
}
