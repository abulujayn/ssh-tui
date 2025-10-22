package helpers

import (
	"ssh-tui/internal/parser"
	"ssh-tui/internal/types"

	"github.com/charmbracelet/lipgloss"
)

// BuildCustomHost builds a parser.SSHHost from raw user input (e.g. user@host)
func BuildCustomHost(raw string) types.SSHHost {
	user, host := parser.ParseUserHost(raw)
	return types.SSHHost{
		Name:     raw,
		HostName: host,
		User:     user,
		Port:     types.DefaultSSHPort,
		Source:   types.SourceCustom,
		Aliases:  nil,
	}
}

// ScrollRange calculates start and end indices for a list given a desired max visible count
func ScrollRange(listLen, cursor, maxVisible int) (start, end int) {
	if maxVisible <= 0 {
		return 0, listLen
	}
	if listLen <= maxVisible {
		return 0, listLen
	}
	start = 0
	if cursor >= maxVisible/2 {
		start = cursor - maxVisible/2
		if start > listLen-maxVisible {
			start = listLen - maxVisible
		}
	}
	end = start + maxVisible
	if end > listLen {
		end = listLen
	}
	return start, end
}

// DeleteWordBackwards deletes the word immediately before the cursor in s and returns the new string and new cursor position.
func DeleteWordBackwards(s string, cursor int) (string, int) {
	if cursor <= 0 || len(s) == 0 {
		return s, cursor
	}
	if cursor > len(s) {
		cursor = len(s)
	}
	// find the start of the previous word
	i := cursor - 1
	for i > 0 && s[i-1] == ' ' {
		i--
	}
	for i > 0 && s[i-1] != ' ' {
		i--
	}
	newS := s[:i] + s[cursor:]
	return newS, i
}

// RenderInputWithCursor returns a string where the cursor position is rendered
func RenderInputWithCursor(s string, cursor int, width int) string {
	// cursorStyle: render the cursor with the same color as the input border so it matches visually
	cursorStyle := lipgloss.NewStyle().Background(lipgloss.Color("86")).Foreground(lipgloss.Color("0")).Bold(true)

	// Normalize cursor bounds
	if cursor < 0 {
		cursor = 0
	}
	if cursor > len(s) {
		cursor = len(s)
	}

	left := ""
	right := ""

	if len(s) > 0 {
		left = s[:cursor]
		if cursor < len(s) {
			right = s[cursor+1:]
		}
	}

	var cursorGlyph string
	if len(s) == 0 {
		// empty input — show a reversed space as cursor
		cursorGlyph = cursorStyle.Render(" ")
	} else if cursor == len(s) {
		// cursor at end — show reversed space after input
		cursorGlyph = cursorStyle.Render(" ")
	} else {
		// cursor on an existing character — render that character reversed
		cursorGlyph = cursorStyle.Render(string(s[cursor]))
	}

	return left + cursorGlyph + right
}
