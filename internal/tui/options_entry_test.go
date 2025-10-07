package tui

import (
	"ssh-tui/internal/parser"
	"testing"
)

func TestBuildSSHCommand_ConfigHost(t *testing.T) {
	// Test host from SSH config
	configHost := &parser.SSHHost{
		Name:     "myserver",
		HostName: "example.com",
		User:     "myuser",
		Port:     "2222",
		Source:   "config",
	}

	// Test without additional options
	command := buildSSHCommand(configHost, "")
	expected := "ssh myserver"
	if command != expected {
		t.Errorf("Expected '%s', got '%s'", expected, command)
	}

	// Test with additional options
	command = buildSSHCommand(configHost, "-L 8080:localhost:80")
	expected = "ssh -L 8080:localhost:80 myserver"
	if command != expected {
		t.Errorf("Expected '%s', got '%s'", expected, command)
	}
}

func TestBuildSSHCommand_KnownHost(t *testing.T) {
	// Test host from known_hosts (should expand options)
	knownHost := &parser.SSHHost{
		Name:     "server.example.com",
		HostName: "server.example.com",
		User:     "admin",
		Port:     "2222",
		Source:   "known_hosts",
	}

	// Test without additional options
	command := buildSSHCommand(knownHost, "")
	expected := "ssh -p 2222 admin@server.example.com"
	if command != expected {
		t.Errorf("Expected '%s', got '%s'", expected, command)
	}

	// Test with additional options
	command = buildSSHCommand(knownHost, "-i ~/.ssh/key")
	expected = "ssh -i ~/.ssh/key -p 2222 admin@server.example.com"
	if command != expected {
		t.Errorf("Expected '%s', got '%s'", expected, command)
	}
}

func TestBuildSSHCommand_ConfigHostWithDefaultPort(t *testing.T) {
	// Test config host with default port 22 (should still use host name only)
	configHost := &parser.SSHHost{
		Name:     "webserver",
		HostName: "web.company.com",
		User:     "deploy",
		Port:     "22", // Default port
		Source:   "config",
	}

	command := buildSSHCommand(configHost, "")
	expected := "ssh webserver"
	if command != expected {
		t.Errorf("Expected '%s', got '%s'", expected, command)
	}
}

func TestBuildSSHCommand_KnownHostWithDefaultPort(t *testing.T) {
	// Test known_hosts with default port (should not include -p 22)
	knownHost := &parser.SSHHost{
		Name:     "server.example.com",
		HostName: "server.example.com",
		User:     "admin",
		Port:     "22", // Default port
		Source:   "known_hosts",
	}

	command := buildSSHCommand(knownHost, "")
	expected := "ssh admin@server.example.com"
	if command != expected {
		t.Errorf("Expected '%s', got '%s'", expected, command)
	}
}
