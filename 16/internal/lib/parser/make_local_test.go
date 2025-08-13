package parser

import (
	"bytes"
	"io"
	"net/url"
	"strings"
	"testing"
)

func TestMakeLocalRewrites(t *testing.T) {
	html := `<html>
    <a href="/p1"></a>
    <img src="img.png">
    <source srcset="a.png 1x, b.png 2x">
    <div style="background:url('bg.png')"></div>
    </html>`
	base, _ := url.Parse("https://example.com")
	p := New()
	out := p.MakeLocal(io.NopCloser(strings.NewReader(html)), base, "text/html", "out")
	buf := new(bytes.Buffer)
	_, _ = io.Copy(buf, out)
	s := buf.String()
	if !strings.Contains(s, "out/p1.html") {
		t.Errorf("expected p1.html rewrite, got %s", s)
	}
}
