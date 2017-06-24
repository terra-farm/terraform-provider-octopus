# Public IP address for access to the target VM.
resource "azurerm_public_ip" "primary" {
	name 					= "tf-octo-acc-test-${var.uniqueness_key}-pip"
	location 				= "${var.region_name}"
	resource_group_name 	= "${var.resource_group_name}"

	public_ip_address_allocation = "static"
}

# The primary network for the target VM.
resource "azurerm_virtual_network" "primary" {
	name 						= "tf-octo-acc-test-${var.uniqueness_key}-network"
	address_space 				= ["10.7.0.0/16"]
	location 					= "${var.region_name}"
	resource_group_name			= "${var.resource_group_name}"
}

# The primary subnet for the target VM.
resource "azurerm_subnet" "primary" {
	name 						= "tf-octo-acc-test-${var.uniqueness_key}-subnet"
	resource_group_name			= "${var.resource_group_name}"
	virtual_network_name 		= "${azurerm_virtual_network.primary.name}"
	address_prefix 				= "10.7.1.0/24"
}

# The primary network adapter for the target VM.
resource "azurerm_network_interface" "primary" {
	name 						= "octo-${var.uniqueness_key}-ni"
	location 					= "${var.region_name}"
	resource_group_name 		= "${var.resource_group_name}"

	ip_configuration {
		name 					= "octo-${var.uniqueness_key}-ni-config"
		subnet_id 				= "${azurerm_subnet.primary.id}"

		# Hook up public IP to private IP.
		public_ip_address_id	= "${azurerm_public_ip.primary.id}"
		
		private_ip_address_allocation = "dynamic"
	}
}
