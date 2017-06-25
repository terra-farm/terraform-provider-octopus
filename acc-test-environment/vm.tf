# The Octopus Deploy virtual machine.
resource "azurerm_virtual_machine" "octo" {
	name 					= "octo-${var.uniqueness_key}"
	location 				= "${var.region_name}"
	resource_group_name 	= "${var.resource_group_name}"
	network_interface_ids 	= [ "${azurerm_network_interface.octo.id}" ]

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
		admin_username 		= "${var.admin_username}"
		admin_password		= "${var.admin_password}"
	}

	os_profile_windows_config {
		provision_vm_agent			= true
		enable_automatic_upgrades	= true

		winrm {
			protocol = "http"
		}
	}

	tags {
		public_ip			= "${azurerm_public_ip.octo.ip_address}"
		private_ip			= "${azurerm_network_interface.octo.private_ip_address}"
	}
}

# Configure WinRM to enable SSL (since it's otherwise unusable, given Azure's default config).
resource "azurerm_virtual_machine_extension" "configure_winrm" {
	name                 = "ConfigureWinRM"
	location             = "${var.region_name}"
	resource_group_name  = "${var.resource_group_name}"
	virtual_machine_name = "${azurerm_virtual_machine.octo.name}"
	publisher            = "Microsoft.Compute"
	type                 = "CustomScriptExtension"
	type_handler_version = "1.8"

	settings = <<SETTINGS
		{
			"fileUris": [
				"https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/201-vm-winrm-windows/ConfigureWinRM.ps1",
				"https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/201-vm-winrm-windows/makecert.exe",
				"https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/201-vm-winrm-windows/winrmconf.cmd"
			],
			"commandToExecute": "Powershell -ExecutionPolicy Unrestricted -file ConfigureWinRM.ps1 -HostName '${azurerm_public_ip.octo.fqdn}'"
		}
SETTINGS
}

# Install and configure the Octopus Deploy server.
resource "null_resource" "octo_server_install" {
	provisioner "file" {
		source		= "scripts/Install-OctopusServer.ps1"
		destination	= "C:\\Install-OctopusServer.ps1"

		connection {
			type 		= "winrm"
			host 		= "${azurerm_public_ip.octo.ip_address}"
			port		= 5986
			user 		= "${var.admin_username}"
			password 	= "${var.admin_password}"
			https		= true
			insecure	= true # We're using a self-signed certificate
		}
	}
	
	provisioner "remote-exec" {
		inline = [
			"C:\\Install-OctopusServer.ps1 -SqlServerHost '${azurerm_sql_server.primary.fully_qualified_domain_name}' -Database '${azurerm_sql_database.octo.name}' -User '${var.admin_username}' -Password '${var.admin_password}'"
		]
		
		connection {
			type 		= "winrm"
			host 		= "${azurerm_public_ip.octo.ip_address}"
			port		= 5986
			user 		= "${var.admin_username}"
			password 	= "${var.admin_password}"
			https		= true
			insecure	= true # We're using a self-signed certificate
		}
	}
	
	depends_on = [
		"azurerm_virtual_machine.octo",
		"azurerm_virtual_machine_extension.configure_winrm",
		"azurerm_public_ip.octo",
		"azurerm_network_security_group.default",
		"azurerm_sql_database.octo"
	]
}
