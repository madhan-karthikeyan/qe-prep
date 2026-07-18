package expression_evaluator

import (
	"math"
	"testing"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  float64
	}{
		{"addition", "1+2", 3},
		{"subtraction", "5-3", 2},
		{"multiplication", "4*3", 12},
		{"division", "10/2", 5},
		{"modulo", "7%3", 1},
		{"exponentiation", "2^3", 8},
		{"precedence", "1+2*3", 7},
		{"parentheses", "(1+2)*3", 9},
		{"nested parens", "((2+3)*4)", 20},
		{"unary minus", "-5", -5},
		{"double unary", "--5", 5},
		{"unary with op", "-3+4", 1},
		{"unary in parens", "-(3+4)", -7},
		{"decimals", "3.5+2.5", 6},
		{"exponent right assoc", "2^3^2", 512},
		{"complex", "2*(3+4)^2", 98},
		{"modulo precedence", "10+7%3", 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.input)
			if err != nil {
				t.Fatalf("Evaluate(%q): %v", tt.input, err)
			}
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Evaluate(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestEvaluateErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"division by zero", "5/0"},
		{"modulo by zero", "5%0"},
		{"empty", ""},
		{"invalid expr", "1+"},
		{"unknown token", "1 @ 2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Evaluate(tt.input)
			if err == nil {
				t.Errorf("expected error for input %q", tt.input)
			}
		})
	}
}

func TestEvaluateExtraTokens(t *testing.T) {
	tokens := []string{"1", "2", "+", "3"}
	_, err := EvaluateRPN(tokens)
	if err == nil {
		t.Error("expected error for extra operands")
	}
}
