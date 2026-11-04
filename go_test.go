package main

import (
	"testing"
	"net/url"
)

func TestHandleRelLinks(t *testing.T) {
	base, _ := url.Parse("https://github.com/golang/go/tree/master/src")

	tests := []struct {
		input    string
		expected string
	}{
		// Absolute URLs pass through unchanged
		{"https://example.com", "https://example.com"},
		{"http://example.com/page", "http://example.com/page"},

		// Root-relative: resolve against the site root, not the current page
		{"/about", "https://github.com/about"},
		{"/", "https://github.com/"},

		// Scheme-relative: a different host, same scheme as the base
		{"//docs.github.com/en", "https://docs.github.com/en"},

		// Directory-relative: replace everything after the last slash
		{"README.md", "https://github.com/golang/go/tree/master/README.md"},
		{"./README.md", "https://github.com/golang/go/tree/master/README.md"},
		{"example.com", "https://github.com/golang/go/tree/master/example.com"},

		// Parent directories
		{"../doc", "https://github.com/golang/go/tree/doc"},
		{"../../..", "https://github.com/golang/"},

		// Query-only: keep the current path, swap the query
		{"?tab=readme", "https://github.com/golang/go/tree/master/src?tab=readme"},

		// Fragments stripped, so these are the same page as without them
		{"#readme", "https://github.com/golang/go/tree/master/src"},
		{"/about#team", "https://github.com/about"},
		{"", "https://github.com/golang/go/tree/master/src"},

		// Not crawlable: return "" to mean "skip"
		{"mailto:hi@example.com", ""},
		{"tel:+15555555555", ""},
		{"javascript:void(0)", ""},
	}

	for _, test := range tests {
		result := handleRelLinks(base, test.input)
		if result != test.expected {
			t.Errorf("handleRelLinks(%q, %q) = %q; want %q", base, test.input, result, test.expected)
		}
	}
}
