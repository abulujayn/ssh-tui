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

	// Selected host info
	hostStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	b.WriteString("Selected host: " + hostStyle.Render(parser.FormatHostDisplay(*m.host)) + "\n\n")

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

	// Main instructions at the bottom
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true)

	b.WriteString(instructionStyle.Render("Press Enter to continue, Esc to go back, Ctrl+C to quit") + "\n\n")

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

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
