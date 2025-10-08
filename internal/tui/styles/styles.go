package styles

import "github.com/charmbracelet/lipgloss"

// Shared styles used across TUI models. Exported so other files can reference them.
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("183")).
			Bold(true)

	SearchStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("203")).
			Bold(true)

	InstructionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Italic(true)

	SelectedContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("86")).
				Foreground(lipgloss.Color("252")).
				Padding(0, 1, 0, 2).
				Margin(0, 2, 0, 0).
				Bold(true)

	SelectedTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("86")).
				Bold(true)

	DetailTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	NormalContainerStyle = lipgloss.NewStyle().
				Padding(0, 0, 0, 3)
)
