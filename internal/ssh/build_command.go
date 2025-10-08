package ssh

import (
	"ssh-tui/internal/parser"
	"strings"
)

// BuildSSHCommand constructs the SSH command string for callers
func BuildSSHCommand(host *parser.SSHHost, options string) string {
	var parts []string

	parts = append(parts, "ssh")

	// If the host is from SSH config, just use the host name directly
	// SSH will handle the configuration expansion automatically
	if host.Source == "config" {
		parts = append(parts, host.Name)
		if options != "" {
			parts = append(parts, options)
		}
		return strings.Join(parts, " ")
	}

	// For hosts from known_hosts or other sources, expand the configuration
	// Add host-specific options from SSH config
	if host.Port != "" && host.Port != "22" {
		parts = append(parts, "-p", host.Port)
	}

	// Construct the connection string
	var target string
	if host.User != "" {
		target = host.User + "@"
	}

	// Use HostName if available, otherwise use Name
	if host.HostName != "" && host.HostName != host.Name {
		target += host.HostName
	} else {
		target += host.Name
	}

	parts = append(parts, target)

	// Add user-provided options after the host/target
	if options != "" {
		parts = append(parts, options)
	}

	return strings.Join(parts, " ")
}
