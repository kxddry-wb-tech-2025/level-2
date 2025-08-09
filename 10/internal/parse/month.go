package parse

import "strings"

// Map of month abbreviations and full names (case-insensitive) to 1..12
var months = map[string]int{
	"JAN": 1,
	"FEB": 2,
	"MAR": 3,
	"APR": 4,
	"MAY": 5,
	"JUN": 6,
	"JUL": 7,
	"AUG": 8,
	"SEP": 9,
	"OCT": 10,
	"NOV": 11,
	"DEC": 12,
}

// Month tries to extract a month number from s.
// Strategy: split on whitespace and punctuation, find the first token matching a known month.
// Returns (month, true) or (0, false).
func Month(s string) (int, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}

	// A quick pass: try start of string 3 or more chars
	upper := strings.ToUpper(s)
	for _, tok := range tokenize(upper) {
		if m, ok := months[tok]; ok {
			return m, true
		}
	}
	return 0, false
}

func tokenize(s string) []string {
	var out []string
	cur := make([]rune, 0, len(s))
	flush := func() {
		if len(cur) > 0 {
			out = append(out, string(cur))
			cur = cur[:0]
		}
	}
	for _, r := range s {
		if r == ' ' || r == '\t' || r == '\n' || r == ',' || r == '.' || r == '-' || r == '/' || r == ':' {
			flush()
		} else {
			cur = append(cur, r)
		}
	}
	flush()
	return out
}
