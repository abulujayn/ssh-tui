package hostselector

import (
	"strings"
	"testing"

	"ssh-tui/internal/types"

	tea "github.com/charmbracelet/bubbletea"
)

func TestHostSelectorModel_Update(t *testing.T) {
	// Create test hosts
	hosts := []types.SSHHost{
		{Name: "host1", HostName: "host1.example.com", User: "user1", Port: "22", Source: types.SourceConfig},
		{Name: "host2", HostName: "host2.example.com", User: "user2", Port: "22", Source: types.SourceConfig},
		{Name: "testhost", HostName: "test.example.com", User: "test", Port: "22", Source: types.SourceKnownHosts},
	}

	model := NewHostSelectorModel(hosts)

	// Test typing 'test' to filter hosts
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	hsModel := updatedModel.(*HostSelectorModel)

	if hsModel.searchInput != "test" {
		t.Errorf("Expected searchInput to be 'test', got %q", hsModel.searchInput)
	}

	if len(hsModel.filteredHosts) != 1 || hsModel.filteredHosts[0].Name != "testhost" {
		t.Errorf("Expected 1 filtered host 'testhost', got %d hosts: %v", len(hsModel.filteredHosts), hsModel.filteredHosts)
	}

	// Test cursor down
	updatedModel, _ = hsModel.Update(tea.KeyMsg{Type: tea.KeyDown})
	hsModel = updatedModel.(*HostSelectorModel)

	// Since only 1 host, cursor should stay at 0
	if hsModel.cursor != 0 {
		t.Errorf("Expected cursor to be 0, got %d", hsModel.cursor)
	}

	// Test enter to select
	updatedModel, _ = hsModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	hsModel = updatedModel.(*HostSelectorModel)

	if !hsModel.selected || hsModel.selectedHost == nil || hsModel.selectedHost.Name != "testhost" {
		t.Errorf("Expected host 'testhost' to be selected")
	}
}

func TestHostSelectorModel_View(t *testing.T) {
	hosts := []types.SSHHost{
		{Name: "host1", HostName: "host1.example.com", User: "user1", Port: "22", Source: types.SourceConfig},
	}

	model := NewHostSelectorModel(hosts)
	model.width = 80
	model.height = 24

	view := model.View()

	// Check that view contains expected elements
	if !strings.Contains(view, "Host selection") {
		t.Errorf("View should contain 'Host selection' title")
	}

	if !strings.Contains(view, "Search:") {
		t.Errorf("View should contain search input")
	}

	if !strings.Contains(view, "host1") {
		t.Errorf("View should contain host name")
	}
}

func TestHostSelectorModel_CustomHost(t *testing.T) {
	model := NewHostSelectorModel([]types.SSHHost{}) // Empty hosts list

	// Type a custom host
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'@'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'.'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

	hsModel := updatedModel.(*HostSelectorModel)

	if hsModel.searchInput != "user@example.com" {
		t.Errorf("Expected searchInput to be 'user@example.com', got %q", hsModel.searchInput)
	}

	// Select the custom host
	updatedModel, _ = hsModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	hsModel = updatedModel.(*HostSelectorModel)

	if !hsModel.selected || hsModel.selectedHost == nil || hsModel.selectedHost.Name != "user@example.com" {
		t.Errorf("Expected custom host 'user@example.com' to be selected")
	}
}
