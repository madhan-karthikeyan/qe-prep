# Expression Evaluator
Difficulty: Hard
Estimated Interview Time: 45 min
Prerequisites: Stack, parsing, Shunting Yard algorithm, RPN

## Problem Statement
Implement a mathematical expression evaluator using the Shunting Yard algorithm to convert infix to RPN, then evaluate the RPN expression.

## Requirements
- Double-based arithmetic: +, -, *, /, ^, %
- Correct operator precedence and associativity
- Parentheses for grouping
- Unary minus support
- Tokenization handles decimals and negative numbers

## Implementation Notes
- ShuntingYard converts infix to RPN token list
- Evaluator evaluates RPN using a double stack
- Unary minus marked as "_" token in RPN
- Right-associative exponentiation (^) with correct precedence

## Test Strategy
- Precedence and associativity verification
- Nested parentheses
- Unary minus (leading, consecutive)
- Decimal and integer numbers
- Division by zero

## Edge Cases
- Division by zero → ArithmeticException
- Modulo by zero → ArithmeticException
- Mismatched parentheses → IllegalArgumentException
- Consecutive operators (e.g., "1+-2")

## Failure Cases
- null/blank expression → IllegalArgumentException
- Invalid character → IllegalArgumentException
- Division by zero → ArithmeticException

## Complexity (Time + Space)
- ShuntingYard: O(n) time, O(n) space
- Evaluator: O(n) time, O(n) space

## Progression Path
Implement tokenization first, then Shunting Yard, then evaluator. Add operators incrementally.

## Common Interview Follow-ups
- Add functions (sin, cos, log)
- Variable substitution
- Expression validation before evaluation

## Possible Production Improvements
- Use BigDecimal for precise arithmetic
- Compile expressions to bytecode for performance
- Support for custom functions and constants
