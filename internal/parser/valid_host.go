package parser

import (
	"net"
	"regexp"
	"strings"
)

// IsValidHost returns true if the input is a valid IP address or domain name,
// optionally prefixed by user@ (e.g., user@host.name)
func IsValidHost(input string) bool {
	input = strings.TrimSpace(input)
	if input == "" {
		return false
	}
	// Split user@host if present
	var hostPart string
	if at := strings.LastIndex(input, "@"); at != -1 {
		hostPart = input[at+1:]
		if hostPart == "" {
			return false
		}
	} else {
		hostPart = input
	}
	// Check if it's a valid IP address
	if net.ParseIP(hostPart) != nil {
		return true
	}
	// Check if it's a valid domain name (RFC 1035, simplified)
	domainPattern := `^(?i)[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*\.[a-z]{2,}$`
	matched, _ := regexp.MatchString(domainPattern, hostPart)
	return matched
}

// ParseUserHost splits a string like user@host into user and host parts.
// If no user is present, user will be empty.
func ParseUserHost(input string) (user, host string) {
	input = strings.TrimSpace(input)
	if at := strings.LastIndex(input, "@"); at != -1 {
		user = input[:at]
		host = input[at+1:]
	} else {
		user = ""
		host = input
	}
	return
}
