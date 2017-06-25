## Database for Octopus Deploy

# SQL Server for the acceptance test environment
resource "azurerm_sql_server" "primary" {
    name                = "tfoctoacctest${var.uniqueness_key}"
    resource_group_name = "${var.resource_group_name}"
    location            = "${var.region_name}"
    version             = "12.0"

    administrator_login          = "${var.admin_username}"
    administrator_login_password = "${var.admin_password}"
}

# Database used by the environment's Octopus server
resource "azurerm_sql_database" "octo" {
    name                = "tfoctoacctest${var.uniqueness_key}"
    resource_group_name = "${var.resource_group_name}"
    location            = "${var.region_name}"

    edition     = "Basic"
    server_name = "${azurerm_sql_server.primary.name}"
}

# Permit T-SQL access from the Octopus server to the SQL server.
resource "azurerm_sql_firewall_rule" "octo_server" {
    name                = "octo_server_${var.uniqueness_key}"
    resource_group_name = "${var.resource_group_name}"
    
    server_name         = "${azurerm_sql_server.primary.name}"
    start_ip_address    = "${azurerm_network_interface.octo.private_ip_address}"
    end_ip_address      = "${azurerm_network_interface.octo.private_ip_address}"
}
