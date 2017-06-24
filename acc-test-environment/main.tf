## Provider configuration

# The Id of the target Azure subscription.
variable "azure_subscription_id"    { sensitive = true }

# The client Id used to authenticate to Azure.
variable "azure_client_id"          { sensitive = true }

# The client secret used to authenticate to Azure.
variable "azure_client_secret"      { sensitive = true }

# The Id of target Azure AD tenant.
variable "azure_tenant_id"          { sensitive = true }

provider "azurerm" {
  subscription_id = "${azure_subscription_id}"
  client_id       = "${azure_client_id}"
  client_secret   = "${azure_client_secret}"
  tenant_id       = "${azure_tenant_id}"
}

## Common configuration

# The name of the target Azure region (i.e. datacenter).
variable "region_name"              { default = "West Central US" }

# The name of the resource group that holds the Octopus server used by acceptance tests.
variable "resource_group_name"      { default = "terraform-provider-octopus-acctest" }

# Used to prevent naming clashes between multiple concurrent deployments.
variable "uniqueness_key" { default = "acctest" }

# TODO: Define other variables
