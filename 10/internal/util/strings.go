package util

import "unicode"

// TrimRightSpace removes trailing Unicode space characters.
func TrimRightSpace(s string) string {
	i := len(s)
	for i > 0 && unicode.IsSpace(rune(s[i-1])) {
		i--
	}
	if i == len(s) {
		return s
	}
	return s[:i]
}
