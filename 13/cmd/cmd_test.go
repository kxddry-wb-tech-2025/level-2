package cmd

import (
	"strings"
	"testing"

	"cut/internal/config"
)

func TestParseFields(t *testing.T) {
	tests := []struct {
		input   string
		want    map[int]struct{}
		wantErr bool
		errMsg  string
	}{
		{"", nil, false, ""},
		{"1", map[int]struct{}{1: {}}, false, ""},
		{"1,3-5", map[int]struct{}{1: {}, 3: {}, 4: {}, 5: {}}, false, ""},
		{"5-3", map[int]struct{}{3: {}, 4: {}, 5: {}}, false, ""},
		{"abc", nil, true, "gotta have at least some columns"},
		{"0", map[int]struct{}{0: {}}, false, ""},
	}

	for _, tt := range tests {
		got, err := parseFields(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("parseFields(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if err != nil && err.Error() != tt.errMsg {
			t.Errorf("parseFields(%q) error = %v, wantErrMsg %v", tt.input, err.Error(), tt.errMsg)
		}

		if len(got) != len(tt.want) {
			t.Errorf("parseFields(%q) got = %v, want %v", tt.input, got, tt.want)
			continue
		}

		for k := range tt.want {
			if _, ok := got[k]; !ok {
				t.Errorf("parseFields(%q) missing key %d in result", tt.input, k)
			}
		}
	}
}

func TestProcess(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		cfg      config.Options
		expected string
	}{
		{
			name:  "Show all fields",
			input: "a,b,c\n1,2,3\n",
			cfg: config.Options{
				ShowAll:   true,
				Delimiter: ',',
			},
			expected: "a\tb\tc\n1\t2\t3\n",
		},
		{
			name:  "Select fields 1 and 3",
			input: "a,b,c\n1,2,3\n",
			cfg: config.Options{
				ShowAll:   false,
				Delimiter: ',',
				Fields:    map[int]struct{}{1: {}, 3: {}},
			},
			expected: "a\tc\n1\t3\n",
		},
		{
			name:  "SepOnly true skips lines without delimiter",
			input: "a,b,c\nlinewithoutdelimiter\n1,2,3\n",
			cfg: config.Options{
				ShowAll:   true,
				Delimiter: ',',
				SepOnly:   true,
			},
			expected: "a\tb\tc\n1\t2\t3\n",
		},
		{
			name:  "Empty input",
			input: "",
			cfg: config.Options{
				ShowAll:   true,
				Delimiter: ',',
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			var sb strings.Builder
			err := Process(r, &sb, tt.cfg)
			if err != nil {
				t.Fatalf("Process() error = %v", err)
			}
			got := sb.String()
			if got != tt.expected {
				t.Errorf("Process() = %q, want %q", got, tt.expected)
			}
		})
	}
}
