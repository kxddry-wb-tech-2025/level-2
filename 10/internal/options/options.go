package options

// Options serves as a struct for config flags
type Options struct {
	KeyColumn            int  // 1-based; 0 means whole line
	Numeric              bool // -n
	Reverse              bool // -r
	Unique               bool // -u
	Month                bool // -M
	IgnoreTrailingBlanks bool // -b (trailing)
	Check                bool // -c
	HumanNumeric         bool // -h
}
