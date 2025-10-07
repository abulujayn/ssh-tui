package main

import (
	"fmt"
	"log"
	"os"
	"ssh-tui/internal/parser"
	"ssh-tui/internal/ssh"
	"ssh-tui/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	// Check if SSH is available
	if err := ssh.CheckSSHAvailable(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Discover hosts
	hosts, err := parser.DiscoverHosts()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering SSH hosts: %v\n", err)
		os.Exit(1)
	}

	if len(hosts) == 0 {
		showNoHostsMessage()
		os.Exit(1)
	}

	// Run the TUI flow
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
	// Step 1: Host Selection
	hostSelector := tui.NewHostSelectorModel(hosts)

	p := tea.NewProgram(hostSelector, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run host selector: %w", err)
	}

	hostModel, ok := finalModel.(*tui.HostSelectorModel)
	if !ok {
		return fmt.Errorf("unexpected model type from host selector")
	}

	selectedHost := hostModel.GetSelectedHost()
	if selectedHost == nil || !hostModel.IsSelected() {
		// User cancelled or no selection made
		return nil
	}

	// Step 2: Options Entry
	optionsEntry := tui.NewOptionsEntryModel(selectedHost)

	p = tea.NewProgram(optionsEntry, tea.WithAltScreen())
	finalModel, err = p.Run()
	if err != nil {
		return fmt.Errorf("failed to run options entry: %w", err)
	}

	optionsModel, ok := finalModel.(*tui.OptionsEntryModel)
	if !ok {
		return fmt.Errorf("unexpected model type from options entry")
	}

	if optionsModel.IsCancelled() {
		// User went back, restart host selection
		return runTUIFlow(hosts)
	}

	if !optionsModel.IsConfirmed() {
		// User cancelled
		return nil
	}

	options := optionsModel.GetOptions()

	// Step 3: Command Preview
	commandPreview := tui.NewCommandPreviewModel(selectedHost, options)

	p = tea.NewProgram(commandPreview, tea.WithAltScreen())
	finalModel, err = p.Run()
	if err != nil {
		return fmt.Errorf("failed to run command preview: %w", err)
	}

	previewModel, ok := finalModel.(*tui.CommandPreviewModel)
	if !ok {
		return fmt.Errorf("unexpected model type from command preview")
	}

	if previewModel.IsCancelled() {
		// User went back, restart from options entry
		return runOptionsFlow(selectedHost, hosts)
	}

	if !previewModel.IsConfirmed() {
		// User cancelled
		return nil
	}

	// Step 4: Execute SSH Command
	command := previewModel.GetCommand()

	// Validate the command before execution
	if err := ssh.ValidateSSHCommand(command); err != nil {
		return fmt.Errorf("invalid SSH command: %w", err)
	}

	// Show a message before executing
	fmt.Printf("Connecting to %s...\n\n", selectedHost.Name)

	// Execute the SSH command
	if err := ssh.ExecuteSSHCommand(command); err != nil {
		return fmt.Errorf("SSH execution failed: %w", err)
	}

	return nil
}

// runOptionsFlow runs the options entry and subsequent steps (for back navigation)
func runOptionsFlow(selectedHost *parser.SSHHost, hosts []parser.SSHHost) error {
	// Step 2: Options Entry
	optionsEntry := tui.NewOptionsEntryModel(selectedHost)

	p := tea.NewProgram(optionsEntry, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run options entry: %w", err)
	}

	optionsModel, ok := finalModel.(*tui.OptionsEntryModel)
	if !ok {
		return fmt.Errorf("unexpected model type from options entry")
	}

	if optionsModel.IsCancelled() {
		// User went back to host selection
		return runTUIFlow(hosts)
	}

	if !optionsModel.IsConfirmed() {
		// User cancelled
		return nil
	}

	options := optionsModel.GetOptions()

	// Step 3: Command Preview
	commandPreview := tui.NewCommandPreviewModel(selectedHost, options)

	p = tea.NewProgram(commandPreview, tea.WithAltScreen())
	finalModel, err = p.Run()
	if err != nil {
		return fmt.Errorf("failed to run command preview: %w", err)
	}

	previewModel, ok := finalModel.(*tui.CommandPreviewModel)
	if !ok {
		return fmt.Errorf("unexpected model type from command preview")
	}

	if previewModel.IsCancelled() {
		// User went back, restart options entry
		return runOptionsFlow(selectedHost, hosts)
	}

	if !previewModel.IsConfirmed() {
		// User cancelled
		return nil
	}

	// Step 4: Execute SSH Command
	command := previewModel.GetCommand()

	// Validate the command before execution
	if err := ssh.ValidateSSHCommand(command); err != nil {
		return fmt.Errorf("invalid SSH command: %w", err)
	}

	// Show a message before executing
	fmt.Printf("Connecting to %s...\n\n", selectedHost.Name)

	// Execute the SSH command
	if err := ssh.ExecuteSSHCommand(command); err != nil {
		return fmt.Errorf("SSH execution failed: %w", err)
	}

	return nil
}

func init() {
	// Set up logging to suppress tea debug output
	log.SetOutput(os.Stderr)
}
