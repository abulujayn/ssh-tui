package parser

import (
	"strings"
	"testing"
)

func TestParseHostPort(t *testing.T) {
	cases := []struct {
		in       string
		wantHost string
		wantPort string
	}{
		{"example.com", "example.com", ""},
		{"example.com:2222", "example.com", "2222"},
		{"[2001:db8::1]:2222", "2001:db8::1", "2222"},
		{"[host.example.com]:2200", "host.example.com", "2200"},
		{"[badbracket:2200", "[badbracket", "2200"},
	}

	for _, c := range cases {
		h, p := parseHostPort(c.in)
		if h != c.wantHost || p != c.wantPort {
			t.Fatalf("parseHostPort(%q) = (%q,%q), want (%q,%q)", c.in, h, p, c.wantHost, c.wantPort)
		}
	}
}

func TestIsValidHost(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"192.168.1.1", true},
		{"user@192.168.1.1", true},
		{"2001:db8::1", true},
		{"user@2001:db8::1", true},
		{"example.com", true},
		{"User@Example.COM", true},
		{"localhost", false},
		{"", false},
		{"user@", false},
		{"invalid_domain", false},
	}

	for _, c := range cases {
		got := IsValidHost(c.in)
		if got != c.want {
			t.Fatalf("IsValidHost(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParseUserHost(t *testing.T) {
	cases := []struct {
		in       string
		wantUser string
		wantHost string
	}{
		{"user@host", "user", "host"},
		{"host", "", "host"},
		{" user@host ", "user", "host"},
		{"user@", "user", ""},
	}

	for _, c := range cases {
		u, h := ParseUserHost(c.in)
		if u != c.wantUser || h != c.wantHost {
			t.Fatalf("ParseUserHost(%q) = (%q,%q), want (%q,%q)", c.in, u, h, c.wantUser, c.wantHost)
		}
	}
}

func TestFilterHosts_OrderAndMatching(t *testing.T) {
	hosts := []SSHHost{
		{Name: "alpha", HostName: "alpha.com", Aliases: []string{"a1"}},
		{Name: "beta", HostName: "beta.com", Aliases: []string{"b1", "search"}},
		{Name: "searchhost", HostName: "searchhost.com", Aliases: nil},
		{Name: "other", HostName: "other.com", Aliases: []string{"ssearch"}},
	}

	filtered := FilterHosts(hosts, "search")

	if len(filtered) < 2 {
		t.Fatalf("expected at least 2 results, got %d", len(filtered))
	}

	// exact alias match should come first (beta has alias "search")
	if filtered[0].Name != "beta" {
		t.Fatalf("expected first result to be 'beta' (exact alias match), got %q", filtered[0].Name)
	}

	// primary prefix match should follow (searchhost)
	if filtered[1].Name != "searchhost" {
		t.Fatalf("expected second result to be 'searchhost' (primary prefix), got %q", filtered[1].Name)
	}
}

func TestFormatHostDisplay_ContainsDetails(t *testing.T) {
	h := SSHHost{
		Name:     "myserver",
		HostName: "example.com",
		User:     "me",
		Port:     "2222",
	}

	out := FormatHostDisplay(h)

	if !strings.HasPrefix(out, "myserver") {
		t.Fatalf("expected output to start with host name, got %q", out)
	}

	if !strings.Contains(out, "host: example.com") {
		t.Fatalf("expected host detail to include hostname, got %q", out)
	}

	if !strings.Contains(out, "user: me") {
		t.Fatalf("expected host detail to include user, got %q", out)
	}

	if !strings.Contains(out, "port: 2222") {
		t.Fatalf("expected host detail to include port, got %q", out)
	}
}
