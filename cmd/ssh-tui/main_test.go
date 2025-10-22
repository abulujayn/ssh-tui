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
