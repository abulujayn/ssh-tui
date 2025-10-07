package parser

import (
	"fmt"
	"sort"
	"strings"
)

// DiscoverHosts discovers all SSH hosts from both config and known_hosts files
func DiscoverHosts() ([]SSHHost, error) {
	var allHosts []SSHHost

	// Parse SSH config
	configHosts, err := ParseSSHConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH config: %w", err)
	}

	// Parse known_hosts
	knownHosts, err := ParseKnownHosts()
	if err != nil {
		return nil, fmt.Errorf("failed to parse known_hosts: %w", err)
	}

	// Merge hosts with deduplication
	hostMap := make(map[string]SSHHost)

	// Add config hosts first (they have priority)
	for _, host := range configHosts {
		key := strings.ToLower(host.Name)
		hostMap[key] = host
	}

	// Add known_hosts entries that aren't already in config
	for _, host := range knownHosts {
		key := strings.ToLower(host.Name)
		if _, exists := hostMap[key]; !exists {
			hostMap[key] = host
		}
	}

	// Convert map back to slice
	for _, host := range hostMap {
		allHosts = append(allHosts, host)
	}

	// Sort hosts alphabetically by name
	sort.Slice(allHosts, func(i, j int) bool {
		return strings.ToLower(allHosts[i].Name) < strings.ToLower(allHosts[j].Name)
	})

	return allHosts, nil
}

// FilterHosts filters hosts by a search term
func FilterHosts(hosts []SSHHost, searchTerm string) []SSHHost {
	if searchTerm == "" {
		return hosts
	}

	var filtered []SSHHost
	searchLower := strings.ToLower(searchTerm)

	for _, host := range hosts {
		// Search in host name and hostname
		if strings.Contains(strings.ToLower(host.Name), searchLower) ||
			strings.Contains(strings.ToLower(host.HostName), searchLower) {
			filtered = append(filtered, host)
		}
	}

	return filtered
}

// FormatHostDisplay formats a host for display in the TUI
func FormatHostDisplay(host SSHHost) string {
	var parts []string

	// Add the host name
	parts = append(parts, host.Name)

	// Add hostname if different from name
	if host.HostName != "" && host.HostName != host.Name {
		parts = append(parts, fmt.Sprintf("(%s)", host.HostName))
	}

	// Add user if specified
	if host.User != "" {
		parts = append(parts, fmt.Sprintf("user: %s", host.User))
	}

	// Add port if specified and not default
	if host.Port != "" && host.Port != "22" {
		parts = append(parts, fmt.Sprintf("port: %s", host.Port))
	}

	// Add source indicator
	if host.Source == "config" {
		parts = append(parts, "[config]")
	} else {
		parts = append(parts, "[known_hosts]")
	}

	return strings.Join(parts, " ")
}
