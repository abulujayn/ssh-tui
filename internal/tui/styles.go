package tui

import "github.com/charmbracelet/lipgloss"

// InstructionStyle returns the shared style used for footer instructions.
func InstructionStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true)
}
