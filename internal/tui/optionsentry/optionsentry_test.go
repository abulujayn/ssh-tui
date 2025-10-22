package optionsentry

import (
	"strings"
	"testing"

	"ssh-tui/internal/types"

	tea "github.com/charmbracelet/bubbletea"
)

func TestOptionsEntryModel_Update(t *testing.T) {
	// Create test host
	host := &types.SSHHost{
		Name:     "testhost",
		HostName: "test.example.com",
		User:     "testuser",
		Port:     "22",
		Source:   types.SourceConfig,
	}

	model := NewOptionsEntryModel(host)

	// Test typing options
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'-'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	optModel := updatedModel.(*OptionsEntryModel)

	if optModel.options != "-p 2222" {
		t.Errorf("Expected options to be '-p 2222', got %q", optModel.options)
	}

	if optModel.cursor != 7 {
		t.Errorf("Expected cursor to be at position 7, got %d", optModel.cursor)
	}

	// Test cursor movement left
	updatedModel, _ = optModel.Update(tea.KeyMsg{Type: tea.KeyLeft})
	optModel = updatedModel.(*OptionsEntryModel)

	if optModel.cursor != 6 {
		t.Errorf("Expected cursor to be at position 6 after left, got %d", optModel.cursor)
	}

	// Test backspace
	updatedModel, _ = optModel.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	optModel = updatedModel.(*OptionsEntryModel)

	if optModel.options != "-p 222" {
		t.Errorf("Expected options to be '-p 222' after backspace, got %q", optModel.options)
	}

	if optModel.cursor != 5 {
		t.Errorf("Expected cursor to be at position 5 after backspace, got %d", optModel.cursor)
	}

	// Test enter to confirm
	updatedModel, _ = optModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	optModel = updatedModel.(*OptionsEntryModel)

	if !optModel.IsConfirmed() {
		t.Errorf("Expected options to be confirmed")
	}

	if optModel.GetOptions() != "-p 222" {
		t.Errorf("Expected trimmed options to be '-p 222', got %q", optModel.GetOptions())
	}
}

func TestOptionsEntryModel_CursorNavigation(t *testing.T) {
	host := &types.SSHHost{Name: "host", HostName: "host.com", Source: types.SourceConfig}
	model := NewOptionsEntryModel(host)

	// Type some text
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
	optModel := updatedModel.(*OptionsEntryModel)

	// Test home key
	updatedModel, _ = optModel.Update(tea.KeyMsg{Type: tea.KeyHome})
	optModel = updatedModel.(*OptionsEntryModel)

	if optModel.cursor != 0 {
		t.Errorf("Expected cursor at home (0), got %d", optModel.cursor)
	}

	// Test end key
	updatedModel, _ = optModel.Update(tea.KeyMsg{Type: tea.KeyEnd})
	optModel = updatedModel.(*OptionsEntryModel)

	if optModel.cursor != 5 {
		t.Errorf("Expected cursor at end (5), got %d", optModel.cursor)
	}

	// Test ctrl+u (delete to beginning)
	updatedModel, _ = optModel.Update(tea.KeyMsg{Type: tea.KeyCtrlU})
	optModel = updatedModel.(*OptionsEntryModel)

	if optModel.options != "" {
		t.Errorf("Expected options to be empty after ctrl+u, got %q", optModel.options)
	}

	if optModel.cursor != 0 {
		t.Errorf("Expected cursor at 0 after ctrl+u, got %d", optModel.cursor)
	}
}

func TestOptionsEntryModel_Cancellation(t *testing.T) {
	host := &types.SSHHost{Name: "host", HostName: "host.com", Source: types.SourceConfig}
	model := NewOptionsEntryModel(host)

	// Type some options
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'-'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	optModel := updatedModel.(*OptionsEntryModel)

	// Test escape to cancel
	updatedModel, _ = optModel.Update(tea.KeyMsg{Type: tea.KeyEsc})
	optModel = updatedModel.(*OptionsEntryModel)

	if !optModel.IsCancelled() {
		t.Errorf("Expected options entry to be cancelled")
	}

	if optModel.IsConfirmed() {
		t.Errorf("Expected options entry not to be confirmed")
	}
}

func TestOptionsEntryModel_View(t *testing.T) {
	host := &types.SSHHost{
		Name:     "testhost",
		HostName: "test.example.com",
		User:     "testuser",
		Port:     "2222",
		Source:   types.SourceConfig,
	}

	model := NewOptionsEntryModel(host)
	model.width = 80
	model.height = 24

	// Type some options
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'-'}})
	updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	model = updatedModel.(*OptionsEntryModel)

	view := model.View()

	// Check that view contains expected elements
	if !strings.Contains(view, "SSH Options & Arguments") {
		t.Errorf("View should contain title")
	}

	if !strings.Contains(view, "Host configuration") {
		t.Errorf("View should contain host configuration table")
	}

	if !strings.Contains(view, "test.example.com") {
		t.Errorf("View should contain hostname")
	}

	if !strings.Contains(view, "testuser") {
		t.Errorf("View should contain username")
	}

	if !strings.Contains(view, "2222") {
		t.Errorf("View should contain port")
	}

	if !strings.Contains(view, "Options:") {
		t.Errorf("View should contain options label")
	}

	if !strings.Contains(view, "-v") {
		t.Errorf("View should contain entered options")
	}

	if !strings.Contains(view, "Command Preview:") {
		t.Errorf("View should contain command preview")
	}

	if !strings.Contains(view, "ssh testhost -v") {
		t.Errorf("View should contain command preview with options")
	}

	if !strings.Contains(view, "Use Enter to execute, Esc to go back") {
		t.Errorf("View should contain instructions")
	}
}

func TestOptionsEntryModel_CustomHost(t *testing.T) {
	// Test with a custom host (no config table)
	host := &types.SSHHost{
		Name:     "user@custom.com",
		HostName: "custom.com",
		User:     "user",
		Port:     "22",
		Source:   types.SourceCustom,
	}

	model := NewOptionsEntryModel(host)
	model.width = 80
	model.height = 24

	view := model.View()

	// For custom hosts, should not show host configuration table
	if strings.Contains(view, "Host configuration") {
		t.Errorf("View should not contain host configuration table for custom hosts")
	}

	if !strings.Contains(view, "user@custom.com") {
		t.Errorf("View should contain custom host name")
	}
}
