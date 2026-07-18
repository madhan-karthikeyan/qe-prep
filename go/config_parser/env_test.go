package config_parser

import (
	"testing"
)

func TestParseEnvBasic(t *testing.T) {
	input := `KEY=value
EMPTY=
# comment
ANOTHER=val`
	got, err := ParseEnv(input)
	if err != nil {
		t.Fatalf("ParseEnv: %v", err)
	}
	if got["KEY"] != "value" {
		t.Errorf("KEY = %q, want 'value'", got["KEY"])
	}
	if got["EMPTY"] != "" {
		t.Errorf("EMPTY = %q, want ''", got["EMPTY"])
	}
	if _, ok := got["# comment"]; ok {
		t.Error("comment line produced a key")
	}
	if got["ANOTHER"] != "val" {
		t.Errorf("ANOTHER = %q, want 'val'", got["ANOTHER"])
	}
}

func TestParseEnvQuoted(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"double quoted", `KEY="hello world"`, "hello world"},
		{"single quoted", "KEY='hello world'", "hello world"},
		{"escape n", `KEY="line1\nline2"`, "line1\nline2"},
		{"escape t", `KEY="col1\tcol2"`, "col1\tcol2"},
		{"escape backslash", `KEY="path\\to"`, "path\\to"},
		{"escape quote", `KEY="say \"hi\""`, `say "hi"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseEnv(tt.input)
			if err != nil {
				t.Fatalf("ParseEnv: %v", err)
			}
			if got["KEY"] != tt.want {
				t.Errorf("KEY = %q, want %q", got["KEY"], tt.want)
			}
		})
	}
}

func TestParseEnvSubstitution(t *testing.T) {
	input := `BASE=/usr/local
PATH=$BASE/bin
LOGDIR=${BASE}/logs
ESCAPED=\$BASE`
	got, err := ParseEnv(input)
	if err != nil {
		t.Fatalf("ParseEnv: %v", err)
	}
	if got["PATH"] != "/usr/local/bin" {
		t.Errorf("PATH = %q, want %q", got["PATH"], "/usr/local/bin")
	}
	if got["LOGDIR"] != "/usr/local/logs" {
		t.Errorf("LOGDIR = %q, want %q", got["LOGDIR"], "/usr/local/logs")
	}
	if got["ESCAPED"] != "$BASE" {
		t.Errorf("ESCAPED = %q, want $BASE", got["ESCAPED"])
	}
}

func TestParseEnvComments(t *testing.T) {
	input := `# This is a comment
KEY=val
 ; another comment style`
	got, err := ParseEnv(input)
	if err != nil {
		t.Fatalf("ParseEnv: %v", err)
	}
	if got["KEY"] != "val" {
		t.Errorf("KEY = %q", got["KEY"])
	}
	if len(got) != 1 {
		t.Errorf("expected 1 key, got %d", len(got))
	}
}

func TestParseEnvEmptyLines(t *testing.T) {
	input := "KEY1=val1\n\n\nKEY2=val2"
	got, err := ParseEnv(input)
	if err != nil {
		t.Fatalf("ParseEnv: %v", err)
	}
	if got["KEY1"] != "val1" || got["KEY2"] != "val2" {
		t.Errorf("unexpected results: %v", got)
	}
}

func TestParseEnvErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"unterminated double", `KEY="unterminated`},
		{"unterminated single", "KEY='unterminated"},
		{"empty key", "=value"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseEnv(tt.input)
			if err == nil {
				t.Errorf("expected error for %q", tt.input)
			}
		})
	}
}

func TestParseEnvWhitespace(t *testing.T) {
	input := "  KEY  =  value  "
	got, err := ParseEnv(input)
	if err != nil {
		t.Fatalf("ParseEnv: %v", err)
	}
	if got["KEY"] != "value" {
		t.Errorf("KEY = %q, want 'value'", got["KEY"])
	}
}
