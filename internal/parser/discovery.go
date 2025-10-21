package parser

import (
	"fmt"
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

	// Merge hosts with deduplication while preserving config order
	hostMap := make(map[string]bool) // Track which hosts we've already seen

	// Add config hosts first (preserving their original order)
	for _, host := range configHosts {
		key := strings.ToLower(host.Name)
		if !hostMap[key] {
			allHosts = append(allHosts, host)
			hostMap[key] = true
		}
	}

	// Add known_hosts entries that aren't already in config
	for _, host := range knownHosts {
		key := strings.ToLower(host.Name)
		if !hostMap[key] {
			allHosts = append(allHosts, host)
			hostMap[key] = true
		}
	}

	// No sorting needed - hosts are already in config order, with known_hosts appended

	return allHosts, nil
}

// FilterHosts filters hosts by a search term
func FilterHosts(hosts []SSHHost, searchTerm string) []SSHHost {
	if searchTerm == "" {
		return hosts
	}

	var filtered []SSHHost
	searchLower := strings.ToLower(searchTerm)

	// Prioritization strategy with aliases:
	// 1) Hosts where an alias exactly matches the search term (highest priority)
	// 2) Primary name/hostname prefix matches
	// 3) Primary name/hostname substring matches
	// 4) Alias prefix matches
	// 5) Alias substring matches
	// Within each bucket, preserve input order.

	var exactAliasMatches []SSHHost
	var primaryPrefix []SSHHost
	var primaryContains []SSHHost
	var aliasPrefix []SSHHost
	var aliasContains []SSHHost

	for _, host := range hosts {
		nameLower := strings.ToLower(host.Name)
		hostLower := strings.ToLower(host.HostName)

		// Check exact alias matches first
		matchedExactAlias := false
		for _, a := range host.Aliases {
			if strings.ToLower(a) == searchLower {
				exactAliasMatches = append(exactAliasMatches, host)
				matchedExactAlias = true
				break
			}
		}
		if matchedExactAlias {
			continue
		}

		// Primary name/hostname prefix
		if strings.HasPrefix(nameLower, searchLower) || strings.HasPrefix(hostLower, searchLower) {
			primaryPrefix = append(primaryPrefix, host)
			continue
		}

		// Primary name/hostname contains
		if strings.Contains(nameLower, searchLower) || strings.Contains(hostLower, searchLower) {
			primaryContains = append(primaryContains, host)
			continue
		}

		// Check aliases for prefix/contains matches
		aliasMatched := false
		for _, a := range host.Aliases {
			aLower := strings.ToLower(a)
			if strings.HasPrefix(aLower, searchLower) {
				aliasPrefix = append(aliasPrefix, host)
				aliasMatched = true
				break
			}
			if strings.Contains(aLower, searchLower) {
				aliasContains = append(aliasContains, host)
				aliasMatched = true
				break
			}
		}
		if aliasMatched {
			continue
		}
	}

	// Append groups in priority order
	filtered = append(filtered, exactAliasMatches...)
	filtered = append(filtered, primaryPrefix...)
	filtered = append(filtered, primaryContains...)
	filtered = append(filtered, aliasPrefix...)
	filtered = append(filtered, aliasContains...)

	return filtered
}

// FormatHostDisplay formats a host for display in the TUI
func FormatHostDisplay(host SSHHost) string {
	var lines []string

	// Main host line - just the name (aliases are handled in TUI styling)
	lines = append(lines, host.Name)

	// Build details line with hostname, user, port
	var details []string

	// Add hostname if different from name
	if host.HostName != "" && host.HostName != host.Name {
		details = append(details, fmt.Sprintf("host: %s", host.HostName))
	}

	if host.Port != "" && host.Port != "22" {
		details = append(details, fmt.Sprintf("port: %s", host.Port))
	}

	if host.User != "" {
		details = append(details, fmt.Sprintf("user: %s", host.User))
	}

	if len(details) > 0 {
		lines = append(lines, strings.Join(details, " • "))
	}

	return strings.Join(lines, "\n")
}
