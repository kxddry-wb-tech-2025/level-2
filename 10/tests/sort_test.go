package tests

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

const prog = "../gosort"

func runCLI(t *testing.T, args []string, input string) (stdout string, stderr string, err error) {
	t.Helper()
	cmd := exec.Command(prog, args...)
	cmd.Stdin = strings.NewReader(input)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	return strings.TrimRight(outBuf.String(), "\n"), strings.TrimRight(errBuf.String(), "\n"), err
}

func TestCLI(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		input       string
		expectedOut string
		expectErr   bool
	}{
		{
			name:        "Базовая сортировка",
			args:        []string{},
			input:       "c\na\nb",
			expectedOut: "a\nb\nc",
		},
		{
			name:        "Сортировка по первому столбцу (-k 1)",
			args:        []string{"-k", "1"},
			input:       "3\tz\n1\ta\n2\tb",
			expectedOut: "1\ta\n2\tb\n3\tz",
		},
		{
			name:        "Сортировка по второму столбцу (-k 2)",
			args:        []string{"-k", "2"},
			input:       "3\tz\n1\ta\n2\tb",
			expectedOut: "1\ta\n2\tb\n3\tz",
		},
		{
			name:        "Числовая сортировка (-n)",
			args:        []string{"-n"},
			input:       "10\n2\n33\n1",
			expectedOut: "1\n2\n10\n33",
		},
		{
			name:        "Числовая сортировка в обратном порядке (-nr)",
			args:        []string{"-nr"},
			input:       "10\n2\n33\n1",
			expectedOut: "33\n10\n2\n1",
		},
		{
			name:        "Удаление дубликатов (-u)",
			args:        []string{"-u"},
			input:       "a\na\nb\nb\nc",
			expectedOut: "a\nb\nc",
		},
		{
			name:        "Сортировка по месяцам (-M)",
			args:        []string{"-M"},
			input:       "Feb\nJan\nDec\nApr",
			expectedOut: "Jan\nFeb\nApr\nDec",
		},
		{
			name:        "Игнорирование хвостовых пробелов (-b)",
			args:        []string{"-b"},
			input:       "a   \n   b\nc",
			expectedOut: "   b\na   \nc",
		},
		{
			name:        "Человекочитаемые размеры (-h)",
			args:        []string{"-h"},
			input:       "2K\n1M\n512\n3K",
			expectedOut: "512\n2K\n3K\n1M",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, errOut, err := runCLI(t, tc.args, tc.input)

			if tc.expectErr {
				if err == nil {
					t.Errorf("\nТест: %s\nОжидалась ошибка, но программа завершилась без неё", tc.name)
				}
				return
			}

			if err != nil {
				t.Errorf("\nТест: %s\nНеожиданная ошибка: %v\nstderr: %s", tc.name, err, errOut)
			}

			if out != tc.expectedOut {
				t.Errorf("\nТест: %s\nОжидалось:\n%q\nПолучилось:\n%q", tc.name, tc.expectedOut, out)
			}
		})
	}
}

func TestCheckSorted(t *testing.T) {
	// Отсортированные
	out, errOut, err := runCLI(t, []string{"-c"}, "a\nb\nc")
	if err != nil {
		t.Errorf("\n-c: Ожидалось, что ошибки не будет, stderr=%q", errOut)
	}
	if out != "" {
		t.Errorf("\n-c: Ожидался пустой вывод, получили %q", out)
	}

	// Неотсортированные
	_, errOut, err = runCLI(t, []string{"-c"}, "b\na\nc")
	if err == nil {
		t.Errorf("\n-c: Ожидалась ошибка при неотсортированных данных, stderr=%q", errOut)
	}
}
