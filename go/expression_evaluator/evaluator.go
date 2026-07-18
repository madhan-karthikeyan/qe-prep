package expression_evaluator

import (
	"fmt"
	"math"
	"strconv"
)

// Evaluate evaluates an infix arithmetic expression and returns the float64
// result. Supported operators: +, -, *, /, ^, %, parentheses, and unary minus.
func Evaluate(expr string) (float64, error) {
	postfix, err := InfixToPostfix(expr)
	if err != nil {
		return 0, fmt.Errorf("infix to postfix: %w", err)
	}
	return EvaluateRPN(postfix)
}

// EvaluateRPN evaluates a Reverse Polish Notation expression.
func EvaluateRPN(tokens []string) (float64, error) {
	var stack []float64

	for _, tok := range tokens {
		if isNumber(tok) {
			val, err := strconv.ParseFloat(tok, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid number %q", tok)
			}
			stack = append(stack, val)
			continue
		}

		if tok == "_" {
			if len(stack) < 1 {
				return 0, fmt.Errorf("not enough operands for unary minus")
			}
			stack[len(stack)-1] = -stack[len(stack)-1]
			continue
		}

		if !IsOperator(tok) {
			return 0, fmt.Errorf("unknown token %q", tok)
		}

		if len(stack) < 2 {
			return 0, fmt.Errorf("not enough operands for %q", tok)
		}
		b := stack[len(stack)-1]
		a := stack[len(stack)-2]
		stack = stack[:len(stack)-2]

		var result float64
		switch tok {
		case "+":
			result = a + b
		case "-":
			result = a - b
		case "*":
			result = a * b
		case "/":
			if b == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			result = a / b
		case "%":
			if b == 0 {
				return 0, fmt.Errorf("modulo by zero")
			}
			result = math.Mod(a, b)
		case "^":
			result = math.Pow(a, b)
		default:
			return 0, fmt.Errorf("unknown operator %q", tok)
		}
		stack = append(stack, result)
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression: %d values remain on stack", len(stack))
	}
	return stack[0], nil
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
