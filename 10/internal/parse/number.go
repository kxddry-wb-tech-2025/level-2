package parse

import (
	"strconv"
	"strings"
)

// FloatLoose tries to parse a float from the string, trimming spaces.
// Returns (value, true) if parsed, otherwise (0, false).
func FloatLoose(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	// Allow underscores like 1_000? Keep it strict.
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}
