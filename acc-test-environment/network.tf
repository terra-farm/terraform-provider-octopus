# The primary network for the acceptance test environment.
resource "azurerm_virtual_network" "primary" {
	name 				= "tf-octo-acc-test-${var.uniqueness_key}-network"
	address_space 		= ["10.7.0.0/16"]
	location 			= "${var.region_name}"
	resource_group_name = "${var.resource_group_name}"
}

# The primary subnet for the acceptance test environment.
resource "azurerm_subnet" "primary" {
	name 					= "tf-octo-acc-test-${var.uniqueness_key}-subnet"
	resource_group_name		= "${var.resource_group_name}"
	virtual_network_name    = "${azurerm_virtual_network.primary.name}"
	address_prefix 			= "10.7.1.0/24"
}

# The primary network adapter for the Octopus Server VM.
resource "azurerm_network_interface" "octo" {
	name 				= "octo-${var.uniqueness_key}-ni"
	location 			= "${var.region_name}"
	resource_group_name = "${var.resource_group_name}"

	ip_configuration {
		name 		= "octo-${var.uniqueness_key}-ni-config"
		subnet_id   = "${azurerm_subnet.primary.id}"

		# Hook up public IP to private IP.
		public_ip_address_id            = "${azurerm_public_ip.octo.id}"
		private_ip_address_allocation   = "dynamic"
	}
}

# The default network security group for the acceptance test environment.
resource "azurerm_network_security_group" "default" {
    name                = "octo-${var.uniqueness_key}-default-nsg"
    location            = "${var.region_name}"
    resource_group_name = "${var.resource_group_name}"

    # Remote Desktop
    security_rule {
        name                        = "rdp"
        priority                    = 100
        direction                   = "Inbound"
        access                      = "Allow"
        protocol                    = "Tcp"
        source_port_range           = "*"
        destination_port_range      = "3389"
        source_address_prefix       = "*"
        destination_address_prefix  = "*"
    }

    # WinRM
    security_rule {
        name                        = "winrm"
        priority                    = 101
        direction                   = "Inbound"
        access                      = "Allow"
        protocol                    = "Tcp"
        source_port_range           = "*"
        destination_port_range      = "5985" # HTTP
        source_address_prefix       = "*"
        destination_address_prefix  = "*"
    }
}

# Public IP address for access to the Octopus Server VM.
resource "azurerm_public_ip" "octo" {
	name 				= "tf-octo-acc-test-${var.uniqueness_key}-pip"
	location 			= "${var.region_name}"
	resource_group_name = "${var.resource_group_name}"

	public_ip_address_allocation = "static"
}
