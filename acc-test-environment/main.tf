## Provider configuration

# The Id of the target Azure subscription.
variable "azure_subscription_id"    { }

# The client Id used to authenticate to Azure.
variable "azure_client_id"          { }

# The client secret used to authenticate to Azure.
variable "azure_client_secret"      { }

# The Id of target Azure AD tenant.
variable "azure_tenant_id"          { }

provider "azurerm" {
  subscription_id = "${var.azure_subscription_id}"
  client_id       = "${var.azure_client_id}"
  client_secret   = "${var.azure_client_secret}"
  tenant_id       = "${var.azure_tenant_id}"
}

## Common configuration

# The name of the target Azure region (i.e. datacenter).
variable "region_name"              { default = "West US" }

# The name of the resource group that holds the Octopus server used by acceptance tests.
variable "resource_group_name"      { default = "terraform-provider-octopus-acctest" }

# The name of the storage account where VM disks (etc) are located.
variable "storage_account_name"     { default = "tfprovideroctopusacctest" }

# Used to prevent naming clashes between multiple concurrent deployments.
variable "uniqueness_key"           { default = "acctest" }

# The instance type for the Octopus Server VM.
variable "octo_vm_instance_type"    { default = "Standard_A3" }

# The administrator username for the Octopus and SQL servers.
variable "admin_username"           { default = "octo-admin" }

# The administrator password for the Octopus and SQL servers.
variable "admin_password"   { }
