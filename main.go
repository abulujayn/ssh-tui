package main

import (
	"fmt"
	"log"
	"os"
	"ssh-tui/internal/parser"
	"ssh-tui/internal/ssh"
	"ssh-tui/internal/tui/hostselector"
	"ssh-tui/internal/tui/optionsentry"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Version is the application version string. Bump this when releasing.
const Version = "1.0.0-alpha"

func main() {

	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "--help":
			_ = ssh.ExecuteSSHCommand("ssh")
			return
		case "--version":
			fmt.Printf("ssh-tui %s\n", Version)
			return
		}
	}

	// If CLI args were provided, treat them as a direct ssh invocation and execute immediately
	if len(os.Args) > 1 {
		cmd := "ssh " + strings.Join(os.Args[1:], " ")

		_ = ssh.ExecuteSSHCommand(cmd)
		return
	}

	hosts, err := parser.DiscoverHosts()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering SSH hosts: %v\n", err)
		os.Exit(1)
	}

	if len(hosts) == 0 {
		showNoHostsMessage()
		os.Exit(1)
	}

	if err := runTUIFlow(hosts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// showNoHostsMessage displays a helpful message when no hosts are found
func showNoHostsMessage() {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("203")).
		Bold(true)

	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	fmt.Println(titleStyle.Render("No SSH hosts found!"))
	fmt.Println()
	fmt.Println(messageStyle.Render("SSH-TUI couldn't find any hosts to connect to."))
	fmt.Println()
	fmt.Println(messageStyle.Render("To use SSH-TUI, you need hosts configured in:"))
	fmt.Println(messageStyle.Render("• ~/.ssh/config - SSH configuration file"))
	fmt.Println(messageStyle.Render("• ~/.ssh/known_hosts - Previously connected hosts"))
	fmt.Println()
	fmt.Println(messageStyle.Render("Example ~/.ssh/config entry:"))
	fmt.Println(messageStyle.Render("Host myserver"))
	fmt.Println(messageStyle.Render("    HostName example.com"))
	fmt.Println(messageStyle.Render("    User myuser"))
	fmt.Println(messageStyle.Render("    Port 22"))
}

// runTUIFlow runs the complete TUI flow
func runTUIFlow(hosts []parser.SSHHost) error {
	hostSelector := hostselector.NewHostSelectorModel(hosts)

	p := tea.NewProgram(hostSelector, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run host selector: %w", err)
	}

	hostModel, ok := finalModel.(*hostselector.HostSelectorModel)
	if !ok {
		return fmt.Errorf("unexpected model type from host selector")
	}

	selectedHost := hostModel.GetSelectedHost()
	if selectedHost == nil || !hostModel.IsSelected() {
		return nil
	}

	if hostModel.OpenOptionsRequested() {
		return runOptionsFlow(selectedHost, hosts)
	}

	// Build command using default/empty options
	command := ssh.BuildSSHCommand(selectedHost, "")

	if err := ssh.ValidateSSHCommand(command); err != nil {
		return fmt.Errorf("invalid SSH command: %w", err)
	}

	if err := ssh.ExecuteSSHCommand(command); err != nil {
		return fmt.Errorf("SSH execution failed: %w", err)
	}

	return nil
}

// runOptionsFlow runs the options entry and subsequent steps (for back navigation)
func runOptionsFlow(selectedHost *parser.SSHHost, hosts []parser.SSHHost) error {
	optionsEntry := optionsentry.NewOptionsEntryModel(selectedHost)

	p := tea.NewProgram(optionsEntry, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run options entry: %w", err)
	}

	optionsModel, ok := finalModel.(*optionsentry.OptionsEntryModel)
	if !ok {
		return fmt.Errorf("unexpected model type from options entry")
	}

	if optionsModel.IsCancelled() {
		return runTUIFlow(hosts)
	}

	if !optionsModel.IsConfirmed() {
		return nil
	}

	command := optionsModel.GetCommand()

	if err := ssh.ValidateSSHCommand(command); err != nil {
		return fmt.Errorf("invalid SSH command: %w", err)
	}

	if err := ssh.ExecuteSSHCommand(command); err != nil {
		return fmt.Errorf("SSH execution failed: %w", err)
	}

	return nil
}

func init() {
	// Set up logging to suppress tea debug output
	log.SetOutput(os.Stderr)
}
