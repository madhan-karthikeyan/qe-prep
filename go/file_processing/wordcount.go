package file_processing

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

// Counts holds the results of counting a file.
type Counts struct {
	Lines      int
	Words      int
	Characters int
	Bytes      int64
}

// WordCount reads from r and returns line, word, character, and byte counts.
func WordCount(r io.Reader) (Counts, error) {
	var c Counts
	buf := make([]byte, 32*1024)
	inWord := false
	for {
		n, err := r.Read(buf)
		if n > 0 {
			c.Bytes += int64(n)
			chunk := buf[:n]
			for len(chunk) > 0 {
				r, size := utf8.DecodeRune(chunk)
				chunk = chunk[size:]
				c.Characters++
				if r == '\n' {
					c.Lines++
				}
				if unicode.IsSpace(r) {
					if inWord {
						c.Words++
						inWord = false
					}
				} else {
					inWord = true
				}
			}
		}
		if err != nil {
			if err == io.EOF {
				if inWord {
					c.Words++
				}
				return c, nil
			}
			return c, err
		}
	}
}

// WordCountReaderText is a convenience wrapper that uses bufio.Scanner for
// simple counting. This exists for cross-validation with WordCount.
func WordCountReaderText(r io.Reader) (Counts, error) {
	var c Counts
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		c.Lines++
		line := scanner.Text()
		c.Characters += utf8.RuneCountInString(line) + 1 // +1 for newline
		c.Bytes += int64(len(line)) + 1
		inWord := false
		for _, r := range line {
			if unicode.IsSpace(r) {
				if inWord {
					c.Words++
					inWord = false
				}
			} else {
				inWord = true
			}
		}
		if inWord {
			c.Words++
		}
	}
	return c, scanner.Err()
}

// fileCounts returns Counts for a single file path.
func fileCounts(path string) (Counts, error) {
	f, err := os.Open(path)
	if err != nil {
		return Counts{}, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	return WordCount(f)
}

// ProcessFiles returns a map from file path to Counts for each of the given
// files.
func ProcessFiles(paths []string) (map[string]Counts, error) {
	results := make(map[string]Counts, len(paths))
	for _, path := range paths {
		c, err := fileCounts(path)
		if err != nil {
			return nil, err
		}
		results[path] = c
	}
	return results, nil
}

// PrintCounts prints wc-style output for each file and a total.
func PrintCounts(w io.Writer, results map[string]Counts) {
	total := Counts{}
	first := true
	for path, c := range results {
		if !first {
			fmt.Fprintln(w)
		}
		first = false
		fmt.Fprintf(w, "%6d %6d %6d %s", c.Lines, c.Words, c.Characters, path)
		total.Lines += c.Lines
		total.Words += c.Words
		total.Characters += c.Characters
		total.Bytes += c.Bytes
	}
	if len(results) > 1 {
		fmt.Fprintf(w, "\n%6d %6d %6d total", total.Lines, total.Words, total.Characters)
	}
}

// RunWC processes files from args and prints results to stdout.
func RunWC(args []string) error {
	if len(args) == 0 {
		c, err := WordCount(os.Stdin)
		if err != nil {
			return err
		}
		PrintCounts(os.Stdout, map[string]Counts{"-": c})
		return nil
	}
	results, err := ProcessFiles(args)
	if err != nil {
		return err
	}
	PrintCounts(os.Stdout, results)
	return nil
}

// readerFromBytes wraps a byte slice as an io.Reader.
type readerFromBytes struct {
	data []byte
}

func (r *readerFromBytes) Read(p []byte) (int, error) {
	if len(r.data) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.data)
	r.data = r.data[n:]
	return n, nil
}

// NewReaderFromBytes creates an io.Reader from a byte slice.
func NewReaderFromBytes(data []byte) io.Reader {
	return bytes.NewReader(data)
}
