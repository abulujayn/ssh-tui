package hostselector

import (
	"ssh-tui/internal/parser"
	"ssh-tui/internal/tui/helpers"

	tea "github.com/charmbracelet/bubbletea"
)

// Update implements the tea.Model interface for host selector
func (m *HostSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			// If there's search input, clear it; otherwise quit the app
			if m.searchInput != "" {
				m.searchInput = ""
				m.filteredHosts = m.hosts
				m.cursor = 0
			} else {
				return m, tea.Quit
			}

		case "enter":
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
					ch := helpers.BuildCustomHost(m.searchInput)
					m.selectedHost = &ch
					m.selected = true
					return m, tea.Quit
				}
			}

		case "o":
			// Treat 'o' as input to search (search is always active)
			m.searchInput += "o"
			m.updateFilter()
			break

		case "tab":
			// Open options for the selected host when possible (works while searching)
			if len(m.filteredHosts) > 0 && m.cursor < len(m.filteredHosts) {
				m.selectedHost = &m.filteredHosts[m.cursor]
				m.selected = true
				m.openOptions = true
				return m, tea.Quit
			}

			// If no filtered hosts, but the user typed a valid custom host, open options
			if len(m.filteredHosts) == 0 && m.searchInput != "" {
				if parser.IsValidHost(m.searchInput) {
					ch := helpers.BuildCustomHost(m.searchInput)
					m.selectedHost = &ch
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

		case "down":
			// Arrow down navigates within the current filtered list
			if m.cursor < len(m.filteredHosts)-1 {
				m.cursor++
			}

		case "backspace":
			if len(m.searchInput) > 0 {
				m.searchInput = m.searchInput[:len(m.searchInput)-1]
				m.updateFilter()
			}

		default:
			// If it's a single printable character, treat it as typing input.
			// If search is not yet active, enable it and start the search with this char.
			if len(msg.String()) == 1 {
				ch := msg.String()
				// Always treat printable characters as search input
				m.searchInput += ch
				m.updateFilter()
				m.cursor = 0
			}
		}
	}

	return m, nil
}
