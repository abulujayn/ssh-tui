package tui

import (
	"fmt"
	"ssh-tui/internal/parser"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HostSelectorModel represents the host selection screen
type HostSelectorModel struct {
	hosts         []parser.SSHHost
	filteredHosts []parser.SSHHost
	cursor        int
	searchInput   string
	searchActive  bool
	selected      bool
	selectedHost  *parser.SSHHost
	width         int
	height        int
}

// NewHostSelectorModel creates a new host selector model
func NewHostSelectorModel(hosts []parser.SSHHost) *HostSelectorModel {
	return &HostSelectorModel{
		hosts:         hosts,
		filteredHosts: hosts,
		cursor:        0,
		searchActive:  false,
		selected:      false,
	}
}

// Init implements the tea.Model interface
func (m *HostSelectorModel) Init() tea.Cmd {
	return nil
}

// Update implements the tea.Model interface
func (m *HostSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "/":
			// Toggle search mode
			m.searchActive = !m.searchActive
			if m.searchActive {
				m.cursor = 0
			}

		case "esc":
			// Exit search mode
			if m.searchActive {
				m.searchActive = false
				m.searchInput = ""
				m.filteredHosts = m.hosts
				m.cursor = 0
			} else {
				return m, tea.Quit
			}

		case "enter":
			if len(m.filteredHosts) > 0 && m.cursor < len(m.filteredHosts) {
				m.selectedHost = &m.filteredHosts[m.cursor]
				m.selected = true
				return m, tea.Quit
			}

		case "up", "k":
			if !m.searchActive && m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if !m.searchActive && m.cursor < len(m.filteredHosts)-1 {
				m.cursor++
			}

		case "backspace":
			if m.searchActive && len(m.searchInput) > 0 {
				m.searchInput = m.searchInput[:len(m.searchInput)-1]
				m.updateFilter()
			}

		default:
			if m.searchActive && len(msg.String()) == 1 {
				m.searchInput += msg.String()
				m.updateFilter()
			}
		}
	}

	return m, nil
}

// View implements the tea.Model interface
func (m *HostSelectorModel) View() string {
	var b strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Render("SSH Host Selector")

	b.WriteString(title + "\n\n")

	// Search bar
	searchStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86"))

	if m.searchActive {
		b.WriteString(searchStyle.Render("🔍 Search: "+m.searchInput+"_") + "\n\n")
	} else {
		if m.searchInput != "" {
			b.WriteString(searchStyle.Render("🔍 Search: "+m.searchInput+" (press / to edit)") + "\n\n")
		} else {
			b.WriteString(searchStyle.Render("🔍 Press / to search") + "\n\n")
		}
	}

	// Host list
	if len(m.filteredHosts) == 0 {
		noHostsStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("203")).
			Bold(true)

		if m.searchInput != "" {
			b.WriteString(noHostsStyle.Render("No hosts found matching: " + m.searchInput))
		} else {
			b.WriteString(noHostsStyle.Render("No SSH hosts found.\nCheck that ~/.ssh/config or ~/.ssh/known_hosts exist and contain host entries."))
		}

		// Instructions at the bottom even when no hosts found
		b.WriteString("\n\n")
		instructions := "Use ↑/↓ or j/k to navigate, / to search, Enter to select, q or Ctrl+C to quit"
		if m.searchActive {
			instructions = "Type to search, Esc to exit search, Enter to select"
		}

		instructionStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)

		b.WriteString(instructionStyle.Render(instructions))
		return b.String()
	}

	// Calculate visible range for scrolling
	maxVisible := m.height - 8 // Account for header, search, and bottom instructions
	if maxVisible < 5 {
		maxVisible = 5
	}

	start := 0
	end := len(m.filteredHosts)

	if len(m.filteredHosts) > maxVisible {
		// Calculate scroll position
		if m.cursor >= maxVisible/2 {
			start = m.cursor - maxVisible/2
			if start > len(m.filteredHosts)-maxVisible {
				start = len(m.filteredHosts) - maxVisible
			}
		}
		end = start + maxVisible
		if end > len(m.filteredHosts) {
			end = len(m.filteredHosts)
		}
	}

	// Render visible hosts
	for i := start; i < end; i++ {
		host := m.filteredHosts[i]
		hostDisplay := parser.FormatHostDisplay(host)

		if i == m.cursor {
			// Selected style
			selectedStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).
				Background(lipgloss.Color("86")).
				Bold(true)

			b.WriteString("▶ " + selectedStyle.Render(hostDisplay) + "\n")
		} else {
			// Normal style
			normalStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

			b.WriteString("  " + normalStyle.Render(hostDisplay) + "\n")
		}
	}

	// Show scroll indicator if needed
	if len(m.filteredHosts) > maxVisible {
		scrollInfo := fmt.Sprintf("\n%d/%d hosts", m.cursor+1, len(m.filteredHosts))
		if start > 0 || end < len(m.filteredHosts) {
			scrollInfo += " (scroll with ↑/↓)"
		}

		scrollStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)

		b.WriteString(scrollStyle.Render(scrollInfo))
	}

	// Instructions at the bottom
	b.WriteString("\n\n")
	instructions := "Use ↑/↓ or j/k to navigate, / to search, Enter to select, q or Ctrl+C to quit"
	if m.searchActive {
		instructions = "Type to search, Esc to exit search, Enter to select"
	}

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true)

	b.WriteString(instructionStyle.Render(instructions))

	return b.String()
}

// updateFilter updates the filtered hosts based on search input
func (m *HostSelectorModel) updateFilter() {
	m.filteredHosts = parser.FilterHosts(m.hosts, m.searchInput)

	// Reset cursor if it's out of bounds
	if m.cursor >= len(m.filteredHosts) {
		m.cursor = 0
	}
	if len(m.filteredHosts) > 0 && m.cursor < 0 {
		m.cursor = 0
	}
}

// GetSelectedHost returns the selected host
func (m *HostSelectorModel) GetSelectedHost() *parser.SSHHost {
	if m.selected {
		return m.selectedHost
	}
	return nil
}

// IsSelected returns whether a host was selected
func (m *HostSelectorModel) IsSelected() bool {
	return m.selected
}
