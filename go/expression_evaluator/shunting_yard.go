package expression_evaluator

import (
	"fmt"
	"strconv"
	"unicode"
)

// Token types for internal representation.
const (
	tokNumber = iota
	tokOperator
	tokLParen
	tokRParen
)

type token struct {
	typ        int
	val        string
	prec       int
	rightAssoc bool
}

var opPrec = map[string]struct {
	prec  int
	right bool
}{
	"^": {4, true},
	"*": {3, false},
	"/": {3, false},
	"%": {3, false},
	"+": {2, false},
	"-": {2, false},
	"_": {5, true},
}

func tokenize(expr string) ([]token, error) {
	var tokens []token
	runes := []rune(expr)
	i := 0

	for i < len(runes) {
		c := runes[i]

		if unicode.IsSpace(c) {
			i++
			continue
		}

		if c == '(' {
			tokens = append(tokens, token{typ: tokLParen, val: "("})
			i++
			continue
		}

		if c == ')' {
			tokens = append(tokens, token{typ: tokRParen, val: ")"})
			i++
			continue
		}

		if c == '+' || c == '-' || c == '*' || c == '/' || c == '^' || c == '%' {
			tokens = append(tokens, token{typ: tokOperator, val: string(c)})
			i++
			continue
		}

		if unicode.IsDigit(c) || c == '.' {
			j := i
			hasDot := false
			for j < len(runes) && (unicode.IsDigit(runes[j]) || (runes[j] == '.' && !hasDot)) {
				if runes[j] == '.' {
					hasDot = true
				}
				j++
			}
			tokens = append(tokens, token{typ: tokNumber, val: string(runes[i:j])})
			i = j
			continue
		}

		return nil, fmt.Errorf("unexpected character %q at position %d", string(c), i)
	}

	return tokens, nil
}

// InfixToPostfix converts an infix expression token list to Reverse Polish
// Notation using the Shunting Yard algorithm. Unary minus is detected and
// emitted as "_".
func InfixToPostfix(expr string) ([]string, error) {
	tokens, err := tokenize(expr)
	if err != nil {
		return nil, err
	}

	var output []string
	var stack []token
	expectUnary := true

	for _, tok := range tokens {
		switch tok.typ {
		case tokNumber:
			output = append(output, tok.val)
			expectUnary = false

		case tokLParen:
			stack = append(stack, tok)
			expectUnary = true

		case tokRParen:
			for len(stack) > 0 && stack[len(stack)-1].val != "(" {
				output = append(output, stack[len(stack)-1].val)
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, fmt.Errorf("mismatched parentheses")
			}
			stack = stack[:len(stack)-1]
			expectUnary = false

		case tokOperator:
			op := tok.val
			if op == "-" && expectUnary {
				op = "_"
			}

			info := opPrec[op]
			if op != "_" {
				for len(stack) > 0 {
					top := stack[len(stack)-1]
					if top.typ == tokLParen {
						break
					}
					topInfo := opPrec[top.val]
					if top.val == "_" || info.prec < topInfo.prec ||
						(info.prec == topInfo.prec && !info.right) {
						output = append(output, top.val)
						stack = stack[:len(stack)-1]
					} else {
						break
					}
				}
			}
			stack = append(stack, token{typ: tokOperator, val: op, prec: info.prec, rightAssoc: info.right})
			expectUnary = true
		}
	}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		if top.val == "(" {
			return nil, fmt.Errorf("mismatched parentheses")
		}
		output = append(output, top.val)
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

// isValidNumber checks if s is a valid integer or float literal.
func isValidNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// IsOperator checks whether s is a known binary operator.
func IsOperator(s string) bool {
	_, ok := opPrec[s]
	return ok && s != "_"
}
