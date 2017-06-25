Param(
    [Parameter(Mandatory = $true)]
    [string] $SqlServerHost,

    [Parameter(Mandatory = $true)]
    [string] $Database,

    [Parameter(Mandatory = $true)]
    [string] $User,

    [Parameter(Mandatory = $true)]
    [string] $Password
)

Import-Module PowerShellGet
Install-Module OctopusDSC -Force
Import-Module OctopusDSC

Configuration Octopus {
    Import-DscResource -Module OctopusDSC

    Node "localhost" {
        cOctopusServer OctopusServer {
            Ensure = "present"
            State = "started"

            Name = "OctopusServer"

            WebListenPrefix = "http://localhost:9081"
            SqlDbConnectionString = "Server=$SqlServerHost;Database=$Database;UID=$User;PWD=$Password"

            OctopusAdminUsername = $User
            OctopusAdminPassword = $Password

            # AllowUpgradeCheck = $false
            # AllowCollectionOfAnonymousUsageStatistics = $false
            # ForceSSL = $false
            # ListenPort = 10943
        }
    }
}
Octopus

Start-DscConfiguration .\Octopus -Verbose -Wait
Test-DscConfiguration
