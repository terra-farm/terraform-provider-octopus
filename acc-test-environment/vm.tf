# The Octopus Deploy virtual machine.
resource "azurerm_virtual_machine" "octo" {
	name 					= "octo-${var.uniqueness_key}"
	location 				= "${var.region_name}"
	resource_group_name 	= "${var.resource_group_name}"
	network_interface_ids 	= [ "${azurerm_network_interface.primary.id}" ]

	vm_size 				= "${var.instance_type}"

	storage_image_reference {
		publisher 			= "Canonical"
		offer 				= "UbuntuServer"
		sku 				= "14.04.2-LTS"
		version 			= "latest"
	}

	storage_os_disk {
		name 				= "octo-${var.uniqueness_key}-osdisk1"
		vhd_uri 			= "https://${var.storage_acct_name}.blob.core.windows.net/${azurerm_storage_container.primary.name}/octo-${var.uniqueness_key}-osdisk1.vhd"
		caching 			= "ReadWrite"
		create_option 		= "FromImage"
	}

	os_profile {
		computer_name 		= "octo-${var.uniqueness_key}"
		admin_username 		= "octo-admin"
		admin_password		= "${var.initial_admin_password}"
	}

	# os_profile_linux_config {
	# 	disable_password_authentication = true
	# 	ssh_keys {
	#       path 				= "/home/ubuntu/.ssh/authorized_keys"
	#       key_data 			= "${var.ssh_key}"
	#     }
	# }

	tags {
		public_ip			= "${azurerm_public_ip.primary.ip_address}"
		private_ip			= "${azurerm_network_interface.primary.private_ip_address}"
	}