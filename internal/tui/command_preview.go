package tui

import (
	"ssh-tui/internal/parser"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CommandPreviewModel represents the SSH command preview screen
type CommandPreviewModel struct {
	host      *parser.SSHHost
	options   string
	command   string
	confirmed bool
	cancelled bool
	width     int
	height    int
}

// NewCommandPreviewModel creates a new command preview model
func NewCommandPreviewModel(host *parser.SSHHost, options string) *CommandPreviewModel {
	command := buildSSHCommand(host, options)

	return &CommandPreviewModel{
		host:      host,
		options:   options,
		command:   command,
		confirmed: false,
		cancelled: false,
	}
}

// Init implements the tea.Model interface
func (m *CommandPreviewModel) Init() tea.Cmd {
	return nil
}

// Update implements the tea.Model interface
func (m *CommandPreviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc", "n", "N":
			m.cancelled = true
			return m, tea.Quit

		case "enter", "y", "Y":
			m.confirmed = true
			return m, tea.Quit
		}
	}

	return m, nil
}

// View implements the tea.Model interface
func (m *CommandPreviewModel) View() string {
	var b strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("183")).
		Bold(true).
		Render("SSH Command Preview")

	b.WriteString(title + "\n\n")

	// Host information
	hostStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	b.WriteString("Host: " + hostStyle.Render(parser.FormatHostDisplay(*m.host)) + "\n")

	if m.options != "" {
		optionsStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("183"))

		b.WriteString("Options: " + optionsStyle.Render(m.options) + "\n")
	}

	b.WriteString("\n")

	// Command preview
	b.WriteString("The following SSH command will be executed:\n\n")

	commandStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(1, 2).
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("252")).
		Bold(true)

	b.WriteString(commandStyle.Render(m.command) + "\n\n")

	// Confirmation prompt
	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")).
		Bold(true)

	b.WriteString(promptStyle.Render("Do you want to proceed?") + "\n\n\n")

	// Instructions at the bottom
	yesStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")).
		Bold(true)

	noStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("203")).
		Bold(true)

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true)

	b.WriteString(yesStyle.Render("[Y]es") + " / " + noStyle.Render("[N]o") + " / " + instructionStyle.Render("Esc to go back") + "\n")
	b.WriteString(instructionStyle.Render("Press Enter to confirm, or Ctrl+C to quit"))

	return b.String()
}

// buildSSHCommand constructs the SSH command string
func buildSSHCommand(host *parser.SSHHost, options string) string {
	var parts []string

	parts = append(parts, "ssh")

	// Add user-provided options first
	if options != "" {
		parts = append(parts, options)
	}

	// Add host-specific options from SSH config
	if host.Port != "" && host.Port != "22" {
		parts = append(parts, "-p", host.Port)
	}

	// Construct the connection string
	var target string
	if host.User != "" {
		target = host.User + "@"
	}

	// Use HostName if available, otherwise use Name
	if host.HostName != "" && host.HostName != host.Name {
		target += host.HostName
	} else {
		target += host.Name
	}

	parts = append(parts, target)

	return strings.Join(parts, " ")
}

// GetCommand returns the SSH command to execute
func (m *CommandPreviewModel) GetCommand() string {
	return m.command
}

// IsConfirmed returns whether the user confirmed the command
func (m *CommandPreviewModel) IsConfirmed() bool {
	return m.confirmed
}

// IsCancelled returns whether the user cancelled the command
func (m *CommandPreviewModel) IsCancelled() bool {
	return m.cancelled
}
