package parser

import (
	"net/url"
	"path/filepath"
	"strings"
	"wget/internal/models"
)

func resolveURL(href string, base *url.URL) string {
	if strings.HasPrefix(href, "data:") ||
		strings.HasPrefix(href, "mailto:") ||
		strings.HasPrefix(href, "javascript:") ||
		strings.HasPrefix(href, "#") {
		return ""
	}

	resolved, err := base.Parse(href)
	if err != nil {
		return ""
	}

	return resolved.String()
}

func getResourceType(u string) models.ResourceType {
	ext := strings.ToLower(filepath.Ext(u))
	switch ext {
	case ".css":
		return models.CSS
	case ".js":
		return models.JS
	case ".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp":
		return models.Image
	case ".html", ".htm":
		return models.HTML
	default:
		return models.Other
	}
}

// GenerateLocalPath generates a local path for a file
func GenerateLocalPath(u *url.URL, contentType, outputDir string) string {
	path := u.Path
	if path == "" || path == "/" {
		path = "/index.html"
	}

	if strings.HasSuffix(path, "/") {
		path += "index.html"
	}

	if filepath.Ext(path) == "" {
		switch contentType {
		case "text/html":
			path += ".html"
		case "text/css":
			path += ".css"
		case "application/javascript", "text/javascript":
			path += ".js"
		}
	}

	if path[0] == '/' {
		path = path[1:]
	}
	return filepath.Join(outputDir, filepath.Clean(path))
}

func splitSrcset(s string) []string {
	var out []string
	inParens := 0
	start := 0
	for i, ch := range s {
		switch ch {
		case '(':
			inParens++
		case ')':
			if inParens > 0 {
				inParens--
			}
		case ',':
			if inParens == 0 {
				out = append(out, strings.TrimSpace(s[start:i]))
				start = i + 1
			}
		}
	}
	if start < len(s) {
		out = append(out, strings.TrimSpace(s[start:]))
	}
	return out
}

func extractCSSURLs(s string) []string {
	var urls []string
	lower := strings.ToLower(s)
	idx := strings.Index(lower, "url(")
	for idx != -1 {
		// find the closing ')'
		endRel := strings.Index(lower[idx:], ")")
		if endRel == -1 {
			break
		}
		start := idx + 4 // len("url(")
		end := idx + endRel
		raw := strings.TrimSpace(s[start:end])
		raw = strings.Trim(raw, `"'`)
		if raw != "" {
			urls = append(urls, raw)
		}
		// search for next occurrence starting after the ')'
		next := strings.Index(lower[end+1:], "url(")
		if next == -1 {
			break
		}
		idx = end + 1 + next
	}
	return urls
}

func rewriteCSSURLs(s string, base *url.URL, contentType, outputDir string) string {
	urls := extractCSSURLs(s)
	for _, u := range urls {
		if abs := resolveURL(u, base); abs != "" {
			if parsed, err := url.Parse(abs); err == nil {
				local := GenerateLocalPath(parsed, contentType, outputDir)
				s = strings.ReplaceAll(s, u, local)
			}
		}
	}
	return s
}
