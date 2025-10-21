package hostselector

import (
	"ssh-tui/internal/parser"

	tea "github.com/charmbracelet/bubbletea"
)

// HostSelectorModel represents the host selection screen
type HostSelectorModel struct {
	hosts         []parser.SSHHost
	filteredHosts []parser.SSHHost
	cursor        int
	searchInput   string
	selected      bool
	selectedHost  *parser.SSHHost
	openOptions   bool
	width         int
	height        int
}

// NewHostSelectorModel creates a new host selector model
func NewHostSelectorModel(hosts []parser.SSHHost) *HostSelectorModel {
	return &HostSelectorModel{
		hosts:         hosts,
		filteredHosts: hosts,
		cursor:        0,
		selected:      false,
	}
}

// Init implements the tea.Model interface
func (m *HostSelectorModel) Init() tea.Cmd {
	return nil
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
