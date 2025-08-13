package parser

import (
	"bytes"
	"io"
	"net/url"
	"path/filepath"
	"strings"

	"wget/internal/models"

	"golang.org/x/net/html"
)

// Parser implements HTML parsing.
type Parser struct{}

// New returns a new Parser.
func New() *Parser {
	return &Parser{}
}

// HTML parses the HTML and returns unique page links and resources.
// Faster by: early tag filtering, on-the-fly dedup, single parse/resolve per attr,
// iterative traversal (no recursion), and a lightweight srcset parser.
func (p *Parser) HTML(htm io.Reader, base *url.URL) ([]*url.URL, []models.Resource, error) {
	doc, err := html.Parse(htm)
	if err != nil {
		return nil, nil, err
	}

	links := make([]*url.URL, 0, 32)
	resources := make([]models.Resource, 0, 64)
	seenLinks := make(map[string]struct{}, 128)
	seenRes := make(map[string]struct{}, 256)

	// Helpers
	resolve := func(raw string) string {
		// keep your existing logic here if resolveURL has extra rules;
		// otherwise you can inline: parsed, _ := url.Parse(raw); return base.ResolveReference(parsed).String()
		if strings.Index(raw, "?") != -1 {
			return resolveURL(raw[:strings.Index(raw, "?")], base)
		}

		return resolveURL(raw, base)
	}

	addLink := func(abs string) {
		if abs == "" {
			return
		}
		if _, ok := seenLinks[abs]; ok {
			return
		}
		u, err := url.Parse(abs)
		if err != nil {
			return
		}
		seenLinks[abs] = struct{}{}
		links = append(links, u)
	}

	addRes := func(abs string) {
		if abs == "" {
			return
		}
		if _, ok := seenRes[abs]; ok {
			return
		}
		u, err := url.Parse(abs)
		if err != nil {
			return
		}
		seenRes[abs] = struct{}{}
		resources = append(resources, models.Resource{
			URL:  u,
			Type: getResourceType(abs),
		})
	}

	getAttr := func(n *html.Node, key string) (string, bool) {
		// In HTML mode, keys are already lowercase.
		for i := range n.Attr {
			a := n.Attr[i]
			if a.Namespace == "" && a.Key == key {
				return a.Val, true
			}
		}
		return "", false
	}

	parseSrcsetFast := func(s string) []string {
		// srcset := comma-separated candidates; for each candidate, URL is up to first whitespace.
		out := make([]string, 0, 4)
		start := 0
		for i := 0; i <= len(s); i++ {
			if i == len(s) || s[i] == ',' {
				part := strings.TrimSpace(s[start:i])
				if part != "" {
					// take up to first whitespace
					j := -1
					for k := 0; k < len(part); k++ {
						switch part[k] {
						case ' ', '\n', '\t', '\r', '\f':
							j = k
							k = len(part)
						}
					}
					if j >= 0 {
						part = part[:j]
					}
					if part != "" {
						out = append(out, part)
					}
				}
				start = i + 1
			}
		}
		return out
	}

	// Iterative pre-order traversal (avoids deep recursion/function-call overhead)
	for n := doc; n != nil; {
		if n.Type == html.ElementNode {
			switch n.Data { // tag names are lowercased by the parser
			// Clickable document links
			case "a", "area":
				if v, ok := getAttr(n, "href"); ok {
					addLink(resolve(v))
				}

			// External resources by href
			case "link":
				if v, ok := getAttr(n, "href"); ok {
					addRes(resolve(v))
				}

			// Sources by src/srcset
			case "img", "script", "iframe", "source", "audio", "video", "track", "embed", "input":
				if v, ok := getAttr(n, "src"); ok {
					addRes(resolve(v))
				}
				if v, ok := getAttr(n, "srcset"); ok {
					for _, u := range parseSrcsetFast(v) {
						addRes(resolve(u))
					}
				}
			}

			// Inline CSS: only scan if it even contains url(
			if v, ok := getAttr(n, "style"); ok && strings.Contains(v, "url(") {
				for _, cssURL := range extractCSSURLs(v) {
					addRes(resolve(cssURL))
				}
			}
		}

		// Move to next node (pre-order)
		if n.FirstChild != nil {
			n = n.FirstChild
			continue
		}
		for n != nil && n.NextSibling == nil {
			n = n.Parent
		}
		if n != nil {
			n = n.NextSibling
		}
	}

	return links, resources, nil
}

// MakeLocal rewrites resource URLs in HTML to local file paths.
func (p *Parser) MakeLocal(htm io.ReadCloser, base *url.URL, contentType, outputDir string) io.ReadCloser {
	data, _ := io.ReadAll(htm)
	_ = htm.Close()
	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return htm
	}
	guessCT := func(tag, attr, path, fallback string) string {
		ext := strings.ToLower(filepath.Ext(path))
		switch {
		case tag == "script" || ext == ".js":
			return "application/javascript"
		case tag == "link" && attr == "href" && (ext == ".css" || ext == ""):
			return "text/css"
		case ext == ".css":
			return "text/css"
		case ext == ".html" || ext == ".htm" || ext == "":
			return "text/html"
		default:
			return fallback
		}
	}

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for i, a := range n.Attr {
				key := strings.ToLower(a.Key)

				switch key {
				case "href", "src":
					abs := resolveURL(a.Val, base)
					if abs == "" {
						continue
					}

					parsed, err := url.Parse(abs)
					if err != nil {
						continue
					}

					if parsed.Host != "" && parsed.Host != base.Host {
						continue
					}

					ct := guessCT(n.Data, key, parsed.Path, contentType)
					n.Attr[i].Val = GenerateLocalPath(parsed, ct, outputDir)
				case "srcset":
					var newParts []string
					for _, part := range splitSrcset(a.Val) {
						fields := strings.Fields(part)
						if len(fields) > 0 {
							if abs := resolveURL(fields[0], base); abs != "" {
								if parsed, err := url.Parse(abs); err == nil {
									ct := "text/html"
									ext := strings.ToLower(filepath.Ext(parsed.Path))
									if ext == ".css" {
										ct = "text/css"
									} else if ext == ".js" {
										ct = "application/javascript"
									}
									fields[0] = GenerateLocalPath(parsed, ct, outputDir)
								}
							}
							newParts = append(newParts, strings.Join(fields, " "))
						}
					}
					n.Attr[i].Val = strings.Join(newParts, ", ")
				case "style":
					// always rewrite css urls
					n.Attr[i].Val = rewriteCSSURLs(a.Val, base, contentType, outputDir)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return io.NopCloser(bytes.NewReader(data))
	}

	return io.NopCloser(bytes.NewReader(buf.Bytes()))
}
