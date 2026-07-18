package expression_evaluator

import (
	"testing"
)

func TestInfixToPostfixBasic(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "simple addition",
			input: "1+2",
			want:  []string{"1", "2", "+"},
		},
		{
			name:  "multiplication precedence",
			input: "1+2*3",
			want:  []string{"1", "2", "3", "*", "+"},
		},
		{
			name:  "parentheses",
			input: "(1+2)*3",
			want:  []string{"1", "2", "+", "3", "*"},
		},
		{
			name:  "exponentiation right-assoc",
			input: "2^3^2",
			want:  []string{"2", "3", "2", "^", "^"},
		},
		{
			name:  "unary minus",
			input: "-5",
			want:  []string{"5", "_"},
		},
		{
			name:  "unary minus with parens",
			input: "-(3+4)",
			want:  []string{"3", "4", "+", "_"},
		},
		{
			name:  "all operators",
			input: "1+2-3*4/5^6%7",
			want:  []string{"1", "2", "+", "3", "4", "*", "5", "6", "^", "/", "7", "%", "-"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InfixToPostfix(tt.input)
			if err != nil {
				t.Fatalf("InfixToPostfix(%q): %v", tt.input, err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestInfixToPostfixErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"mismatched open paren", "((1+2)"},
		{"mismatched close paren", "(1+2))"},
		{"unknown char", "1 @ 2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := InfixToPostfix(tt.input)
			if err == nil {
				t.Errorf("expected error for input %q", tt.input)
			}
		})
	}
}
