package parser

import (
	"bufio"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// ParseKnownHosts parses the SSH known_hosts file and returns a list of hosts
func ParseKnownHosts() ([]SSHHost, error) {
	var hosts []SSHHost

	// Get the SSH known_hosts path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return hosts, err
	}

	knownHostsPath := filepath.Join(homeDir, ".ssh", "known_hosts")

	// Check if known_hosts file exists
	if _, err := os.Stat(knownHostsPath); os.IsNotExist(err) {
		return hosts, nil // Return empty slice, not an error
	}

	file, err := os.Open(knownHostsPath)
	if err != nil {
		return hosts, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	hostMap := make(map[string]bool) // To avoid duplicates

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Skip special SSH entries (certificate authority, revoked keys, etc.)
		if strings.HasPrefix(line, "@cert-authority") || strings.HasPrefix(line, "@revoked") {
			continue
		}

		// Parse known_hosts line format: host[,host2,...] key-type key [comment]
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		// First part contains the host(s)
		hostPart := parts[0]

		// Skip hashed hostnames (start with |1|)
		if strings.HasPrefix(hostPart, "|1|") {
			continue
		}

		// Handle multiple hosts separated by commas
		hostNames := strings.Split(hostPart, ",")

		for _, hostName := range hostNames {
			hostName = strings.TrimSpace(hostName)
			if hostName == "" {
				continue
			}

			// Extract hostname and port if present
			host, port := parseHostPort(hostName)

			// Skip if we've already seen this host
			if hostMap[host] {
				continue
			}

			// Skip IP addresses and localhost
			if net.ParseIP(host) != nil || host == "localhost" || host == "127.0.0.1" || host == "::1" {
				continue
			}

			hostMap[host] = true

			sshHost := SSHHost{
				Name:     host,
				HostName: host,
				Port:     port,
				Source:   "known_hosts",
			}

			hosts = append(hosts, sshHost)
		}
	}

	return hosts, scanner.Err()
}

// parseHostPort extracts hostname and port from a known_hosts entry
// Handles formats like: hostname, [hostname]:port, hostname:port
func parseHostPort(hostEntry string) (string, string) {
	// Handle bracketed format [hostname]:port
	if strings.HasPrefix(hostEntry, "[") && strings.Contains(hostEntry, "]:") {
		endBracket := strings.Index(hostEntry, "]:")
		if endBracket > 1 {
			host := hostEntry[1:endBracket]
			port := hostEntry[endBracket+2:]
			return host, port
		}
	}

	// Handle regular hostname:port format (but be careful about IPv6)
	if strings.Contains(hostEntry, ":") && !strings.Contains(hostEntry, "::") {
		parts := strings.Split(hostEntry, ":")
		if len(parts) == 2 {
			return parts[0], parts[1]
		}
	}

	// Just a hostname without port
	return hostEntry, ""
}
