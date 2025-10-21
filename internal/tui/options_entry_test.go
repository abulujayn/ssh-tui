package tui

import (
	"ssh-tui/internal/parser"
	"ssh-tui/internal/ssh"
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
	command := ssh.BuildSSHCommand(configHost, "")
	expected := "ssh myserver"
	if command != expected {
		t.Errorf("Expected '%s', got '%s'", expected, command)
	}

	// Test with additional options
	command = ssh.BuildSSHCommand(configHost, "-L 8080:localhost:80")
	expected = "ssh myserver -L 8080:localhost:80"
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
		Source:   parser.SourceKnownHosts,
	}

	// Test without additional options
	command := ssh.BuildSSHCommand(knownHost, "")
	expected := "ssh -p 2222 admin@server.example.com"
	if command != expected {
		t.Errorf("Expected '%s', got '%s'", expected, command)
	}

	// Test with additional options
	command = ssh.BuildSSHCommand(knownHost, "-i ~/.ssh/key")
	expected = "ssh -p 2222 admin@server.example.com -i ~/.ssh/key"
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

	command := ssh.BuildSSHCommand(configHost, "")
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
		Port:     parser.DefaultSSHPort, // Default port
		Source:   parser.SourceKnownHosts,
	}

	command := ssh.BuildSSHCommand(knownHost, "")
	expected := "ssh admin@server.example.com"
	if command != expected {
		t.Errorf("Expected '%s', got '%s'", expected, command)
	}
}
