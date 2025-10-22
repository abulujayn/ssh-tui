package main

import (
	"os/exec"
	"strings"
	"testing"
)

// TestVersionFlag builds the binary and verifies that --version prints the Version string.
func TestVersionFlag(t *testing.T) {
	// Run using `go run` to execute main with --version
	run := exec.Command("go", "run", "./main.go", "--version")
	out, err := run.CombinedOutput()
	if err != nil {
		t.Fatalf("running go run failed: %v\n%s", err, string(out))
	}

	outStr := strings.TrimSpace(string(out))
	expected := "ssh-tui " + Version
	if outStr != expected {
		t.Fatalf("unexpected version output: got %q want %q", outStr, expected)
	}
}

// TestHelpFlag tests the --help flag
func TestHelpFlag(t *testing.T) {
	run := exec.Command("go", "run", "./main.go", "--help")
	out, err := run.CombinedOutput()
	if err != nil {
		t.Fatalf("running go run --help failed: %v\n%s", err, string(out))
	}

	outStr := string(out)
	if !strings.Contains(outStr, "ssh") {
		t.Fatalf("expected help output to contain 'ssh', got %q", outStr)
	}
}

// TestDirectSSHCommand tests passing arguments directly to SSH
func TestDirectSSHCommand(t *testing.T) {
	run := exec.Command("go", "run", "./main.go", "user@host")
	out, err := run.CombinedOutput()
	outStr := string(out)
	// Should print the command
	if !strings.Contains(outStr, "ssh user@host") {
		t.Fatalf("expected output to contain the SSH command, got %q", outStr)
	}
	// err may be nil or not, depending on ssh availability
	_ = err
}
