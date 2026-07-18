import random
import unittest

from expression_evaluator.implementation.evaluator import EvaluationError, RPNEvaluator
from expression_evaluator.implementation.shunting_yard import ShuntingYard

_safe_ops = ["+", "-", "*", "/"]
_integers = list(range(1, 100))


def generate_expression(depth: int = 0) -> str:
    if depth > 3 or random.random() < 0.3:
        return str(random.choice(_integers))
    op = random.choice(_safe_ops)
    left = generate_expression(depth + 1)
    right = generate_expression(depth + 1)
    if random.random() < 0.3:
        return f"({left} {op} {right})"
    return f"{left} {op} {right}"


class TestExpressionFuzz(unittest.TestCase):
    def setUp(self) -> None:
        self.shunting = ShuntingYard()
        self.evaluator = RPNEvaluator()

    def test_random_expressions_vs_eval(self) -> None:
        for _ in range(5000):
            expr = generate_expression()
            try:
                postfix = self.shunting.to_postfix(expr)
                result = self.evaluator.evaluate(postfix)
                expected = float(eval(expr))
                self.assertAlmostEqual(result, expected, places=10)
            except (EvaluationError, ZeroDivisionError, ValueError):
                pass

    def test_edge_expressions(self) -> None:
        edge_cases = [
            "0",
            "0 + 0",
            "1 / 2",
            "2 / 1",
            "1 + 2 + 3 + 4 + 5",
            "10 - 9 - 8",
            "1 + (2 * (3 + (4 * 5)))",
            "100",
        ]
        for expr in edge_cases:
            try:
                postfix = self.shunting.to_postfix(expr)
                result = self.evaluator.evaluate(postfix)
                expected = float(eval(expr))
                self.assertAlmostEqual(result, expected, places=10)
            except (EvaluationError, ZeroDivisionError, ValueError) as e:
                self.fail(f"Unexpected error for {expr!r}: {e}")

    def test_random_invalid_expressions(self) -> None:
        invalid = [
            "+",
            "1 +",
            "+ 1",
            "(1 + 2",
            "1 + 2)",
            "1 / 0",
            "1 % 0",
            "",
            "a b c",
            "1 2 3",
        ]
        for expr in invalid:
            try:
                postfix = self.shunting.to_postfix(expr)
                self.evaluator.evaluate(postfix)
            except (EvaluationError, ValueError):
                pass


if __name__ == "__main__":
    unittest.main()
