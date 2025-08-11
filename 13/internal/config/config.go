package config

// Options is a structure that stores app flags.
type Options struct {
	Fields    map[int]struct{}
	ShowAll   bool
	Delimiter rune
	SepOnly   bool
}
