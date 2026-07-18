package config_parser

import (
	"testing"
)

func TestParseINIBasic(t *testing.T) {
	input := `[server]
host=localhost
port=8080

[database]
name=testdb
user=admin`
	got, err := ParseINI(input)
	if err != nil {
		t.Fatalf("ParseINI: %v", err)
	}
	if got["server"]["host"] != "localhost" {
		t.Errorf("server.host = %q", got["server"]["host"])
	}
	if got["server"]["port"] != "8080" {
		t.Errorf("server.port = %q", got["server"]["port"])
	}
	if got["database"]["name"] != "testdb" {
		t.Errorf("database.name = %q", got["database"]["name"])
	}
	if got["database"]["user"] != "admin" {
		t.Errorf("database.user = %q", got["database"]["user"])
	}
}

func TestParseINIComments(t *testing.T) {
	input := `# top comment
; semicolon comment
[section]
key=value
# inline comment`
	got, err := ParseINI(input)
	if err != nil {
		t.Fatalf("ParseINI: %v", err)
	}
	if got["section"]["key"] != "value" {
		t.Errorf("key = %q", got["section"]["key"])
	}
}

func TestParseINIBlankLines(t *testing.T) {
	input := "[section]\n\n\nkey=val"
	got, err := ParseINI(input)
	if err != nil {
		t.Fatalf("ParseINI: %v", err)
	}
	if got["section"]["key"] != "val" {
		t.Errorf("key = %q", got["section"]["key"])
	}
}

func TestParseINIMultipleSections(t *testing.T) {
	input := `[section1]
a=1

[section2]
b=2
c=3`
	got, err := ParseINI(input)
	if err != nil {
		t.Fatalf("ParseINI: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 sections, got %d", len(got))
	}
	if got["section1"]["a"] != "1" {
		t.Errorf("section1.a = %q", got["section1"]["a"])
	}
	if got["section2"]["b"] != "2" || got["section2"]["c"] != "3" {
		t.Errorf("section2 = %v", got["section2"])
	}
}

func TestParseINIErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"no section", "key=val"},
		{"malformed section", "[section"},
		{"empty key", "[s]\n=val"},
		{"no equals", "[s]\nkey"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseINI(tt.input)
			if err == nil {
				t.Errorf("expected error for %q", tt.input)
			}
		})
	}
}

func TestParseINIValueWithEquals(t *testing.T) {
	input := `[section]
connection=postgres://user:pass@host/db?sslmode=verify-full`
	got, err := ParseINI(input)
	if err != nil {
		t.Fatalf("ParseINI: %v", err)
	}
	want := "postgres://user:pass@host/db?sslmode=verify-full"
	if got["section"]["connection"] != want {
		t.Errorf("connection = %q, want %q", got["section"]["connection"], want)
	}
}
