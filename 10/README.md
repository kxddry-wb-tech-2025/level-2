# l2.10: sort utility


> A simplified analogue of the UNIX sort utility.
> The program reads lines from STDIN or a file and returns sorted strings.

---

## Usage

```bash
./gosort -flags [file.txt]
```

---

## Flags

- `-k N` - sorts according to the N-th column, default delimiter is tab.

- `-n` - sorts by integer value
- `-r` - sorts in reverse
- `-u` - only shows unique lines

Additionally,

- `-M` - sorts based on months in lines (Jan, May, Mar, etc.)
- `-b` - ignores trailing blanks
- `-c` - only checks whether the lines are sorted
- `-h` - sorts based on human-readable suffixes (K for Kilobytes, etc.)

---
## Quickstart

### 1. Compile the program

```bash
git clone https://github.com/kxddry-wb-tech-2025/level-2
cd level-2/10
go build .
```

### 2. Run tests

```bash
go test gosort/tests
```
