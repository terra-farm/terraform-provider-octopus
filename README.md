# terraform-octopus
A plugin for Terraform to control / integrate with [Octopus Deploy](https://octopus.com/).

This is a work in progress. Currently, the following resource types are supported:

* `octopus_environment`: An Octopus Deploy environment
* `octopus_variable`: A Octopus Deploy variable (currently only project-level variables are supported)

Note that if a variable already exists with the specified name and scope, the provider will start managing the existing variable.

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

resource "octopus_environment" "my_environment" {
	name         = "MyEnvironment"
}

resource "octopus_variable" "my_variable" {
	# This is the Id (or slug) of the project in which the variable is defined.
	project      = "Projects-501"

	name         = "MyVariable"
	value        = "My Variable Value"

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
