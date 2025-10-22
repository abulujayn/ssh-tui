package optionsentry

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"ssh-tui/internal/tui/helpers"
	"ssh-tui/internal/tui/ui"
	"ssh-tui/internal/types"
)

// View implements the tea.Model interface for options entry
func (m *OptionsEntryModel) View() string {
	var b strings.Builder

	b.WriteString(ui.TitleStyle.Render("SSH Options & Arguments") + "\n\n")

	// Selected host info in table format
	b.WriteString(m.renderHostInfoTable())
	if m.host.Source == types.SourceConfig {
		b.WriteString("\n\n")
	}

	b.WriteString(ui.TitleStyle.Render("Options:") + "\n")

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(0, 1).
		Width(max(60, m.width-10))

	// Create the input display with cursor (delegated to helper)
	rendered := helpers.RenderInputWithCursor(m.options, m.cursor, max(60, m.width-10))
	b.WriteString(inputStyle.Render(rendered))

	b.WriteString("\n")

	b.WriteString(ui.InstructionStyle.Render(ui.ExamplesText) + "\n\n")

	b.WriteString(ui.TitleStyle.Render("Command Preview:") + "\n")

	// Show the current command that would be executed
	currentCommand := m.GetCommand()
	b.WriteString(currentCommand + "\n\n")

	b.WriteString(ui.InstructionStyle.Render("Use Enter to execute, Esc to go back") + "\n\n")

	return b.String()
}

// renderHostInfoTable renders the selected host information in a table format
func (m *OptionsEntryModel) renderHostInfoTable() string {
	// Only show table for hosts from sshconfig
	if m.host.Source != types.SourceConfig {
		return ""
	}

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

	tableContent.WriteString(headerStyle.Render("Host configuration"))
	tableContent.WriteString("\n\n")

	// Host name
	m.addTableRow(&tableContent, labelStyle, valueStyle, "Host:", m.host.HostName)

	// Port
	port := m.host.Port
	if port == "" {
		port = types.DefaultSSHPort // Default SSH port
	}
	m.addTableRow(&tableContent, labelStyle, valueStyle, "Port:", port)

	// User
	if m.host.User != "" {
		m.addTableRow(&tableContent, labelStyle, valueStyle, "User:", m.host.User)
	}

	return tableStyle.Render(tableContent.String())
}

// addTableRow adds a labeled row to the table content
func (m *OptionsEntryModel) addTableRow(content *strings.Builder, labelStyle, valueStyle lipgloss.Style, label, value string) {
	content.WriteString(labelStyle.Render(label))
	content.WriteString("  ")
	content.WriteString(valueStyle.Render(value))
	content.WriteString("\n")
}
