package file_processing

import (
	"io"
	"strings"
	"testing"
)

func TestParseCSV(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		delim   rune
		want    [][]string
		wantErr bool
	}{
		{
			name:  "simple",
			input: "a,b,c\n1,2,3",
			delim: ',',
			want:  [][]string{{"a", "b", "c"}, {"1", "2", "3"}},
		},
		{
			name:  "quoted with comma",
			input: `"hello, world",foo`,
			delim: ',',
			want:  [][]string{{"hello, world", "foo"}},
		},
		{
			name:  "escaped quote",
			input: `"say ""hello""",bar`,
			delim: ',',
			want:  [][]string{{`say "hello"`, "bar"}},
		},
		{
			name:  "tab delimited",
			input: "a\tb\tc\n1\t2\t3",
			delim: '\t',
			want:  [][]string{{"a", "b", "c"}, {"1", "2", "3"}},
		},
		{
			name:  "empty fields",
			input: ",,,\n1,,3,",
			delim: ',',
			want:  [][]string{{"", "", "", ""}, {"1", "", "3", ""}},
		},
		{
			name:  "trailing newline",
			input: "a,b\n",
			delim: ',',
			want:  [][]string{{"a", "b"}},
		},
		{
			name:  "carriage return newline",
			input: "a,b\r\n1,2",
			delim: ',',
			want:  [][]string{{"a", "b"}, {"1", "2"}},
		},
		{
			name:    "unterminated quote",
			input:   `"unterminated`,
			delim:   ',',
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseCSV(strings.NewReader(tc.input), tc.delim)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("got %d records, want %d", len(got), len(tc.want))
			}
			for i := range got {
				if len(got[i]) != len(tc.want[i]) {
					t.Fatalf("record %d: got %d fields, want %d", i, len(got[i]), len(tc.want[i]))
				}
				for j := range got[i] {
					if got[i][j] != tc.want[i][j] {
						t.Errorf("record %d field %d: got %q, want %q", i, j, got[i][j], tc.want[i][j])
					}
				}
			}
		})
	}
}

func TestParseCSVCallback(t *testing.T) {
	t.Parallel()
	var records [][]string
	err := ParseCSVCallback(strings.NewReader("a,b\n1,2"), ',', func(record []string) error {
		records = append(records, append([]string{}, record...))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 {
		t.Fatalf("got %d records, want 2", len(records))
	}
}

func TestDetectHeader(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		recs [][]string
		want bool
	}{
		{"header", [][]string{{"Name", "Age", "City"}}, true},
		{"data row", [][]string{{"Alice", "30", "NYC"}}, false},
		{"empty", [][]string{}, false},
		{"mixed", [][]string{{"Name", "30"}}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := DetectHeader(tc.recs); got != tc.want {
				t.Errorf("DetectHeader = %v, want %v", got, tc.want)
			}
		})
	}
}

func FuzzParseCSV(f *testing.F) {
	f.Add("a,b,c\n1,2,3")
	f.Add(`"hello","world"`)
	f.Add("")
	f.Add(",")
	f.Add("a,\"b\"\n")
	f.Fuzz(func(t *testing.T, input string) {
		records, err := ParseCSV(strings.NewReader(input), ',')
		if err != nil && err != io.EOF {
			// valid parse errors are acceptable; just don't panic
			return
		}
		_ = records
	})
}
