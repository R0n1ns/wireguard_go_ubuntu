###WireGuard Management Utility Documentation
This project is a utility for managing WireGuard server and client configurations programmatically. It provides various functionalities such as adding, activating, stopping, and deleting clients. The configuration data for WireGuard is manipulated through Go, and it interacts with the underlying system via the command line.

###Features
Start, stop, and delete WireGuard clients dynamically.
Generate server keys and configure WireGuard.
Automatic handling of WireGuard restart and interface setup.
Add and manage clients via dynamically generated configurations.
Retrieve client statuses and other related information.

###Project Structure
PeerConfig: Structure representing the configuration of a peer.
Client: Represents a WireGuard client, including keys, status, and configuration.
WireGuardConfig: Manages the WireGuard server and clients, including server keys, interface details, and client handling.

###Methods
##Client Management
#StopClient(id int)
Stops the WireGuard client by removing its configuration from the WireGuard server configuration file.
#ActClient(id int)
Activates a WireGuard client by adding its configuration back to the WireGuard configuration file.
#DeleteClient(id int)
Stops and deletes a client from the system. The client is removed from memory and the configuration file.
#AllClients() string
Retrieves a formatted string listing all clients, their statuses, and their assigned addresses.
##WireGuard Server Management
#Autostart()
Runs the initialization process for the WireGuard server, including generating keys, assigning ports, and starting WireGuard.
#GenServerKeys()
Generates the server's private and public keys, stores them in variables and system files.
#RandomPort()
Assigns a random port for the WireGuard server.
#GetIPAndInterfaceName() error
Detects the active network interface and retrieves the server's IP address.
#GenerateWireGuardConfig()
Generates the configuration file for the WireGuard server based on the provided templates.
#WireguardStart()
Starts the WireGuard service and ensures proper network forwarding and firewall settings.
##Utility Functions
#restWireguard()
Restarts the WireGuard service to apply changes in configurations.
#AddWireguardClient(clientID int)
Adds a new WireGuard client to the configuration, generating client keys and updating the server configuration.

###How to Use
#Initialize WireGuard Configuration:
Call Autostart() to set up the server, generate keys, and configure the system.

#Add Clients:
Use AddWireguardClient(clientID int) to create new clients and append their configuration to the WireGuard server.

#Manage Clients:

To stop a client, call StopClient(id int).
To activate a client, use ActClient(id int).
To delete a client, use DeleteClient(id int).
View Client Status:
Call AllClients() to list all clients along with their statuses.

###Requirements
Go (Golang) installed
WireGuard installed and configured on the system (wg command available)
Administrative (root) permissions to modify system files (/etc/wireguard)
Installation
Install Go and WireGuard on your server.
Clone the repository and navigate to the project folder.
Configure system settings, ensuring that WireGuard is properly installed.
Run the application using Go commands:
