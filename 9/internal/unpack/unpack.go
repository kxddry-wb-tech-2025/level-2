package unpack

import (
	"errors"
	"strconv"
	"unicode"
)

var (
	// ErrEndsBackslash comes up when the input string ends with a backslash
	ErrEndsBackslash = errors.New("unpack: ends with a backslash")

	// ErrNoCharToMultiply comes up when the input string has digits with no characters prior
	ErrNoCharToMultiply = errors.New("unpack: no character to multiply")
)

// String unpacks a string
func String(s string) (string, error) {
	var res []rune
	runes := []rune(s)
	n := len(runes)
	if n == 0 {
		return "", nil
	}

	var last rune
	var hasChar bool
	i := 0

	for i < n {
		ch := runes[i]

		if ch == '\\' {
			i++
			if i >= n {
				return "", ErrEndsBackslash
			}
			last = runes[i]
			res = append(res, last)
			hasChar = true
			i++
			continue
		}

		if unicode.IsLetter(ch) || !unicode.IsDigit(ch) {
			last = ch
			res = append(res, last)
			hasChar = true
			i++
			continue
		}

		if unicode.IsDigit(ch) {
			if !hasChar {
				return "", ErrNoCharToMultiply
			}
			j := i
			for j < n && unicode.IsDigit(runes[j]) {
				j++
			}
			numStr := string(runes[i:j])
			count, err := strconv.Atoi(numStr)
			if err != nil {
				return "", err
			}
			if count < 1 {
				i = j
				continue
			}
			for range count - 1 {
				res = append(res, last)
			}
			i = j
			continue
		}

		res = append(res, last)
		hasChar = true
		i++
	}

	if !hasChar {
		return "", nil
	}

	return string(res), nil
}
