package file_processing

import (
	"errors"
	"io"
	"strings"
	"unicode/utf8"
)

var (
	ErrUnexpectedQuote = errors.New("unexpected quote")
	ErrUnterminated    = errors.New("unterminated quoted field")
)

// csvParser reads CSV data from an io.Reader with a configurable delimiter.
type csvParser struct {
	reader    io.Reader
	delimiter rune
	buf       []byte
	offset    int
}

// newCSVParser creates a new csvParser.
func newCSVParser(r io.Reader, delimiter rune) *csvParser {
	if delimiter == 0 {
		delimiter = ','
	}
	return &csvParser{
		reader:    r,
		delimiter: delimiter,
	}
}

// readRune reads the next rune from the buffer, refilling from the reader if
// needed.
func (p *csvParser) readRune() (rune, int, error) {
	if p.offset >= len(p.buf) {
		tmp := make([]byte, 4096)
		n, err := p.reader.Read(tmp)
		if n > 0 {
			p.buf = tmp[:n]
			p.offset = 0
		}
		if err != nil {
			return 0, 0, err
		}
	}
	r, size := utf8.DecodeRune(p.buf[p.offset:])
	p.offset += size
	return r, size, nil
}

// unreadRune unreads the last rune read.
func (p *csvParser) unreadRune() {
	if p.offset > 0 {
		r, size := utf8.DecodeLastRune(p.buf[:p.offset])
		p.offset -= size
		_ = r
	}
}

// parseField reads one CSV field, handling quotes.
func (p *csvParser) parseField() (string, error) {
	r, _, err := p.readRune()
	if err != nil {
		return "", err
	}
	if r == '"' {
		return p.parseQuotedField()
	}
	var b strings.Builder
	if r == p.delimiter || r == '\n' || r == '\r' {
		p.unreadRune()
		return b.String(), nil
	}
	b.WriteRune(r)
	for {
		r, _, err := p.readRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return "", err
		}
		if r == p.delimiter || r == '\n' || r == '\r' {
			p.unreadRune()
			break
		}
		b.WriteRune(r)
	}
	return b.String(), nil
}

// parseQuotedField reads a quoted CSV field, handling escaped quotes.
func (p *csvParser) parseQuotedField() (string, error) {
	var b strings.Builder
	for {
		r, _, err := p.readRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return "", ErrUnterminated
			}
			return "", err
		}
		if r == '"' {
			next, _, err := p.readRune()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return "", err
			}
			if next == '"' {
				b.WriteRune('"')
				continue
			}
			if next == p.delimiter || next == '\n' || next == '\r' {
				p.unreadRune()
				break
			}
			return "", ErrUnexpectedQuote
		}
		b.WriteRune(r)
	}
	return b.String(), nil
}

// ParseCSV reads all records from r. The first record may be a header.
func ParseCSV(r io.Reader, delimiter rune) ([][]string, error) {
	p := newCSVParser(r, delimiter)
	var records [][]string
	for {
		record, err := p.parseRecord()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

// parseRecord reads one CSV record (a slice of fields).
func (p *csvParser) parseRecord() ([]string, error) {
	var fields []string
	lastSep := p.delimiter
	for {
		field, err := p.parseField()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if len(fields) == 0 {
					return nil, err
				}
				if lastSep == p.delimiter {
					fields = append(fields, "")
				}
				return fields, nil
			}
			return nil, err
		}
		fields = append(fields, field)
		r, _, err := p.readRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if lastSep == p.delimiter && field == "" {
					fields = append(fields, "")
				}
				return fields, nil
			}
			return nil, err
		}
		if r == '\n' {
			return fields, nil
		}
		if r == '\r' {
			r2, _, err := p.readRune()
			if err == nil && r2 != '\n' {
				p.unreadRune()
			}
			return fields, nil
		}
		lastSep = r
	}
}

// ParseCSVCallback reads CSV records and calls callback for each record.
func ParseCSVCallback(r io.Reader, delimiter rune, callback func([]string) error) error {
	p := newCSVParser(r, delimiter)
	for {
		record, err := p.parseRecord()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		if err := callback(record); err != nil {
			return err
		}
	}
}

// DetectHeader returns true if the first record looks like a header (all fields
// contain no numbers and match a typical header pattern).
func DetectHeader(records [][]string) bool {
	if len(records) == 0 {
		return false
	}
	for _, field := range records[0] {
		if field == "" {
			return false
		}
		hasLetter := false
		for _, r := range field {
			if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && r != '_' {
				return false
			}
			hasLetter = true
		}
		if !hasLetter {
			return false
		}
	}
	return true
}
