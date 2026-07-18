package file_processing

import (
	"os"
	"strings"
	"testing"
)

func TestWordCount(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		data string
		want Counts
	}{
		{
			name: "empty",
			data: "",
			want: Counts{},
		},
		{
			name: "one line",
			data: "hello world",
			want: Counts{Lines: 0, Words: 2, Characters: 11, Bytes: 11},
		},
		{
			name: "one line with newline",
			data: "hello world\n",
			want: Counts{Lines: 1, Words: 2, Characters: 12, Bytes: 12},
		},
		{
			name: "multiple lines",
			data: "hello world\nfoo bar baz\n",
			want: Counts{Lines: 2, Words: 5, Characters: 24, Bytes: 24},
		},
		{
			name: "unicode",
			data: "héllo wörld\n",
			want: Counts{Lines: 1, Words: 2, Characters: 12, Bytes: 14},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := WordCount(strings.NewReader(tc.data))
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Errorf("WordCount = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestProcessFiles(t *testing.T) {
	dir := t.TempDir()
	aPath := dir + "/a.txt"
	bPath := dir + "/b.txt"
	if err := os.WriteFile(aPath, []byte("hello world\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bPath, []byte("foo bar baz\n"), 0644); err != nil {
		t.Fatal(err)
	}
	results, err := ProcessFiles([]string{aPath, bPath})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	c := results[aPath]
	if c.Lines != 1 || c.Words != 2 {
		t.Errorf("a.txt: got %+v, want Lines=1 Words=2", c)
	}
	c = results[bPath]
	if c.Lines != 1 || c.Words != 3 {
		t.Errorf("b.txt: got %+v, want Lines=1 Words=3", c)
	}
}

func TestWordCountUnicode(t *testing.T) {
	input := "你好 世界\n"
	c, err := WordCount(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if c.Characters != 6 {
		t.Errorf("expected 6 chars (2 hanzi + space + 2 hanzi + newline), got %d", c.Characters)
	}
}
