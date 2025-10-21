package optionsentry

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"ssh-tui/internal/tui/helpers"
	"ssh-tui/internal/tui/labels"
	"ssh-tui/internal/tui/styles"
)

// View implements the tea.Model interface for options entry
func (m *OptionsEntryModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(styles.TitleStyle.Render("SSH Options & Arguments") + "\n\n")

	// Selected host info in table format
	b.WriteString(m.renderHostInfoTable())
	if m.host.Source == "config" {
		b.WriteString("\n\n")
	}

	// Options label
	b.WriteString(styles.TitleStyle.Render("Options:") + "\n")

	// Input field
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(0, 1).
		Width(max(60, m.width-10))

	// Create the input display with cursor (delegated to helper)
	rendered := helpers.RenderInputWithCursor(m.options, m.cursor, max(60, m.width-10))
	b.WriteString(inputStyle.Render(rendered))

	b.WriteString("\n")

	// Examples label (muted description under the input)
	b.WriteString(styles.InstructionStyle.Render(labels.ExamplesText) + "\n\n")

	// Command Preview Section
	b.WriteString(styles.TitleStyle.Render("Command Preview:") + "\n")

	// Show the current command that would be executed as plain text
	currentCommand := m.GetCommand()
	// Render as plain default text (no styling)
	b.WriteString(currentCommand + "\n\n")

	// Main instructions at the bottom
	b.WriteString(styles.InstructionStyle.Render("Use Enter to execute, Esc to go back") + "\n\n")

	return b.String()
}

// renderHostInfoTable renders the selected host information in a table format
func (m *OptionsEntryModel) renderHostInfoTable() string {
	// Only show table for hosts from sshconfig
	if m.host.Source != "config" {
		return ""
	}

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
	tableContent.WriteString(headerStyle.Render("Host configuration"))
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
