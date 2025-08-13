# wget

### go utility mirroring ```wget -m``` functionality

---
## Quickstart

- Clone the project and compile the application

```bash
git clone https://github.com/kxddry-wb-tech-2025/level-2
cd level-2/16
go build .
./wget [flags] <url>
```

### Flags

- `-d N` - set max depth to N (default: 2)
- `-o <dir>` - output directory (default: hostname)
- `-w N` - number of workers (default: 5)
- `-t <10s>` - request timeout (default: 10 seconds)
- `-a <...>` - user-agent (default: "Wget/1.0")
- `-r` - ignore robots.txt file (default: false)
