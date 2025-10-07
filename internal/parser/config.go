package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// SSHHost represents a parsed SSH host with its configuration
type SSHHost struct {
	Name     string
	HostName string
	User     string
	Port     string
	Source   string // "config" or "known_hosts"
}

// ParseSSHConfig parses the SSH config file and returns a list of hosts
func ParseSSHConfig() ([]SSHHost, error) {
	var hosts []SSHHost

	// Get the SSH config path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return hosts, err
	}

	configPath := filepath.Join(homeDir, ".ssh", "config")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return hosts, nil // Return empty slice, not an error
	}

	file, err := os.Open(configPath)
	if err != nil {
		return hosts, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentHost SSHHost
	var inHostSection bool

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split the line into key and value
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		key := strings.ToLower(parts[0])
		value := strings.Join(parts[1:], " ")

		switch key {
		case "host":
			// Save previous host if it exists
			if inHostSection && currentHost.Name != "" {
				currentHost.Source = "config"
				hosts = append(hosts, currentHost)
			}

			// Start new host section
			currentHost = SSHHost{Name: value}
			inHostSection = true

			// Skip wildcard entries
			if strings.Contains(value, "*") || strings.Contains(value, "?") {
				inHostSection = false
				continue
			}

		case "hostname":
			if inHostSection {
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
		currentHost.Source = "config"
		hosts = append(hosts, currentHost)
	}

	return hosts, scanner.Err()
}
