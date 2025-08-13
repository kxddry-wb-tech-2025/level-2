package parser

import (
	"net/url"
	"path/filepath"
	"testing"
)

func TestResolveURL(t *testing.T) {
	base, _ := url.Parse("https://example.com/dir/")
	tests := []struct {
		in   string
		want string
	}{
		{"page.html", "https://example.com/dir/page.html"},
		{"/root", "https://example.com/root"},
		{"//example.com/x", "https://example.com/x"},
		{"data:image/png;base64,xxx", ""},
		{"mailto:test@example.com", ""},
		{"javascript:alert(1)", ""},
		{"#frag", ""},
	}
	for _, tt := range tests {
		got := resolveURL(tt.in, base)
		if got != tt.want {
			t.Errorf("%q → got %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestGenerateLocalPath(t *testing.T) {
	base, _ := url.Parse("https://example.com/foo/")
	got := GenerateLocalPath(base, "text/html", "out")
	want := filepath.Join("out", "foo", "index.html")
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
