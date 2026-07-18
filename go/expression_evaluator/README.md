# Expression Evaluator
Difficulty: Medium
Estimated Interview Time: 35 min
Prerequisites: stacks, operator precedence, parsing

## Problem Statement
Implement an arithmetic expression evaluator using the Shunting Yard algorithm
and RPN evaluation.

## Requirements
- Support +, -, *, /, ^, %, parentheses
- Unary minus detection
- Correct precedence and associativity
- Float64 evaluation

## Implementation Notes
- Shunting Yard infix-to-postfix conversion
- Stack-based RPN evaluator
- Unary minus emitted as "_" token

## Test Strategy
- Table-driven tests for each operator
- Precedence and associativity tests
- Error cases (division by zero, invalid tokens)

## Edge Cases
- Unary minus at expression start
- Unary minus after operator or parenthesis
- Floating point numbers
- Nested parentheses

## Failure Cases
- Division by zero
- Modulo by zero
- Mismatched parentheses
- Unknown characters

## Complexity (Time + Space)
- Conversion: O(n) time, O(n) space
- Evaluation: O(n) time, O(n) space

## Progression Path
- Add functions (sin, cos, sqrt)
- Add variables
- Add string operations

## Common Interview Follow-ups
- Operator precedence tables
- Right vs left associativity
- Error recovery

## Possible Production Improvements
- Use math/big for arbitrary precision
- Add expression AST for optimization
- Support conditional expressions
