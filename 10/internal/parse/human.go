package parse

import (
	"math"
	"strings"
	"unicode"
)

var pow1024 = []float64{
	1,
	1024,
	math.Pow(1024, 2),
	math.Pow(1024, 3),
	math.Pow(1024, 4),
	math.Pow(1024, 5),
	math.Pow(1024, 6),
}

// HumanNumber parses "1K", "1.5M", "2G", "10T", optionally with trailing "B" or "iB".
// Uses powers of 1024 like GNU sort -h.
// Returns (value, true) if parsed, else (0, false).
func HumanNumber(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}

	// separate numeric prefix and suffix
	i := 0
	for i < len(s) && (unicode.IsDigit(rune(s[i])) || s[i] == '.' || s[i] == '+' || s[i] == '-') {
		i++
	}
	numStr := s[:i]
	suf := strings.TrimSpace(s[i:])

	num, ok := FloatLoose(numStr)
	if !ok {
		return 0, false
	}
	if suf == "" {
		return num, true
	}
	// Normalize suffix
	suf = strings.ToUpper(suf)
	// allow trailing B or iB
	suf = strings.TrimSuffix(suf, "B")
	suf = strings.TrimSuffix(suf, "IB")

	mul := 1.0
	switch suf {
	case "":
		mul = 1
	case "K":
		mul = pow1024[1]
	case "M":
		mul = pow1024[2]
	case "G":
		mul = pow1024[3]
	case "T":
		mul = pow1024[4]
	case "P":
		mul = pow1024[5]
	case "E":
		mul = pow1024[6]
	default:
		// Unknown suffix -> not human numeric
		return 0, false
	}

	return num * mul, true
}
