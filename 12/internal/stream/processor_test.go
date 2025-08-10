package stream

import (
	"grep/internal/config"
	"regexp"
	"strings"
	"testing"
)

func TestNewProcessor(t *testing.T) {
	flags := &config.Flags{
		After:        2,
		Before:       1,
		OnlyCount:    true,
		IgnoreCase:   false,
		Invert:       false,
		FixedString:  true,
		PrintNumbers: true,
	}

	processor := NewProcessor(flags)

	if processor == nil {
		t.Fatal("NewProcessor returned nil")
	}

	if processor.opt != flags {
		t.Error("NewProcessor did not set flags correctly")
	}
}

func TestProcessorIsMatch(t *testing.T) {
	tests := []struct {
		name     string
		flags    config.Flags
		line     string
		pattern  string
		regex    *regexp.Regexp
		expected bool
	}{
		{
			name:     "fixed string match case sensitive",
			flags:    config.Flags{FixedString: true, IgnoreCase: false},
			line:     "Hello World",
			pattern:  "Hello",
			regex:    nil,
			expected: true,
		},
		{
			name:     "fixed string no match case sensitive",
			flags:    config.Flags{FixedString: true, IgnoreCase: false},
			line:     "Hello World",
			pattern:  "hello",
			regex:    nil,
			expected: false,
		},
		{
			name:     "fixed string match case insensitive",
			flags:    config.Flags{FixedString: true, IgnoreCase: true},
			line:     "Hello World",
			pattern:  "hello",
			regex:    nil,
			expected: true,
		},
		{
			name:     "fixed string no match case insensitive",
			flags:    config.Flags{FixedString: true, IgnoreCase: true},
			line:     "Hello World",
			pattern:  "xyz",
			regex:    nil,
			expected: false,
		},
		{
			name:     "regex match",
			flags:    config.Flags{FixedString: false},
			line:     "test123",
			pattern:  "",
			regex:    regexp.MustCompile(`\d+`),
			expected: true,
		},
		{
			name:     "regex no match",
			flags:    config.Flags{FixedString: false},
			line:     "testABC",
			pattern:  "",
			regex:    regexp.MustCompile(`\d+`),
			expected: false,
		},
		{
			name:     "regex nil returns false",
			flags:    config.Flags{FixedString: false},
			line:     "test123",
			pattern:  "",
			regex:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewProcessor(&tt.flags)
			result := processor.isMatch(tt.line, tt.pattern, tt.regex)
			if result != tt.expected {
				t.Errorf("isMatch() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestProcessStreamBasicMatching(t *testing.T) {
	tests := []struct {
		name          string
		flags         config.Flags
		input         string
		pattern       string
		regex         *regexp.Regexp
		expectedCount int
	}{
		{
			name:          "simple fixed string match",
			flags:         config.Flags{FixedString: true},
			input:         "line1\ntest line\nline3\ntest again",
			pattern:       "test",
			regex:         nil,
			expectedCount: 2,
		},
		{
			name:          "case insensitive match",
			flags:         config.Flags{FixedString: true, IgnoreCase: true},
			input:         "Test\nTEST\ntest\nother",
			pattern:       "test",
			regex:         nil,
			expectedCount: 3,
		},
		{
			name:          "regex match",
			flags:         config.Flags{FixedString: false},
			input:         "abc123\ndef456\nghi\njkl789",
			pattern:       "",
			regex:         regexp.MustCompile(`\d+`),
			expectedCount: 3,
		},
		{
			name:          "inverted match",
			flags:         config.Flags{FixedString: true, Invert: true},
			input:         "match\nnomatch\nmatch\nnomatch",
			pattern:       "match",
			regex:         nil,
			expectedCount: 0,
		},
		{
			name:          "no matches",
			flags:         config.Flags{FixedString: true},
			input:         "line1\nline2\nline3",
			pattern:       "notfound",
			regex:         nil,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewProcessor(&tt.flags)
			reader := strings.NewReader(tt.input)

			count, err := processor.ProcessStream(reader, tt.pattern, tt.regex)

			if err != nil {
				t.Errorf("ProcessStream() returned error: %v", err)
			}

			if count != tt.expectedCount {
				t.Errorf("ProcessStream() count = %d, expected %d", count, tt.expectedCount)
			}
		})
	}
}

func TestProcessStreamContextLines(t *testing.T) {
	input := "line1\nline2\nmatch\nline4\nline5\nmatch2\nline7"

	tests := []struct {
		name    string
		flags   config.Flags
		pattern string
	}{
		{
			name: "before context",
			flags: config.Flags{
				FixedString: true,
				Before:      2,
			},
			pattern: "match",
		},
		{
			name: "after context",
			flags: config.Flags{
				FixedString: true,
				After:       2,
			},
			pattern: "match",
		},
		{
			name: "before and after context",
			flags: config.Flags{
				FixedString: true,
				Before:      1,
				After:       1,
			},
			pattern: "match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewProcessor(&tt.flags)
			reader := strings.NewReader(input)

			count, err := processor.ProcessStream(reader, tt.pattern, nil)

			if err != nil {
				t.Errorf("ProcessStream() returned error: %v", err)
			}

			// Should find 2 matches (match and match2)
			if count != 2 {
				t.Errorf("ProcessStream() count = %d, expected 2", count)
			}
		})
	}
}

func TestProcessStreamOnlyCount(t *testing.T) {
	input := "test\nline\ntest\nother"
	flags := config.Flags{
		FixedString: true,
		OnlyCount:   true,
	}

	processor := NewProcessor(&flags)
	reader := strings.NewReader(input)

	count, err := processor.ProcessStream(reader, "test", nil)

	if err != nil {
		t.Errorf("ProcessStream() returned error: %v", err)
	}

	if count != 2 {
		t.Errorf("ProcessStream() count = %d, expected 2", count)
	}
}

func TestProcessStreamPrintNumbers(t *testing.T) {
	input := "line1\nmatch\nline3"
	flags := config.Flags{
		FixedString:  true,
		PrintNumbers: true,
	}

	processor := NewProcessor(&flags)
	reader := strings.NewReader(input)

	count, err := processor.ProcessStream(reader, "match", nil)

	if err != nil {
		t.Errorf("ProcessStream() returned error: %v", err)
	}

	if count != 1 {
		t.Errorf("ProcessStream() count = %d, expected 1", count)
	}
}

func TestProcessStreamEmptyInput(t *testing.T) {
	flags := config.Flags{FixedString: true}
	processor := NewProcessor(&flags)
	reader := strings.NewReader("")

	count, err := processor.ProcessStream(reader, "test", nil)

	if err != nil {
		t.Errorf("ProcessStream() returned error: %v", err)
	}

	if count != 0 {
		t.Errorf("ProcessStream() count = %d, expected 0", count)
	}
}

func TestProcessStreamSingleLine(t *testing.T) {
	flags := config.Flags{FixedString: true}
	processor := NewProcessor(&flags)
	reader := strings.NewReader("single test line")

	count, err := processor.ProcessStream(reader, "test", nil)

	if err != nil {
		t.Errorf("ProcessStream() returned error: %v", err)
	}

	if count != 1 {
		t.Errorf("ProcessStream() count = %d, expected 1", count)
	}
}

func TestProcessStreamLargeInput(t *testing.T) {
	// Test with input larger than buffer size
	var builder strings.Builder
	for i := 0; i < 1000; i++ {
		if i%100 == 0 {
			builder.WriteString("match\n")
		} else {
			builder.WriteString("line\n")
		}
	}

	flags := config.Flags{FixedString: true}
	processor := NewProcessor(&flags)
	reader := strings.NewReader(builder.String())

	count, err := processor.ProcessStream(reader, "match", nil)

	if err != nil {
		t.Errorf("ProcessStream() returned error: %v", err)
	}

	if count != 10 {
		t.Errorf("ProcessStream() count = %d, expected 10", count)
	}
}

func TestProcessStreamComplexRegex(t *testing.T) {
	input := "user@example.com\ninvalid-email\ntest@domain.org\nnot-an-email"
	flags := config.Flags{FixedString: false}
	processor := NewProcessor(&flags)
	reader := strings.NewReader(input)

	// Email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	count, err := processor.ProcessStream(reader, "", emailRegex)

	if err != nil {
		t.Errorf("ProcessStream() returned error: %v", err)
	}

	if count != 2 {
		t.Errorf("ProcessStream() count = %d, expected 2", count)
	}
}

func TestProcessStreamContextOverlap(t *testing.T) {
	// Test case where context lines from different matches might overlap
	input := "line1\nmatch1\nline3\nmatch2\nline5"
	flags := config.Flags{
		FixedString: true,
		Before:      1,
		After:       1,
	}

	processor := NewProcessor(&flags)
	reader := strings.NewReader(input)

	count, err := processor.ProcessStream(reader, "match", nil)

	if err != nil {
		t.Errorf("ProcessStream() returned error: %v", err)
	}

	if count != 2 {
		t.Errorf("ProcessStream() count = %d, expected 2", count)
	}
}

// Benchmark tests
func BenchmarkProcessStreamFixedString(b *testing.B) {
	input := strings.Repeat("line without match\n", 1000) +
		strings.Repeat("line with test match\n", 100)
	flags := config.Flags{FixedString: true}
	processor := NewProcessor(&flags)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(input)
		_, _ = processor.ProcessStream(reader, "test", nil)
	}
}

func BenchmarkProcessStreamRegex(b *testing.B) {
	input := strings.Repeat("line123\n", 1000) +
		strings.Repeat("lineABC\n", 100)
	flags := config.Flags{FixedString: false}
	processor := NewProcessor(&flags)
	regex := regexp.MustCompile(`\d+`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(input)
		_, _ = processor.ProcessStream(reader, "", regex)
	}
}
