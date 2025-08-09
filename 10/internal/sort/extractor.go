package sort

import (
	"strings"
	"unicode"

	"gosort/internal/options"
)

// Extractor extracts a key from a line and tells you whether the keys are equal in two lines.
type Extractor interface {
	Key(line string) string
	EqualKey(a, b string) bool
}

type extractor struct {
	keyColumn int
	trimRight bool
}

// NewExtractor creates a new Extractor.
func NewExtractor(opt options.Options) Extractor {
	return &extractor{
		keyColumn: opt.KeyColumn,
		trimRight: opt.IgnoreTrailingBlanks,
	}
}

func (e *extractor) Key(line string) string {
	var key string
	if e.keyColumn <= 0 {
		key = line
	} else {
		key = getTabColumn(line, e.keyColumn)
	}
	if e.trimRight {
		key = strings.TrimRightFunc(key, unicode.IsSpace)
	}
	return key
}

func (e *extractor) EqualKey(a, b string) bool {
	if e.trimRight {
		a = strings.TrimRightFunc(a, unicode.IsSpace)
		b = strings.TrimRightFunc(b, unicode.IsSpace)
	}
	return a == b
}

// getTabColumn returns the 1-based column split by tabs.
// If not present, returns empty string.
func getTabColumn(s string, col int) string {
	if col <= 1 {
		// first column is chars before first '\t' or whole line
		i := strings.IndexByte(s, '\t')
		if i == -1 {
			return s
		}
		return s[:i]
	}
	i := 0
	start := 0
	for idx := 0; idx < len(s); idx++ {
		if s[idx] == '\t' {
			i++
			if i == col-1 {
				start = idx + 1
			} else if i == col {
				return s[start:idx]
			}
		}
	}
	if i == col-1 {
		return s[start:]
	}
	return ""
}
