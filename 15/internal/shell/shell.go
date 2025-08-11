package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"shell/internal/shell/builtins"
	"strings"
)

// Shell represents the main structure responsible for the usage of go-shell
type Shell struct {
	reader *bufio.Reader
}

// New creates a new Shell
func New() *Shell {
	return &Shell{
		reader: bufio.NewReader(os.Stdin),
	}
}

// Run makes the shell run
func (s *Shell) Run() {
	c := make(chan os.Signal, 1)
	defer close(c)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			fmt.Println()
			fmt.Print("go-shell $ ")
		}
	}()

	fmt.Println("Go Shell -- type 'exit' or Ctrl+D to quit")

	for {
		fmt.Print("go-shell $ ")
		input, err := s.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nexit")
				break
			}
			fmt.Println("error reading input", err.Error())
			continue
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		if input == "exit" {
			fmt.Println("exit")
			break
		}
		s.exec(os.Stdin, os.Stdout, input)
	}

}

func (s *Shell) exec(in io.Reader, out io.Writer, input string) bool {

	if strings.Contains(input, "&&") {
		return s.execCond(input, "&&")
	}
	if strings.Contains(input, "||") {
		return s.execCond(input, "||")
	}

	if strings.Contains(input, ">") || strings.Contains(input, "<") {
		return s.execRedirect(input)
	}

	// Parse and execute single command
	args := s.parseCommand(input)
	if len(args) == 0 {
		return true
	}
	cmd := args[0]
	args1 := args[1:]

	switch cmd {
	case "cd":
		return builtins.Cd(in, out, args1...)
	case "exit":
		os.Exit(0)
	case "ps":
		return builtins.Ps(in, out, args1...)
	case "kill":
		return builtins.Kill(args1...)
	case "pwd":
		return builtins.Pwd(out)
	case "exec":
		return s.execOutside(in, out, args1...)
	case "echo":
		return builtins.Echo(in, out, args1...)
	default:
		return s.execOutside(in, out, args...)
	}
	return true
}

func (s *Shell) execRedirect(input string) bool {
	var inputFile, outputFile string

	if idx := strings.Index(input, ">"); idx != -1 {
		left := strings.TrimSpace(input[:idx])
		right := strings.TrimSpace(input[idx+1:])
		input = left
		outputFile = right
	}
	if idx := strings.Index(input, "<"); idx != -1 {
		left := strings.TrimSpace(input[:idx])
		right := strings.TrimSpace(input[idx+1:])
		input = left
		inputFile = right
	}

	args := s.parseCommand(strings.TrimSpace(input))
	if len(args) == 0 {
		return false
	}

	var stdin io.Reader = os.Stdin
	var stdout io.Writer = os.Stdout

	if inputFile != "" {
		f, err := os.Open(inputFile)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
			return false
		}
		defer func() { _ = f.Close() }()
		stdin = f
	}

	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			return false
		}
		defer func() { _ = f.Close() }()
		stdout = f
	}

	return s.exec(stdin, stdout, strings.Join(args, " "))
}

func (s *Shell) parseCommand(input string) []string {
	input = s.expandVariables(input)

	// Simple parsing - split by spaces (could be improved for quoted strings)
	return strings.Fields(input)
}

func (s *Shell) execPipelines(input string) {
	ss := strings.Split(input, "|")
	for i := range ss {
		ss[i] = strings.TrimSpace(ss[i])

	}
}

func (s *Shell) execOutside(in io.Reader, out io.Writer, args ...string) bool {
	if len(args) == 0 {
		return true
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "exec: %v\n", err)
		return false
	}
	return true
}

func (s *Shell) execCond(input string, cond string) bool {
	commands := strings.Split(input, cond)
	if len(commands) < 2 {
		fmt.Println("parse error: not enough arguments")
		return false
	}

	success := cond == "&&"
	for _, command := range commands {
		if (cond == "&&" && success) || (cond == "||" && !success) {
			success = s.exec(os.Stdin, os.Stdout, command)
		}
	}
	return success
}

func (s *Shell) expandVariables(input string) string {
	return os.ExpandEnv(input)
}
