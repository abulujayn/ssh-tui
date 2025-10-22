package parser

import (
	"net"
	"regexp"
	"strings"
)

var domainRegex = regexp.MustCompile(`^(?i)[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*\.[a-z]{2,}$`)

// IsValidSSHOption returns true if the SSH option string is safe (no shell metacharacters)
func IsValidSSHOption(option string) bool {
	// Disallow common shell metacharacters that could lead to command injection
	dangerousChars := ";|&`$()<>[]{}\"'\\"
	return !strings.ContainsAny(option, dangerousChars)
}

// IsValidHost returns true if the input is a valid IP address or domain name
func IsValidHost(input string) bool {
	input = strings.TrimSpace(input)
	if input == "" {
		return false
	}

	var hostPart string
	if at := strings.LastIndex(input, "@"); at != -1 {
		hostPart = input[at+1:]
		if hostPart == "" {
			return false
		}
	} else {
		hostPart = input
	}

	if net.ParseIP(hostPart) != nil {
		return true
	}
	// Check if it's a valid domain name (RFC 1035, simplified)
	return domainRegex.MatchString(hostPart)
}

// ParseUserHost splits a string like user@host into user and host parts
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
