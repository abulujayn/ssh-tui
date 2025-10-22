package ui

import "github.com/charmbracelet/lipgloss"

// Shared UI text constants to avoid duplication across views.
const (
	InstructionNav = "Use \u2191/\u2193 to navigate, Tab for options, Enter to connect"
	TabForOptions  = "Tab for options"
	ExamplesText   = "Examples: -L 8080:localhost:80 -i ~/.ssh/id_rsa -p 2222 -X"
	SearchLabel    = "Search: "
)

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
