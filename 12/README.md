# grep: a CLI utility for searching texts, written in Go


## Usage:

1. Compile the application
```bash
git clone https://github.com/kxddry-wb-tech-2025/level-2
cd 12
go build ./cmd/grep
```

2. Run the application
```bash
grep [FILENAME (optional)] [flags] 
```

Usage:

```bash
Usage:
  grep [FILENAME (optional)] [flags]

Flags:
  -A, --after-context int    show N lines after each found expression
  -B, --before-context int   show N lines before each found expression
  -C, --context int          show N lines before and after each found expression
  -c, --count                show only matching count
  -F, --fixed-string         fix string instead of regexp
  -h, --help                 help for grep
  -i, --ignore-case          ignore case matching
  -v, --invert               invert matching
  -n, --print-numbers        print line numbers

```