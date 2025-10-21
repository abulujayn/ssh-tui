package optionsentry

import (
	"ssh-tui/internal/tui/helpers"

	tea "github.com/charmbracelet/bubbletea"
)

// Update implements the tea.Model interface for options entry
func (m *OptionsEntryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			m.cancelled = true
			return m, tea.Quit

		case "enter":
			m.confirmed = true
			return m, tea.Quit

		case "left", "ctrl+b":
			if m.cursor > 0 {
				m.cursor--
			}

		case "right", "ctrl+f":
			if m.cursor < len(m.options) {
				m.cursor++
			}

		case "home", "ctrl+a":
			m.cursor = 0

		case "end", "ctrl+e":
			m.cursor = len(m.options)

		case "backspace", "ctrl+h":
			if m.cursor > 0 && len(m.options) > 0 {
				m.options = m.options[:m.cursor-1] + m.options[m.cursor:]
				m.cursor--
			}

		case "delete", "ctrl+d":
			if m.cursor < len(m.options) {
				m.options = m.options[:m.cursor] + m.options[m.cursor+1:]
			}

		case "ctrl+u":
			// Delete from cursor to beginning
			m.options = m.options[m.cursor:]
			m.cursor = 0

		case "ctrl+k":
			// Delete from cursor to end
			m.options = m.options[:m.cursor]

		case "ctrl+w":
			// Delete word backwards
			if m.cursor > 0 {
				m.options, m.cursor = helpers.DeleteWordBackwards(m.options, m.cursor)
			}

		default:
			// Handle regular character input
			if len(msg.String()) == 1 {
				char := msg.String()
				m.options = m.options[:m.cursor] + char + m.options[m.cursor:]
				m.cursor++
			}
		}
	}

	return m, nil
}
