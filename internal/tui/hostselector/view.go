package hostselector

import (
	"fmt"
	"ssh-tui/internal/parser"
	"ssh-tui/internal/types"
	"strings"

	"ssh-tui/internal/tui/helpers"
	"ssh-tui/internal/tui/ui"

	"github.com/charmbracelet/lipgloss"
)

// View implements the tea.Model interface
func (m *HostSelectorModel) View() string {
	var b strings.Builder

	b.WriteString(ui.TitleStyle.Render("Host selection") + "\n\n")

	renderedSearch := helpers.RenderInputWithCursor(m.searchInput, len(m.searchInput), 40)
	b.WriteString(ui.SearchStyle.Render("Search: "+renderedSearch) + "\n\n")

	// If no hosts in the filtered list, show helpful messages and return early
	if len(m.filteredHosts) == 0 {
		if m.searchInput != "" {
			if parser.IsValidHost(m.searchInput) {
				b.WriteString(ui.TitleStyle.Render("Press Enter to connect to custom host: " + m.searchInput))
			} else {
				b.WriteString(ui.ErrorStyle.Render("No hosts found matching: " + m.searchInput))
			}
		} else {
			b.WriteString(ui.ErrorStyle.Render("No SSH hosts found.\nCheck that ~/.ssh/config or ~/.ssh/known_hosts exist and contain host entries."))
		}

		b.WriteString("\n\n")
		b.WriteString(ui.InstructionStyle.Render(ui.InstructionNav))
		return b.String()
	}

	// Calculate visible range for scrolling
	linesPerHost := 1
	maxVisible := (m.height - 8) / linesPerHost
	if maxVisible < 3 {
		maxVisible = 3
	}

	start, end := helpers.ScrollRange(len(m.filteredHosts), m.cursor, maxVisible)

	// Render visible hosts
	for i := start; i < end; i++ {
		host := m.filteredHosts[i]
		hostDisplay := parser.FormatHostDisplay(host)
		lines := strings.Split(hostDisplay, "\n")

		if i == m.cursor {
			var content strings.Builder
			styledHostLine := m.formatHostLineWithAliasesSelectedEnhanced(host, ui.SelectedTextStyle)
			content.WriteString(styledHostLine)
			for j := 1; j < len(lines); j++ {
				content.WriteString("\n" + ui.DetailTextStyle.Render(lines[j]))
			}
			b.WriteString(ui.SelectedContainerStyle.Render(content.String()) + "\n")
		} else {
			styledHostLine := m.formatHostLineWithAliases(host, ui.NormalStyle, ui.DetailTextStyle)
			b.WriteString(ui.NormalContainerStyle.Render(styledHostLine) + "\n")
		}
	}

	// Scroll indicator
	if len(m.filteredHosts) > maxVisible {
		scrollInfo := fmt.Sprintf("\n%d/%d hosts", m.cursor+1, len(m.filteredHosts))
		if start > 0 || end < len(m.filteredHosts) {
			scrollInfo += " (scroll with \u2191/\u2193)"
		}
		b.WriteString(ui.InstructionStyle.Render(scrollInfo))
	}

	b.WriteString("\n\n")
	b.WriteString(ui.InstructionStyle.Render(ui.InstructionNav))

	return b.String()
}

// formatHostLineWithAliases formats the host name line with styled aliases
func (m *HostSelectorModel) formatHostLineWithAliases(host types.SSHHost, normalStyle, aliasStyle lipgloss.Style) string {
	hostName := normalStyle.Render(host.Name)

	if len(host.Aliases) > 0 {
		aliasesStr := " [" + strings.Join(host.Aliases, ", ") + "]"
		hostName += aliasStyle.Render(aliasesStr)
	}

	return hostName
}

// formatHostLineWithAliasesSelectedEnhanced formats the host name line for enhanced selected state
func (m *HostSelectorModel) formatHostLineWithAliasesSelectedEnhanced(host types.SSHHost, selectedStyle lipgloss.Style) string {
	if len(host.Aliases) > 0 {
		// For enhanced selected items, use accent color for aliases
		aliasStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("183"))

		hostName := selectedStyle.Render(host.Name)
		aliasesStr := " [" + strings.Join(host.Aliases, ", ") + "]"
		hostName += aliasStyle.Render(aliasesStr)
		return hostName
	}

	return selectedStyle.Render(host.Name)
}
