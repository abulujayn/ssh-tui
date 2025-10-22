package ssh

import (
	"ssh-tui/internal/types"
	"strings"
)

// BuildSSHCommand constructs the SSH command string for callers
func BuildSSHCommand(host *types.SSHHost, options string) string {
	var parts []string

	parts = append(parts, "ssh")

	// If the host is from SSH config, just use the host name directly
	if host.Source == types.SourceConfig {
		parts = append(parts, host.Name)
		if options != "" {
			parts = append(parts, options)
		}
		return strings.Join(parts, " ")
	}

	// For hosts from known_hosts or other sources, expand the configuration
	if host.Port != "" && host.Port != types.DefaultSSHPort {
		parts = append(parts, "-p", host.Port)
	}

	var target string
	if host.User != "" {
		target = host.User + "@"
	}

	if host.HostName != "" && host.HostName != host.Name {
		target += host.HostName
	} else {
		target += host.Name
	}

	parts = append(parts, target)

	if options != "" {
		parts = append(parts, options)
	}

	return strings.Join(parts, " ")
}
