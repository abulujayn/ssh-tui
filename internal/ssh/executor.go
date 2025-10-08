package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

// ExecuteSSHCommand executes the SSH command with proper handling for different platforms
func ExecuteSSHCommand(command string) error {
	// Parse the command into parts
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	// Find ssh executable
	sshPath, err := exec.LookPath("ssh")
	if err != nil {
		return fmt.Errorf("ssh command not found in PATH: %w", err)
	}

	// Prepare arguments (skip the "ssh" part since we have the full path)
	args := parts[1:]

	// Print only the full command before executing so the user can see exactly what's run.
	// Style it with ANSI color codes to differentiate from the SSH session output.
	// Bold cyan: \x1b[1;36m ... \x1b[0m
	fmt.Println("\x1b[1;36m" + strings.Join(parts, " ") + "\x1b[0m")

	// Execute based on platform
	switch runtime.GOOS {
	case "windows":
		return executeSSHWindows(sshPath, args)
	default:
		return executeSSHUnix(sshPath, args)
	}
}

// executeSSHWindows executes SSH on Windows
func executeSSHWindows(sshPath string, args []string) error {
	cmd := exec.Command(sshPath, args...)

	// Connect to stdin, stdout, and stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and wait for it to complete
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// SSH returned with non-zero exit code
			return fmt.Errorf("SSH connection failed with exit code %d", exitError.ExitCode())
		}
		return fmt.Errorf("failed to execute SSH command: %w", err)
	}

	return nil
}

// executeSSHUnix executes SSH on Unix-like systems using exec.Syscall for proper signal handling
func executeSSHUnix(sshPath string, args []string) error {
	// Prepare arguments for execve (include program name as first argument)
	execArgs := append([]string{sshPath}, args...)

	// Get current environment
	env := os.Environ()

	// Execute the command, replacing the current process
	err := syscall.Exec(sshPath, execArgs, env)

	// If we reach here, exec failed
	return fmt.Errorf("failed to execute SSH command: %w", err)
}

// ValidateSSHCommand performs basic validation of the SSH command
func ValidateSSHCommand(command string) error {
	if command == "" {
		return fmt.Errorf("command is empty")
	}

	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("command contains no arguments")
	}

	if parts[0] != "ssh" {
		return fmt.Errorf("command must start with 'ssh'")
	}

	if len(parts) < 2 {
		return fmt.Errorf("SSH command missing target host")
	}

	// Basic check for host argument (last non-option argument)
	var hostFound bool
	for i := len(parts) - 1; i >= 1; i-- {
		if !strings.HasPrefix(parts[i], "-") {
			hostFound = true
			break
		}
	}

	if !hostFound {
		return fmt.Errorf("no target host specified in SSH command")
	}

	return nil
}

// CheckSSHAvailable checks if SSH is available on the system
func CheckSSHAvailable() error {
	_, err := exec.LookPath("ssh")
	if err != nil {
		return fmt.Errorf("SSH is not available on this system. Please install OpenSSH: %w", err)
	}
	return nil
}
