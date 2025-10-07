package tui

import (
	"ssh-tui/internal/parser"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// OptionsEntryModel represents the SSH options entry screen
type OptionsEntryModel struct {
	host      *parser.SSHHost
	options   string
	cursor    int
	confirmed bool
	cancelled bool
	width     int
	height    int
}

// NewOptionsEntryModel creates a new options entry model
func NewOptionsEntryModel(host *parser.SSHHost) *OptionsEntryModel {
	return &OptionsEntryModel{
		host:      host,
		options:   "",
		cursor:    len(""),
		confirmed: false,
		cancelled: false,
	}
}

// Init implements the tea.Model interface
func (m *OptionsEntryModel) Init() tea.Cmd {
	return nil
}

// Update implements the tea.Model interface
func (m *OptionsEntryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			m.cancelled = true
			return m, tea.Quit

		case "enter":
			m.confirmed = true
			return m, tea.Quit

		case "left", "ctrl+b":
			if m.cursor > 0 {
				m.cursor--
			}

		case "right", "ctrl+f":
			if m.cursor < len(m.options) {
				m.cursor++
			}

		case "home", "ctrl+a":
			m.cursor = 0

		case "end", "ctrl+e":
			m.cursor = len(m.options)

		case "backspace", "ctrl+h":
			if m.cursor > 0 && len(m.options) > 0 {
				m.options = m.options[:m.cursor-1] + m.options[m.cursor:]
				m.cursor--
			}

		case "delete", "ctrl+d":
			if m.cursor < len(m.options) {
				m.options = m.options[:m.cursor] + m.options[m.cursor+1:]
			}

		case "ctrl+u":
			// Delete from cursor to beginning
			m.options = m.options[m.cursor:]
			m.cursor = 0

		case "ctrl+k":
			// Delete from cursor to end
			m.options = m.options[:m.cursor]

		case "ctrl+w":
			// Delete word backwards
			if m.cursor > 0 {
				// Find start of current word
				start := m.cursor - 1
				for start > 0 && m.options[start-1] != ' ' {
					start--
				}
				m.options = m.options[:start] + m.options[m.cursor:]
				m.cursor = start
			}

		default:
			// Handle regular character input
			if len(msg.String()) == 1 {
				char := msg.String()
				// Insert character at cursor position
				m.options = m.options[:m.cursor] + char + m.options[m.cursor:]
				m.cursor++
			}
		}
	}

	return m, nil
}

// View implements the tea.Model interface
func (m *OptionsEntryModel) View() string {
	var b strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("183")).
		Bold(true).
		Render("SSH Options & Arguments")

	b.WriteString(title + "\n\n")

	// Selected host info in table format
	b.WriteString(m.renderHostInfoTable() + "\n\n")

	// Description
	descriptionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	b.WriteString(descriptionStyle.Render("Enter SSH options and arguments (optional):") + "\n")
	b.WriteString(descriptionStyle.Render("Examples: -L 8080:localhost:80 -i ~/.ssh/id_rsa -p 2222 -X") + "\n\n")

	// Input field
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(0, 1).
		Width(max(60, m.width-10))

	// Create the input display with cursor
	var inputDisplay string
	if len(m.options) == 0 {
		inputDisplay = " "
	} else {
		inputDisplay = m.options
	}

	// Add cursor indicator
	if m.cursor <= len(inputDisplay) {
		cursorStyle := lipgloss.NewStyle().
			Background(lipgloss.Color("86")).
			Foreground(lipgloss.Color("0"))

		if m.cursor == len(inputDisplay) {
			inputDisplay += cursorStyle.Render(" ")
		} else {
			char := string(inputDisplay[m.cursor])
			inputDisplay = inputDisplay[:m.cursor] +
				cursorStyle.Render(char) +
				inputDisplay[m.cursor+1:]
		}
	}

	b.WriteString(inputStyle.Render(inputDisplay))

	b.WriteString("\n\n")

	// Command Preview Section
	commandPreviewTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("183")).
		Bold(true).
		Render("Command Preview:")

	b.WriteString(commandPreviewTitle + "\n")

	// Show the current command that would be executed
	currentCommand := m.GetCommand()
	commandStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(0, 1).
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("252")).
		Bold(true).
		Width(max(60, m.width-10))

	b.WriteString(commandStyle.Render(currentCommand) + "\n\n")

	// Main instructions at the bottom
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true)

	executeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")).
		Bold(true)

	b.WriteString(executeStyle.Render("Press Enter to execute") + " • " +
		instructionStyle.Render("Esc to go back • Ctrl+C to quit") + "\n\n")

	// Keyboard shortcuts help
	shortcuts := []string{
		"Ctrl+A: Home",
		"Ctrl+E: End",
		"Ctrl+U: Clear to beginning",
		"Ctrl+K: Clear to end",
		"Ctrl+W: Delete word back",
	}

	shortcutStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true)

	b.WriteString(shortcutStyle.Render(strings.Join(shortcuts, " • ")))

	return b.String()
}

// GetOptions returns the entered options
func (m *OptionsEntryModel) GetOptions() string {
	return strings.TrimSpace(m.options)
}

// IsConfirmed returns whether the user confirmed the options
func (m *OptionsEntryModel) IsConfirmed() bool {
	return m.confirmed
}

// IsCancelled returns whether the user cancelled the options entry
func (m *OptionsEntryModel) IsCancelled() bool {
	return m.cancelled
}

// renderHostInfoTable renders the selected host information in a table format
func (m *OptionsEntryModel) renderHostInfoTable() string {
	// Table styles
	tableStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 2)

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("183")).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true).
		Width(12).
		Align(lipgloss.Left)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))

	var tableContent strings.Builder

	// Table header
	tableContent.WriteString(headerStyle.Render("Current Configuration"))
	tableContent.WriteString("\n\n")

	// Host name
	tableContent.WriteString(labelStyle.Render("Host:"))
	tableContent.WriteString("  ")
	tableContent.WriteString(valueStyle.Render(m.host.HostName))
	tableContent.WriteString("\n")

	// Port (always show)
	port := m.host.Port
	if port == "" {
		port = "22" // Default SSH port
	}
	tableContent.WriteString(labelStyle.Render("Port:"))
	tableContent.WriteString("  ")
	tableContent.WriteString(valueStyle.Render(port))
	tableContent.WriteString("\n")

	// User
	if m.host.User != "" {
		tableContent.WriteString(labelStyle.Render("User:"))
		tableContent.WriteString("  ")
		tableContent.WriteString(valueStyle.Render(m.host.User))
		tableContent.WriteString("\n")
	}

	return tableStyle.Render(tableContent.String())
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// GetCommand returns the SSH command that would be executed with current options
func (m *OptionsEntryModel) GetCommand() string {
	return buildSSHCommand(m.host, m.options)
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
