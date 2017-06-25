# The Octopus Deploy virtual machine.
resource "azurerm_virtual_machine" "octo" {
	name 					= "octo-${var.uniqueness_key}"
	location 				= "${var.region_name}"
	resource_group_name 	= "${var.resource_group_name}"
	network_interface_ids 	= [ "${azurerm_network_interface.primary.id}" ]

	vm_size 				= "${var.octo_vm_instance_type}"

	storage_image_reference {
		publisher = "MicrosoftWindowsServer"
		offer     = "WindowsServer"
		sku       = "2012-R2-Datacenter"
		version   = "latest"
	}

	storage_os_disk {
		name 				= "octo-${var.uniqueness_key}-osdisk1"
		vhd_uri 			= "https://${var.storage_account_name}.blob.core.windows.net/${azurerm_storage_container.primary.name}/octo-${var.uniqueness_key}-osdisk1.vhd"
		caching 			= "ReadWrite"
		create_option 		= "FromImage"
	}

	os_profile {
		computer_name 		= "octo-${var.uniqueness_key}"
		admin_username 		= "octo-admin"
		admin_password		= "${var.initial_admin_password}"
	}

	os_profile_windows_config {
		provision_vm_agent			= true
		enable_automatic_upgrades	= true

		winrm {
			protocol = "http"
		}
	}

	tags {
		public_ip			= "${azurerm_public_ip.primary.ip_address}"
		private_ip			= "${azurerm_network_interface.primary.private_ip_address}"
	}
}
