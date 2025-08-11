package shell

import (
	"bytes"
	"os"
	"reflect"
	"testing"
)

func TestExpandVariables(t *testing.T) {
	const envKey = "SHELL_TEST_FOO"
	_ = os.Setenv(envKey, "bar")
	defer func() { _ = os.Unsetenv(envKey) }()

	s := New()
	in := "$" + envKey + " baz $MISSING"
	got := s.expandVariables(in)
	want := "bar baz "
	if got != want {
		t.Fatalf("expandVariables(%q) = %q; want %q", in, got, want)
	}
}

func TestParseCommand_VariableExpansion(t *testing.T) {
	const envKey = "SHELL_TEST_FOO2"
	_ = os.Setenv(envKey, "xyz")
	defer func() { _ = os.Unsetenv(envKey) }()

	s := New()
	got := s.parseCommand("$" + envKey + " one two")
	want := []string{"xyz", "one", "two"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseCommand(...) = %#v; want %#v", got, want)
	}
}

func TestExecOutside_RunExternal(t *testing.T) {
	s := New()
	var out bytes.Buffer

	ok := s.execOutside(nil, &out, "sh", "-c", "printf hi")
	if !ok {
		t.Fatalf("execOutside returned false")
	}
	if out.String() != "hi" {
		t.Fatalf("execOutside wrote %q; want %q", out.String(), "hi")
	}
}

func TestExecCond_AndOrBehavior(t *testing.T) {
	s := New()

	if !s.execCond("pwd && pwd", "&&") {
		t.Fatalf("expected 'pwd && pwd' to succeed")
	}

	if s.execCond("unknowncommand && pwd", "&&") {
		t.Fatalf("expected 'unknowncommand && pwd' to fail")
	}

	if !s.execCond("unknowncommand || pwd", "||") {
		t.Fatalf("expected 'unknowncommand || pwd' to succeed")
	}

	if !s.execCond("pwd || unknowncommand", "||") {
		t.Fatalf("expected 'pwd || unknowncommand' to succeed")
	}
}
