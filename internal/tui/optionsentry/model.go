package optionsentry

import (
	"ssh-tui/internal/ssh"
	"ssh-tui/internal/types"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// OptionsEntryModel represents the SSH options entry screen
type OptionsEntryModel struct {
	host      *types.SSHHost
	options   string
	cursor    int
	confirmed bool
	cancelled bool
	width     int
	height    int
}

// NewOptionsEntryModel creates a new options entry model
func NewOptionsEntryModel(host *types.SSHHost) *OptionsEntryModel {
	return &OptionsEntryModel{
		host:      host,
		options:   "",
		cursor:    0,
		confirmed: false,
		cancelled: false,
	}
}

// Init implements the tea.Model interface
func (m *OptionsEntryModel) Init() tea.Cmd {
	return nil
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

// GetCommand returns the SSH command that would be executed with current options
func (m *OptionsEntryModel) GetCommand() string {
	return ssh.BuildSSHCommand(m.host, m.options)
}
