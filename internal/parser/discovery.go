package parser

import (
	"fmt"
	"strings"
)

// matchesTerm checks if the search term matches any of the given strings in various ways
// Returns the priority level: 1=exact, 2=prefix, 3=contains, 0=no match
func matchesTerm(searchTerm string, candidates []string) int {
	searchLower := strings.ToLower(searchTerm)
	for _, cand := range candidates {
		if strings.EqualFold(cand, searchTerm) {
			return 1 // exact
		}
		candLower := strings.ToLower(cand)
		if strings.HasPrefix(candLower, searchLower) {
			return 2 // prefix
		}
		if strings.Contains(candLower, searchLower) {
			return 3 // contains
		}
	}
	return 0 // no match
}

// DiscoverHosts discovers all SSH hosts from both config and known_hosts files
func DiscoverHosts() ([]SSHHost, error) {
	var allHosts []SSHHost

	configHosts, err := ParseSSHConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH config: %w", err)
	}

	knownHosts, err := ParseKnownHosts()
	if err != nil {
		return nil, fmt.Errorf("failed to parse known_hosts: %w", err)
	}

	// Merge hosts with deduplication while preserving config order
	hostMap := make(map[string]bool)

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

	return allHosts, nil
}

// FilterHosts filters hosts by a search term
func FilterHosts(hosts []SSHHost, searchTerm string) []SSHHost {
	if searchTerm == "" {
		return hosts
	}

	var filtered []SSHHost

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
		// Check exact alias match first
		if match := matchesTerm(searchTerm, host.Aliases); match == 1 {
			exactAliasMatches = append(exactAliasMatches, host)
			continue
		}

		// Primary name/hostname matches
		primaryCandidates := []string{host.Name, host.HostName}
		if match := matchesTerm(searchTerm, primaryCandidates); match > 0 {
			if match == 2 {
				primaryPrefix = append(primaryPrefix, host)
			} else {
				primaryContains = append(primaryContains, host)
			}
			continue
		}

		// Alias prefix/contains matches
		if match := matchesTerm(searchTerm, host.Aliases); match > 0 {
			if match == 2 {
				aliasPrefix = append(aliasPrefix, host)
			} else {
				aliasContains = append(aliasContains, host)
			}
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

	lines = append(lines, host.Name)

	var details []string

	if host.HostName != "" && host.HostName != host.Name {
		details = append(details, fmt.Sprintf("host: %s", host.HostName))
	}

	if host.Port != "" && host.Port != DefaultSSHPort {
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
