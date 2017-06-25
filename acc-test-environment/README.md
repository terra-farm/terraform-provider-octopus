# Acceptance test environment

This configuration creates (in Azure) the Octopus Deploy server used in Terraform acceptance tests.

You can create credentials for use in tests using the [Azure cross-platform CLI](https://www.terraform.io/docs/providers/azurerm/index.html#creating-credentials-using-the-azure-cli).

You'll need to supply a couple of values in `terraform.tfvars`:

* `azure_subscription_id` - The Id of your Azure subscription.
* `azure_client_id` - The client Id you created in the Azure CLI.
* `azure_client_secret` - The client secret you created in the Azure CLI.
* `azure_tenant_id` - The name of the Azure AD tenant that the subscription belongs to.

There are also a couple of option values in `main.tf` that you can override by supplying their values in `terraform.tfvars`:

* `region_name` - The name of the target region where you will be deploying the environment.  
**Note** - the VM and the storage account _must_ be the same location.
* `resource_group_name` - The name of the resource group where the environment will be deployed (must already exist; the configuration will not create it).
* `storage_account_name` - The name of the Azure Storage account where VM disks (etc) will be stored.
* `uniqueness_key` - A unique value added to resource names so that your deployment doesn't clash with other deployments.
