import math
from typing import Any


class EvaluationError(Exception):
    pass


class RPNEvaluator:
    def __init__(self) -> None:
        self._functions: dict[str, Any] = {
            "sin": math.sin,
            "cos": math.cos,
            "sqrt": math.sqrt,
            "abs": abs,
            "max": max,
            "min": min,
        }

    def evaluate(self, postfix: list[str]) -> float:
        stack: list[float] = []
        for token in postfix:
            if token == "_u":
                if not stack:
                    raise EvaluationError("Missing operand for unary minus")
                stack.append(-stack.pop())
            elif token == "+":
                b, a = self._pop_two(stack)
                stack.append(a + b)
            elif token == "-":
                b, a = self._pop_two(stack)
                stack.append(a - b)
            elif token == "*":
                b, a = self._pop_two(stack)
                stack.append(a * b)
            elif token == "/":
                b, a = self._pop_two(stack)
                if b == 0:
                    raise EvaluationError("Division by zero")
                stack.append(a / b)
            elif token == "%":
                b, a = self._pop_two(stack)
                if b == 0:
                    raise EvaluationError("Modulo by zero")
                stack.append(a % b)
            elif token == "^":
                b, a = self._pop_two(stack)
                stack.append(a ** b)
            elif token in self._functions:
                if not stack:
                    raise EvaluationError(f"Missing argument for {token}")
                stack.append(self._functions[token](stack.pop()))
            else:
                try:
                    stack.append(float(token))
                except ValueError:
                    raise EvaluationError(f"Invalid token: {token!r}")

        if len(stack) != 1:
            raise EvaluationError(
                f"Invalid expression: {len(stack)} values remaining"
            )
        return stack[0]

    def _pop_two(self, stack: list[float]) -> tuple[float, float]:
        if len(stack) < 2:
            raise EvaluationError("Missing operands")
        return stack.pop(), stack.pop()
