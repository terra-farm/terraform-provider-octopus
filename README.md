# terraform-octopus
A plugin for Terraform to control / integrate with [Octopus Deploy](https://octopus.com/).

This is a work in progress. More providers and data-sources are planned, as well as a [Provisioner](https://www.terraform.io/docs/provisioners/index.html) to install the Octopus tentacle.

Tested against Octopus Deploy v3.3.17.

The following resource types are currently supported:

* `octopus_environment`: Creates and manages an Octopus Deploy environment
* `octopus_variable`: Creates and manages an Octopus Deploy variable (currently only project-level variables are supported)

Note that variables are matched on both name and combined scopes (Environments, Roles, Machines, Actions). If a variable already exists with the specified name and scopes, the provider will start managing the existing variable.

The following data-source types are currently supported:
* `octopus_environment`: Tracks an existing Octopus Deploy environment
* `octopus_project`: Tracks an existing Octopus Deploy project
* `octopus_variable`: Tracks an existing Octopus Deploy variable (currently only project-level variables are supported)

Data-sources are similar to variables, except they are read-only. The provider will read and track their state but never modify it.

To get started:

* On windows, create / update `$HOME\terraform.rc`
* On Linux / OSX, create / update `~/.terraformrc`

And add the following contents:

```hcl
providers {
	octopus = "path-to-the-folder/containing/terraform-provider-octopus"
}
```

Create a folder containing a single `.tf` file:

```hcl
#
# This configuration will create an Octopus environment called "MyEnvironment" and configure a project-level variable named "MyVariable" to be scoped to it.
#

provider "octopus" {
	server_url = "https://my-octopus-server/"
	api_key = "my-octopus-api-key"
}

# Projects are a data source - the provider can read from them but not create or manage them.
data "octopus_project" "my_project" {
	slug = "terraformtest" # The last segment of the URL in the browser when viewing the project home page.
}

resource "octopus_environment" "my_environment" {
	name         = "MyEnvironment"
}

resource "octopus_variable" "my_variable" {
	# This is the Id (or slug) of the project in which the variable is defined.
	project      = "${data.octopus_project.my_project.id}"

	name         = "MyVariable"
	value        = "Hello World"

	# The scopes (environment, role, machine, action) to which the variable applies.
	environments = ["${octopus_environment.my_environment.id}"]
}
```

1. Run `terraform plan -out tf.plan`.
2. Verify that everything looks ok.
3. Run `terraform apply tf.plan`
4. Have a look around and, when it's time to clean up...
5. Run `terraform plan -destroy -out tf.plan`
6. Verify that everything looks ok.
7. Run `terraform apply tf.plan`
