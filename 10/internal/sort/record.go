package sort

// Record contains a line with some specs
type Record struct {
	Line     string
	Key      string
	Num      float64
	HasNum   bool
	Month    int
	HasMonth bool
	SortMode Mode
}

// Mode is taken as a separate type to avoid errors
type Mode int

const (
	// ModeString means "consider lines strings"
	ModeString Mode = iota
	// ModeNumeric means "consider lines numbers"
	ModeNumeric
	// ModeHuman means "consider lines with human-readable output"
	ModeHuman
	// ModeMonth means "parse months"
	ModeMonth
)
