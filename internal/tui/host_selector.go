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
	// If true, user requested to open the options screen after selection.
	openOptions bool
	width       int
	height      int
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
		case "ctrl+c":
			return m, tea.Quit

		case "/":
			if !m.searchActive {
				// Toggle search mode
				m.searchActive = true
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
			if m.searchActive {
				// If there are filtered hosts, select the focused one
				if len(m.filteredHosts) > 0 && m.cursor < len(m.filteredHosts) {
					m.selectedHost = &m.filteredHosts[m.cursor]
					m.selected = true
					m.openOptions = false // Execute straight away by default on Enter
					return m, tea.Quit
				}
				// If no hosts found, treat search input as custom host if valid
				if m.searchInput != "" && len(m.filteredHosts) == 0 {
					if parser.IsValidHost(m.searchInput) {
						// Parse user and host from input
						user, host := parser.ParseUserHost(m.searchInput)
						customHost := parser.SSHHost{
							Name:     m.searchInput,
							HostName: host,
							User:     user,
							Port:     "22",
							Source:   "custom",
							Aliases:  nil,
						}
						m.selectedHost = &customHost
						m.selected = true
						return m, tea.Quit
					}
				}
			} else if len(m.filteredHosts) > 0 && m.cursor < len(m.filteredHosts) {
				// Not in search mode: select the focused host
				m.selectedHost = &m.filteredHosts[m.cursor]
				m.selected = true
				m.openOptions = false // Execute straight away by default on Enter
				return m, tea.Quit
			}

		case "o":
			// If search is active, treat 'o' as input (so users can type the letter 'o')
			if m.searchActive {
				m.searchInput += "o"
				m.updateFilter()
				break
			}

			// Open options for the selected host when not searching
			if len(m.filteredHosts) > 0 && m.cursor < len(m.filteredHosts) {
				m.selectedHost = &m.filteredHosts[m.cursor]
				m.selected = true
				m.openOptions = true
				return m, tea.Quit
			}

		case "tab":
			// Open options for the selected host when possible (works during search)
			if len(m.filteredHosts) > 0 && m.cursor < len(m.filteredHosts) {
				m.selectedHost = &m.filteredHosts[m.cursor]
				m.selected = true
				m.openOptions = true
				return m, tea.Quit
			}

			// If searching and no filtered hosts, but the user typed a valid custom host,
			// treat Tab as request to open options for that custom host.
			if m.searchActive && len(m.filteredHosts) == 0 && m.searchInput != "" {
				if parser.IsValidHost(m.searchInput) {
					user, host := parser.ParseUserHost(m.searchInput)
					customHost := parser.SSHHost{
						Name:     m.searchInput,
						HostName: host,
						User:     user,
						Port:     "22",
						Source:   "custom",
						Aliases:  nil,
					}
					m.selectedHost = &customHost
					m.selected = true
					m.openOptions = true
					return m, tea.Quit
				}
			}

		case "up":
			// Arrow up navigates within the current filtered list
			if m.cursor > 0 {
				m.cursor--
			}
		case "k":
			if m.searchActive {
				// treat as input
				m.searchInput += "k"
				m.updateFilter()
			} else if m.cursor > 0 {
				// 'k' navigation when not searching
				m.cursor--
			}

		case "down":
			// Arrow down navigates within the current filtered list
			if m.cursor < len(m.filteredHosts)-1 {
				m.cursor++
			}
		case "j":
			if m.searchActive {
				// treat as input
				m.searchInput += "j"
				m.updateFilter()
			} else if m.cursor < len(m.filteredHosts)-1 {
				// 'j' navigation when not searching
				m.cursor++
			}

		case "backspace":
			if m.searchActive && len(m.searchInput) > 0 {
				m.searchInput = m.searchInput[:len(m.searchInput)-1]
				m.updateFilter()
			}

		case "q":
			if m.searchActive {
				// treat as input
				m.searchInput += "q"
				m.updateFilter()
			} else {
				return m, tea.Quit
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
		Foreground(lipgloss.Color("183")).
		Bold(true).
		Render("Host selection")

	b.WriteString(title + "\n\n")

	// Search bar
	searchStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86"))

	if m.searchActive {
		b.WriteString(searchStyle.Render("Search: "+m.searchInput+"_") + "\n\n")
	} else {
		if m.searchInput != "" {
			b.WriteString(searchStyle.Render("Search: "+m.searchInput+" (press / to edit)") + "\n\n")
		} else {
			b.WriteString(searchStyle.Render("Press / to search") + "\n\n")
		}
	}

	// Host list
	if len(m.filteredHosts) == 0 {
		// Styles
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("203")). // red
			Bold(true)
		titleStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("183")). // purple (same as title)
			Bold(true)

		if m.searchInput != "" {
			if parser.IsValidHost(m.searchInput) {
				b.WriteString(titleStyle.Render("Press Enter to connect to custom host: " + m.searchInput))
			} else {
				b.WriteString(errorStyle.Render("No hosts found matching: " + m.searchInput))
			}
		} else {
			b.WriteString(errorStyle.Render("No SSH hosts found.\nCheck that ~/.ssh/config or ~/.ssh/known_hosts exist and contain host entries."))
		}

		// Instructions at the bottom even when no hosts found
		b.WriteString("\n\n")
		instructions := "Use ↑/↓ or j/k to navigate, / to search, Tab for options, Enter to connect, q to quit"
		if m.searchActive {
			instructions = "Type to search, Esc to exit search, Tab for options, Enter to connect"
		}

		instructionStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)

		b.WriteString(instructionStyle.Render(instructions))
		return b.String()
	}

	// Calculate visible range for scrolling - account for single-line hosts (1 line each)
	// Only focused host shows details, others show just the name
	// Use 1 line per host to reduce vertical spacing between entries
	linesPerHost := 1
	maxVisible := (m.height - 8) / linesPerHost // Account for header, search, and bottom instructions
	if maxVisible < 3 {
		maxVisible = 3
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

		// Split the host display into lines for proper styling
		lines := strings.Split(hostDisplay, "\n")

		if i == m.cursor {
			// Selected style - show all details with enhanced visual styling
			selectedContainerStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("86")).
				Foreground(lipgloss.Color("252")).
				Padding(0, 1, 0, 2).
				Margin(0, 2, 0, 0).
				Bold(true)

			selectedTextStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("86")).
				Bold(true)

			// Detail text style for host information
			detailTextStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("245"))

			// Create content for the styled container
			var content strings.Builder

			// Apply selected style to first line with cursor (with styled aliases)
			styledHostLine := m.formatHostLineWithAliasesSelectedEnhanced(host, selectedTextStyle)
			content.WriteString(styledHostLine)

			// Apply detail style to additional lines with subtle styling
			for j := 1; j < len(lines); j++ {
				content.WriteString("\n" + detailTextStyle.Render(lines[j]))
			}

			// Add a subtle hint for opening options below the server info
			content.WriteString("\n" + InstructionStyle().Render("Tab for options"))

			// Render the entire selection in a styled container
			b.WriteString(selectedContainerStyle.Render(content.String()) + "\n")
		} else {
			// Normal style - only show host name, no details
			normalStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

			detailStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("241"))

			// Add subtle padding for better alignment with focused entries
			normalContainerStyle := lipgloss.NewStyle().
				Padding(0, 0, 0, 3)

			// Render only the first line (host name) with styled aliases
			styledHostLine := m.formatHostLineWithAliases(host, normalStyle, detailStyle)
			b.WriteString(normalContainerStyle.Render(styledHostLine) + "\n")

			// Skip rendering additional lines (details) for non-focused hosts
		}

		// Intentionally avoid adding extra blank lines between hosts to keep list compact
	}

	// Show scroll indicator if needed
	if len(m.filteredHosts) > maxVisible {
		scrollInfo := fmt.Sprintf("\n%d/%d hosts", m.cursor+1, len(m.filteredHosts))
		if start > 0 || end < len(m.filteredHosts) {
			scrollInfo += " (scroll with ↑/↓)"
		}

		b.WriteString(InstructionStyle().Render(scrollInfo))
	}

	// Instructions at the bottom
	b.WriteString("\n\n")
	instructions := "Use ↑/↓ or j/k to navigate, / to search, Enter to connect, q to quit"
	if m.searchActive {
		instructions = "Type to search, Esc to exit search, Enter to connect"
	}

	b.WriteString(InstructionStyle().Render(instructions))

	return b.String()
}

// updateFilter updates the filtered hosts based on search input
func (m *HostSelectorModel) updateFilter() {
	m.filteredHosts = parser.FilterHosts(m.hosts, m.searchInput)

	// Whenever the filter changes (search input modified), reset focus to the first entry
	m.cursor = 0
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

// OpenOptionsRequested indicates whether the user asked to open options after selection
func (m *HostSelectorModel) OpenOptionsRequested() bool {
	return m.openOptions
}

// formatHostLineWithAliases formats the host name line with styled aliases
func (m *HostSelectorModel) formatHostLineWithAliases(host parser.SSHHost, normalStyle, aliasStyle lipgloss.Style) string {
	hostName := normalStyle.Render(host.Name)

	if len(host.Aliases) > 0 {
		aliasesStr := " [" + strings.Join(host.Aliases, ", ") + "]"
		hostName += aliasStyle.Render(aliasesStr)
	}

	return hostName
}

// formatHostLineWithAliasesSelectedEnhanced formats the host name line for enhanced selected state
func (m *HostSelectorModel) formatHostLineWithAliasesSelectedEnhanced(host parser.SSHHost, selectedStyle lipgloss.Style) string {
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
