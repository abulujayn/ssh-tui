package ssh

import (
	"os"
	"testing"

	"ssh-tui/internal/parser"
)

func TestValidateSSHCommand(t *testing.T) {
	cases := []struct {
		in      string
		wantErr bool
	}{
		{"", true},
		{"notssh host", true},
		{"ssh", true},
		{"ssh -p 2222 -i key", false}, // current implementation treats final token as host (even if used as option arg)
		{"ssh -p -C", true},           // last tokens are all option-like, no non-option host
		{"ssh user@host", false},
		{"ssh -p 2222 user@host", false},
		{"ssh   user@host", false},      // multiple spaces
		{"ssh\tuser@host", false},       // tabs
		{"ssh host with spaces", false}, // host with spaces (though invalid, but validation allows)
		{"ssh -o option=value host", false},
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
	cfgHost := &parser.SSHHost{Name: "cfg", HostName: "cfg.example", Source: "config"}
	cmd := BuildSSHCommand(cfgHost, "-L 8080:localhost:80 -i ~/.ssh/id_rsa")
	if cmd != "ssh cfg -L 8080:localhost:80 -i ~/.ssh/id_rsa" {
		t.Fatalf("unexpected command: %q", cmd)
	}

	// known_hosts-like with no user, HostName == Name, default port
	known := &parser.SSHHost{Name: "host.local", HostName: "host.local", Port: parser.DefaultSSHPort, Source: parser.SourceKnownHosts}
	cmd = BuildSSHCommand(known, "")
	if cmd != "ssh host.local" {
		t.Fatalf("unexpected command for known host: %q", cmd)
	}

	// known_hosts-like with user and non-default port
	known2 := &parser.SSHHost{Name: "srv", HostName: "srv.example", User: "admin", Port: "2200", Source: parser.SourceKnownHosts}
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
