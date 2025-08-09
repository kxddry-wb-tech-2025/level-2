package unpack_test

import (
	"testing"

	"l2.9/internal/unpack"
)

func TestUnpackString(t *testing.T) {
	tests := []struct {
		Input   string
		Want    string
		WantErr error
	}{
		{"a4bc2d5e", "aaaabccddddde", nil},
		{"abcd", "abcd", nil},
		{"45", "", unpack.ErrNoCharToMultiply},
		{"", "", nil},

		{"qwe\\4\\5", "qwe45", nil},
		{"qwe\\45", "qwe44444", nil},
	}

	for _, test := range tests {
		got, gotErr := unpack.String(test.Input)
		if gotErr != test.WantErr {
			t.Errorf("String(%q): got error %v, want %v", test.Input, gotErr, test.WantErr)
		}
		if got != test.Want {
			t.Errorf("String(%q): got %q, want %q", test.Input, got, test.Want)
		}
	}
}
