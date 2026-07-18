# Expression Evaluator
Difficulty: Hard
Estimated Interview Time: 40 min
Prerequisites: stacks, parsing, RPN

## Problem Statement
Implement an expression evaluator using the shunting-yard algorithm to convert infix to postfix (RPN), then evaluate the postfix expression.

## Requirements
- Infix to postfix conversion via shunting-yard
- Operators: +, -, *, /, ^, %, parentheses
- Unary minus support
- Integer and float support
- Configurable operator precedence
- Error handling for division by zero and invalid tokens

## Implementation Notes
- Shunting-yard handles operator precedence and associativity
- Unary minus detected by context (after operator or at start)
- RPN evaluator uses a stack with named errors

## Test Strategy (Unit/Fuzz)
- Unit: basic ops, precedence, parentheses, unary minus, division by zero
- Fuzz: random valid expressions compared with Python's eval for a safe subset

## Edge Cases
- Consecutive unary minus (--5)
- Floating point numbers
- Nested parentheses
- Empty expressions
- Functions (sin, cos, sqrt)

## Failure Cases
- Division by zero → EvaluationError
- Mismatched parentheses → ValueError
- Unknown characters → ValueError
- Missing operands → EvaluationError

## Complexity
- O(n) time, O(n) space (n = expression length)

## Progression Path
- Basic: two-operand calculator
- Intermediate: shunting-yard + RPN evaluator
- Advanced: functions, variables
- Production: AST-based evaluator, JIT compilation

## Common Interview Follow-ups
- How would you add variable assignment?
- How would you add custom functions?
- How would you implement short-circuit evaluation?

## Possible Production Improvements
- AST-based compilation for faster evaluation
- Support for user-defined functions
- Type coercion (int vs float)
- Expression caching
