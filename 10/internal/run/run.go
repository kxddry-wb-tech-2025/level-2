package run

import (
	"bufio"
	"fmt"
	"gosort/internal/util"
	"io"
	"os"
	"sort"

	"gosort/internal/options"
	"gosort/internal/reader"
	isort "gosort/internal/sort"
)

// Run runs the program
func Run(opt options.Options, args []string) error {
	in, err := reader.OpenInput(args)
	if err != nil {
		return err
	}
	defer func() {
		if in != os.Stdin {
			_ = in.Close()
		}
	}()

	extractor := isort.NewExtractor(opt)
	comparator := isort.NewComparator(opt)

	if opt.Check {
		return checkSorted(in, extractor, comparator, opt)
	}

	lines, err := reader.ReadAllLines(in)
	if err != nil {
		return err
	}

	if opt.IgnoreTrailingBlanks {
		for i := range lines {
			lines[i] = util.TrimRightSpace(lines[i])
		}
	}

	records := isort.BuildRecords(lines, extractor, comparator)

	sort.Slice(records, func(i, j int) bool {
		cmp := comparator.Compare(records[i], records[j])
		if opt.Reverse {
			return cmp > 0
		}
		return cmp < 0
	})

	if opt.Unique {
		records = isort.Unique(records, comparator)
	}

	// Write to stdout
	w := bufio.NewWriterSize(os.Stdout, 64*1024)
	for _, rec := range records {
		_, _ = w.WriteString(rec.Line)
		_, _ = w.WriteString("\n")
	}
	_ = w.Flush()
	return nil
}

func checkSorted(in io.Reader, extractor isort.Extractor, comparator isort.Comparator, opt options.Options) error {
	rd := bufio.NewReaderSize(in, 64*1024)

	readLine := func() (string, error) {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			if len(line) == 0 {
				return "", io.EOF
			}
			// strip eol
			if n := len(line); n > 0 {
				if line[n-1] == '\n' {
					line = line[:n-1]
					n--
				}
				if n > 0 && line[n-1] == '\r' {
					line = line[:n-1]
				}
			}
			return line, nil
		}
		if err != nil {
			return "", err
		}
		// strip eol
		if n := len(line); n > 0 {
			if line[n-1] == '\n' {
				line = line[:n-1]
				n--
			}
			if n > 0 && line[n-1] == '\r' {
				line = line[:n-1]
			}
		}
		return line, nil
	}

	var (
		prevLine string
		lineNo   = 0
	)

	// read first line
	first, err := readLine()
	if err == io.EOF {
		// empty input is sorted
		return nil
	}
	if err != nil {
		return err
	}
	lineNo = 1
	prevLine = first
	prevRec := isort.MakeRecord(prevLine, extractor, comparator)

	for {
		curLine, err := readLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		lineNo++
		curRec := isort.MakeRecord(curLine, extractor, comparator)
		cmp := comparator.Compare(prevRec, curRec)
		if opt.Reverse {
			cmp = -cmp
		}
		if cmp > 0 {
			// disorder found at current line
			_, _ = fmt.Fprintf(os.Stdout, "Data is not sorted at line %d\n", lineNo)
			return fmt.Errorf("not sorted at line %d", lineNo)
		}
		prevRec = curRec
	}

	// sorted
	return nil
}
