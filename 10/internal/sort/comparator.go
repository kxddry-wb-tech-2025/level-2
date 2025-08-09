package sort

import (
	"gosort/internal/options"
	"gosort/internal/parse"
)

// Comparator is an interface that can compare records
type Comparator interface {
	Compare(a, b Record) int
	Mode() Mode
}

type comparator struct {
	mode      Mode
	extractor Extractor
}

// NewComparator creates a comparator
func NewComparator(opt options.Options) Comparator {
	var mode Mode
	switch {
	case opt.HumanNumeric:
		mode = ModeHuman
	case opt.Numeric:
		mode = ModeNumeric
	case opt.Month:
		mode = ModeMonth
	default:
		mode = ModeString
	}
	return &comparator{
		mode: mode,
	}
}

// Mode returns mode.
func (c *comparator) Mode() Mode { return c.mode }

// MakeRecord creates a record and returns it.
func MakeRecord(line string, extractor Extractor, cmp Comparator) Record {
	key := extractor.Key(line)
	rec := Record{
		Line:     line,
		Key:      key,
		SortMode: cmp.Mode(),
	}
	switch cmp.Mode() {
	case ModeNumeric:
		if v, ok := parse.FloatLoose(key); ok {
			rec.Num = v
			rec.HasNum = true
		}
	case ModeHuman:
		if v, ok := parse.HumanNumber(key); ok {
			rec.Num = v
			rec.HasNum = true
		}
	case ModeMonth:
		if m, ok := parse.Month(key); ok {
			rec.Month = m
			rec.HasMonth = true
		}
	default:
	}
	return rec
}

// BuildRecords creates records.
func BuildRecords(lines []string, extractor Extractor, cmp Comparator) []Record {
	out := make([]Record, len(lines))
	for i, ln := range lines {
		out[i] = MakeRecord(ln, extractor, cmp)
	}
	return out
}

// Compare compares records.
func (c *comparator) Compare(a, b Record) int {
	switch c.mode {
	case ModeNumeric, ModeHuman:
		if a.HasNum && b.HasNum {
			if a.Num < b.Num {
				return -1
			} else if a.Num > b.Num {
				return 1
			}
			// tie-break by key then full line for stability
			if a.Key < b.Key {
				return -1
			} else if a.Key > b.Key {
				return 1
			}
			if a.Line < b.Line {
				return -1
			} else if a.Line > b.Line {
				return 1
			}
			return 0
		}
		// If one has number and other doesn't, put non-parsed after parsed
		if a.HasNum && !b.HasNum {
			return -1
		}
		if !a.HasNum && b.HasNum {
			return 1
		}
		// both no numbers -> fallback to lexicographic key, then full line
		return compareStrings(a.Key, b.Key, a.Line, b.Line)
	case ModeMonth:
		if a.HasMonth && b.HasMonth {
			if a.Month < b.Month {
				return -1
			} else if a.Month > b.Month {
				return 1
			}
			// tie-breakers
			return compareStrings(a.Key, b.Key, a.Line, b.Line)
		}
		if a.HasMonth && !b.HasMonth {
			return -1
		}
		if !a.HasMonth && b.HasMonth {
			return 1
		}
		return compareStrings(a.Key, b.Key, a.Line, b.Line)
	default:
		return compareStrings(a.Key, b.Key, a.Line, b.Line)
	}
}

func compareStrings(ka, kb, la, lb string) int {
	if ka < kb {
		return -1
	} else if ka > kb {
		return 1
	}
	if la < lb {
		return -1
	} else if la > lb {
		return 1
	}
	return 0
}

// Unique returns unique records
func Unique(records []Record, cmp Comparator) []Record {
	if len(records) == 0 {
		return records
	}
	out := records[:0]
	out = append(out, records[0])
	for i := 1; i < len(records); i++ {
		prev := out[len(out)-1]
		cur := records[i]
		// unique by key equality for the chosen mode
		eq := false
		switch cmp.Mode() {
		case ModeNumeric, ModeHuman:
			eq = prev.HasNum && cur.HasNum && prev.Num == cur.Num && prev.Key == cur.Key
		case ModeMonth:
			eq = prev.HasMonth && cur.HasMonth && prev.Month == cur.Month && prev.Key == cur.Key
		default:
			eq = prev.Key == cur.Key
		}
		if !eq {
			out = append(out, cur)
		}
	}
	return out
}
