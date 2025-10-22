package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"ssh-tui/internal/parser"
)

// ExecuteSSHCommand executes the SSH command with proper handling for different platforms
func ExecuteSSHCommand(command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	sshPath, err := exec.LookPath("ssh")
	if err != nil {
		return fmt.Errorf("ssh command not found in PATH: %w", err)
	}

	args := parts[1:]

	// Print only the full command before executing so the user can see exactly what's run
	if command != "ssh" {
		fmt.Println("\x1b[1;36m" + strings.Join(parts, " ") + "\x1b[0m")
	}

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

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("SSH connection failed with exit code %d", exitError.ExitCode())
		}
		return fmt.Errorf("failed to execute SSH command: %w", err)
	}

	return nil
}

// executeSSHUnix executes SSH on Unix-like systems using exec.Syscall for proper signal handling
func executeSSHUnix(sshPath string, args []string) error {
	execArgs := append([]string{sshPath}, args...)

	env := os.Environ()

	err := syscall.Exec(sshPath, execArgs, env)

	return fmt.Errorf("failed to execute SSH command: %w", err)
}

// ValidateSSHCommand performs comprehensive validation of the SSH command
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

	// Find the host argument (last non-option argument)
	var host string
	for i := len(parts) - 1; i >= 1; i-- {
		if !strings.HasPrefix(parts[i], "-") {
			host = parts[i]
			break
		}
	}

	if host == "" {
		return fmt.Errorf("no target host specified in SSH command")
	}

	// Validate the host
	if !parser.IsValidHost(host) {
		return fmt.Errorf("invalid host: %s", host)
	}

	// Check for port option if present
	for i := 1; i < len(parts)-1; i++ {
		if parts[i] == "-p" && i+1 < len(parts) {
			portStr := parts[i+1]
			port, err := strconv.Atoi(portStr)
			if err != nil || port < 1 || port > 65535 {
				return fmt.Errorf("invalid port: %s", portStr)
			}
		}
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
