# Storage configuration

resource "azurerm_storage_container" "primary" {
	name 					= "tf-octo-acc-test-${var.uniqueness_key}"
	resource_group_name 	= "${var.resource_group_name}"
	storage_account_name 	= "${var.storage_acct_name}"
	container_access_type 	= "private"
}
