package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"log"
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
			"username": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "",
				ConflictsWith: []string{"api_key"},
				Description:   "The user name used to authenticate to the Octopus Deploy API (if not specified, then the OCTOPUS_USER environment variable will be used).",
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Default:   "",
				ConflictsWith: []string{
					"api_key",
				},
				Description: "The password used to authenticate to the Octopus Deploy API (if not specified, then the OCTOPUS_PASSWORD environment variable will be used).",
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
		},

		// Provider configuration
		ConfigureFunc: configureProvider,
	}
}

// Configure the provider.
// Returns the provider's compute API client.
func configureProvider(providerSettings *schema.ResourceData) (interface{}, error) {
	server := providerSettings.Get("server").(string)

	username := providerSettings.Get("username").(string)
	password := providerSettings.Get("password").(string)
	apiKey := providerSettings.Get("api_key").(string)

	if isEmpty(apiKey) {
		apiKey = os.Getenv("OCTOPUS_API_KEY")
	}

	if isEmpty(username) {
		apiKey = os.Getenv("OCTOPUS_USER")
	}

	if isEmpty(password) {
		apiKey = os.Getenv("OCTOPUS_PASSWORD")
	}

	var (
		client *octopus.Client
		err    error
	)
	if !isEmpty(username) {
		log.Printf("Octopus Deploy provider configured for HTTP Basic authentication (server = '%s', user = '%s').", server, username)

		client, err = octopus.NewClientWithBasicAuth(server, username, password)
		if err != nil {
			return nil, err
		}
	} else if !isEmpty(apiKey) {
		log.Printf("Octopus Deploy provider configured for API key authentication (server = '%s', API key = '%s').", server, apiKey)

		client, err = octopus.NewClientWithAPIKey(server, apiKey)
		if err != nil {
			return nil, err
		}
	}

	if client == nil {
		return nil, fmt.Errorf("Neither the 'username' nor the 'api_key' property was not specified for the 'octopus' provider, and the 'OCTOPUS_API_KEY' environment variable is not present. Please supply either one of these to configure the user name / API key used to authenticate to Octopus Deploy.")
	}

	return client, nil
}
