 <h1>WireGuard Configuration Documentation</h1>
    <p>This project provides a utility for managing WireGuard VPN configurations. The code includes functionality to generate server and client keys, manage clients, and start or stop the WireGuard service.</p>

    <h2>PeerConfig Structure</h2>
    <p>The <code>PeerConfig</code> structure defines the configuration of a WireGuard peer:</p>
    <ul>
        <li><strong>PublicKey:</strong> The peer's public key.</li>
        <li><strong>AllowedIPs:</strong> The IP addresses allowed for the peer.</li>
        <li><strong>Endpoint:</strong> The peer's endpoint (IP and port).</li>
    </ul>

    <h2>Client Structure</h2>
    <p>The <code>Client</code> structure defines the configuration of a WireGuard client:</p>
    <ul>
        <li><strong>Id:</strong> The client's ID.</li>
        <li><strong>Status:</strong> The client's status (active/inactive).</li>
        <li><strong>AddressClient:</strong> The IP address assigned to the client.</li>
        <li><strong>PubkeyPath:</strong> The path to the public key file for the client.</li>
        <li><strong>PrivkeyPath:</strong> The path to the private key file for the client.</li>
        <li><strong>PrivateClientKey:</strong> The private key for the client.</li>
        <li><strong>PublicClientKey:</strong> The public key for the client.</li>
        <li><strong>Peer:</strong> The peer configuration for the client.</li>
        <li><strong>PeerStr:</strong> The peer configuration string to be appended to the WireGuard configuration.</li>
        <li><strong>Config:</strong> The WireGuard configuration for the client.</li>
        <li><strong>TgId:</strong> The client's Telegram ID for bot communication.</li>
    </ul>

    <h2>WireGuardConfig Structure</h2>
    <p>The <code>WireGuardConfig</code> structure defines the server configuration:</p>
    <ul>
        <li><strong>PrivateKey:</strong> The private key of the WireGuard server.</li>
        <li><strong>PublicKey:</strong> The public key of the WireGuard server.</li>
        <li><strong>Endpoint:</strong> The server's endpoint (IP and port).</li>
        <li><strong>ListenPort:</strong> The server's listening port.</li>
        <li><strong>InterName:</strong> The network interface name.</li>
        <li><strong>BotToken:</strong> The Telegram bot token for communication.</li>
        <li><strong>Clients:</strong> A map of clients managed by the server.</li>
    </ul>

    <h2>Client Management Methods</h2>
    <p>The code provides several methods for managing WireGuard clients:</p>
    <ul>
        <li><strong>StopClient(id int):</strong> Stops a client and removes its configuration from the WireGuard configuration file.</li>
        <li><strong>ActClient(id int):</strong> Activates a client by appending its configuration to the WireGuard configuration file.</li>
        <li><strong>DeleteClient(id int):</strong> Deletes a client by stopping it and removing it from the client map.</li>
        <li><strong>AllClients() string:</strong> Returns the status of all clients as a formatted string.</li>
    </ul>

    <h2>Server Management Methods</h2>
    <p>The <code>WireGuardConfig</code> structure includes several methods for managing the WireGuard server:</p>
    <ul>
        <li><strong>Autostart():</strong> Initializes and starts the WireGuard service, generating keys and configuration files.</li>
        <li><strong>GenServerKeys():</strong> Generates the server's private and public keys and saves them to files.</li>
        <li><strong>RandomPort():</strong> Randomly selects a port for the WireGuard server to listen on.</li>
        <li><strong>GetIPAndInterfaceName() error:</strong> Retrieves the server's IP address and network interface name.</li>
        <li><strong>GenerateWireGuardConfig():</strong> Generates the WireGuard configuration file for the server.</li>
        <li><strong>WireguardStart():</strong> Starts the WireGuard service and enables port forwarding and UFW rules.</li>
    </ul>

    <h2>Utility Functions</h2>
    <p>Additional utility functions include:</p>
    <ul>
        <li><strong>restWireguard():</strong> Restarts the WireGuard service after any configuration changes.</li>
        <li><strong>isWiredInterface(name string) bool:</strong> Checks if a network interface is wired (Ethernet).</li>
        <li><strong>isWirelessInterface(name string) bool:</strong> Checks if a network interface is wireless (Wi-Fi).</li>
    </ul>

    <h2>Adding a Client</h2>
    <p>The <code>AddWireguardClient(clientID int)</code> function adds a new WireGuard client, generates keys for the client, and appends the client's configuration to the WireGuard configuration file.</p>
