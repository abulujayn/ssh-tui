package tui

import (
	"ssh-tui/internal/parser"
	"ssh-tui/internal/tui/helpers"
	"testing"
)

func TestBuildCustomHost(t *testing.T) {
	cases := []struct {
		in   string
		want parser.SSHHost
	}{
		{"user@example.com", parser.SSHHost{Name: "user@example.com", HostName: "example.com", User: "user", Port: "22", Source: "custom"}},
		{"example.org", parser.SSHHost{Name: "example.org", HostName: "example.org", User: "", Port: "22", Source: "custom"}},
	}

	for _, c := range cases {
		got := helpers.BuildCustomHost(c.in)
		// compare relevant fields
		if got.Name != c.want.Name || got.HostName != c.want.HostName || got.User != c.want.User || got.Port != c.want.Port || got.Source != c.want.Source {
			t.Fatalf("BuildCustomHost(%q) = %+v, want %+v", c.in, got, c.want)
		}
	}
}

func TestScrollRange(t *testing.T) {
	cases := []struct {
		listLen, cursor, maxVisible int
		wantStart, wantEnd          int
	}{
		{0, 0, 3, 0, 0},
		{2, 0, 5, 0, 2},
		{10, 0, 3, 0, 3},
		{10, 1, 3, 0, 3},
		{10, 2, 3, 1, 4},
		{10, 9, 3, 7, 10},
	}

	for _, c := range cases {
		s, e := helpers.ScrollRange(c.listLen, c.cursor, c.maxVisible)
		if s != c.wantStart || e != c.wantEnd {
			t.Fatalf("ScrollRange(%d,%d,%d) = %d,%d want %d,%d", c.listLen, c.cursor, c.maxVisible, s, e, c.wantStart, c.wantEnd)
		}
	}
}

func TestDeleteWordBackwards(t *testing.T) {
	cases := []struct {
		s          string
		cursor     int
		wantS      string
		wantCursor int
	}{
		{"hello world", 11, "hello ", 6},
		{"hello world", 5, " world", 0},
		{"one two three", 8, "one three", 4},
		{"leading", 0, "leading", 0},
		{"a b c", 3, " c", 0},
	}

	for _, c := range cases {
		gotS, gotCursor := helpers.DeleteWordBackwards(c.s, c.cursor)
		if gotS != c.wantS || gotCursor != c.wantCursor {
			t.Fatalf("DeleteWordBackwards(%q,%d) = %q,%d want %q,%d", c.s, c.cursor, gotS, gotCursor, c.wantS, c.wantCursor)
		}
	}
}
