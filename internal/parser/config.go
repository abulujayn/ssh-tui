package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"ssh-tui/internal/types"
)

// ParseSSHConfig parses the SSH config file and returns a list of hosts
func ParseSSHConfig() ([]types.SSHHost, error) {
	var hosts []types.SSHHost

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return hosts, err
	}
	configPath := filepath.Join(homeDir, ".ssh", "config")

	// Check if config file exists; return empty if not (config is optional)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return hosts, nil
	}
	file, err := os.Open(configPath)
	if err != nil {
		return hosts, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentHost types.SSHHost
	var inHostSection bool

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		key := strings.ToLower(parts[0])
		value := strings.Join(parts[1:], " ")

		switch key {
		case "host":
			if inHostSection && currentHost.Name != "" {
				currentHost.Source = types.SourceConfig
				hosts = append(hosts, currentHost)
			}

			// Parse multiple hostnames (space-delimited)
			hostNames := strings.Fields(value)
			if len(hostNames) == 0 {
				inHostSection = false
				continue
			}

			// Use the first hostname as the primary name
			primaryHostName := hostNames[0]

			// Skip wildcard entries
			if strings.Contains(primaryHostName, "*") || strings.Contains(primaryHostName, "?") {
				inHostSection = false
				continue
			}

			// Collect aliases (all hostnames except the first one)
			var aliases []string
			if len(hostNames) > 1 {
				aliases = hostNames[1:]
			}

			currentHost = types.SSHHost{Name: primaryHostName, Aliases: aliases}
			inHostSection = true

		case "hostname":
			if inHostSection && IsValidHost(value) {
				currentHost.HostName = value
			}
		case "user":
			if inHostSection {
				currentHost.User = value
			}
		case "port":
			if inHostSection {
				currentHost.Port = value
			}
		}
	}

	// Add the last host if it exists
	if inHostSection && currentHost.Name != "" {
		currentHost.Source = types.SourceConfig
		hosts = append(hosts, currentHost)
	}

	return hosts, scanner.Err()
}
