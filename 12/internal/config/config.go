package config

// Flags represent a config structure for program flags
type Flags struct {
	After        int
	Before       int
	OnlyCount    bool
	IgnoreCase   bool
	Invert       bool
	FixedString  bool
	PrintNumbers bool
}
