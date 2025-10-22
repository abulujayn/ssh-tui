package ssh

import (
	"os"
	"testing"

	"ssh-tui/internal/types"
)

func TestValidateSSHCommand(t *testing.T) {
	cases := []struct {
		in      string
		wantErr bool
	}{
		{"", true},
		{"notssh host", true},
		{"ssh", true},
		{"ssh -p 2222 -i ~/.ssh/key example.com", false}, // valid host
		{"ssh -p -C", true}, // last tokens are all option-like, no non-option host
		{"ssh user@example.com", false},
		{"ssh -p 2222 user@example.com", false},
		{"ssh   user@example.com", false}, // multiple spaces
		{"ssh\tuser@example.com", false},  // tabs
		{"ssh example.com", false},        // valid host
		{"ssh -o option=value example.com", false},
		{"ssh invalid_host", true},       // invalid host
		{"ssh user@invalid..host", true}, // invalid domain
		{"ssh -p 0 host", true},          // invalid port 0
		{"ssh -p 99999 host", true},      // invalid port > 65535
		{"ssh -p abc host", true},        // non-numeric port
		{"ssh 192.168.1.1", false},       // valid IP
		{"ssh user@192.168.1.1", false},  // valid user@IP
		{"ssh example.com", false},       // valid domain
		{"ssh sub.example.com", false},   // valid subdomain
		{"ssh -p 22 valid.host", false},  // valid with port
	}

	for _, c := range cases {
		err := ValidateSSHCommand(c.in)
		if (err != nil) != c.wantErr {
			t.Fatalf("ValidateSSHCommand(%q) error = %v, wantErr=%v", c.in, err, c.wantErr)
		}
	}
}

func TestBuildSSHCommand_AdditionalCases(t *testing.T) {
	// config source with options that include spaces
	cfgHost := &types.SSHHost{Name: "cfg", HostName: "cfg.example", Source: types.SourceConfig}
	cmd := BuildSSHCommand(cfgHost, "-L 8080:localhost:80 -i ~/.ssh/id_rsa")
	if cmd != "ssh cfg -L 8080:localhost:80 -i ~/.ssh/id_rsa" {
		t.Fatalf("unexpected command: %q", cmd)
	}

	// known_hosts-like with no user, HostName == Name, default port
	known := &types.SSHHost{Name: "host.local", HostName: "host.local", Port: types.DefaultSSHPort, Source: types.SourceKnownHosts}
	cmd = BuildSSHCommand(known, "")
	if cmd != "ssh host.local" {
		t.Fatalf("unexpected command for known host: %q", cmd)
	}

	// known_hosts-like with user and non-default port
	known2 := &types.SSHHost{Name: "srv", HostName: "srv.example", User: "admin", Port: "2200", Source: types.SourceKnownHosts}
	cmd = BuildSSHCommand(known2, "-v")
	if cmd != "ssh -p 2200 admin@srv.example -v" {
		t.Fatalf("unexpected command for known2: %q", cmd)
	}

	// Test with empty options
	cmd = BuildSSHCommand(known2, "")
	if cmd != "ssh -p 2200 admin@srv.example" {
		t.Fatalf("unexpected command with empty options: %q", cmd)
	}
}

func TestCheckSSHAvailable_Failure(t *testing.T) {
	// Save original PATH and clear it to force LookPath to fail
	orig := os.Getenv("PATH")
	defer os.Setenv("PATH", orig)

	os.Setenv("PATH", "")

	if err := CheckSSHAvailable(); err == nil {
		t.Fatalf("expected CheckSSHAvailable to fail when PATH is empty")
	}
}

func TestBuildSSHCommand_ConfigHost(t *testing.T) {
	// Test host from SSH config
	configHost := &types.SSHHost{
		Name:     "myserver",
		HostName: "example.com",
		User:     "myuser",
		Port:     "2222",
		Source:   "config",
	}

	// Test without additional options
	command := BuildSSHCommand(configHost, "")
	expected := "ssh myserver"
	if command != expected {
		t.Errorf("Expected '%s', got '%s'", expected, command)
	}

	// Test with additional options
	command = BuildSSHCommand(configHost, "-L 8080:localhost:80")
	expected = "ssh myserver -L 8080:localhost:80"
	if command != expected {
		t.Errorf("Expected '%s', got '%s'", expected, command)
	}
}

func TestBuildSSHCommand_KnownHost(t *testing.T) {
	// Test host from known_hosts (should expand options)
	knownHost := &types.SSHHost{
		Name:     "server.example.com",
		HostName: "server.example.com",
		User:     "admin",
		Port:     "2222",
		Source:   types.SourceKnownHosts,
	}

	// Test without additional options
	command := BuildSSHCommand(knownHost, "")
	expected := "ssh -p 2222 admin@server.example.com"
	if command != expected {
		t.Errorf("Expected '%s', got '%s'", expected, command)
	}

	// Test with additional options
	command = BuildSSHCommand(knownHost, "-i ~/.ssh/key")
	expected = "ssh -p 2222 admin@server.example.com -i ~/.ssh/key"
	if command != expected {
		t.Errorf("Expected '%s', got '%s'", expected, command)
	}
}

func TestBuildSSHCommand_ConfigHostWithDefaultPort(t *testing.T) {
	// Test config host with default port 22 (should still use host name only)
	configHost := &types.SSHHost{
		Name:     "webserver",
		HostName: "web.company.com",
		User:     "deploy",
		Port:     "22", // Default port
		Source:   "config",
	}

	command := BuildSSHCommand(configHost, "")
	expected := "ssh webserver"
	if command != expected {
		t.Errorf("Expected '%s', got '%s'", expected, command)
	}
}

func TestBuildSSHCommand_KnownHostWithDefaultPort(t *testing.T) {
	// Test known_hosts with default port (should not include -p 22)
	knownHost := &types.SSHHost{
		Name:     "server.example.com",
		HostName: "server.example.com",
		User:     "admin",
		Port:     types.DefaultSSHPort, // Default port
		Source:   types.SourceKnownHosts,
	}

	command := BuildSSHCommand(knownHost, "")
	expected := "ssh admin@server.example.com"
	if command != expected {
		t.Errorf("Expected '%s', got '%s'", expected, command)
	}
}
