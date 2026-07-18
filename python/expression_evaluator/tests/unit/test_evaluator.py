import unittest

from expression_evaluator.implementation.evaluator import EvaluationError, RPNEvaluator
from expression_evaluator.implementation.shunting_yard import ShuntingYard


class TestExpressionEvaluator(unittest.TestCase):
    def setUp(self) -> None:
        self.shunting = ShuntingYard()
        self.evaluator = RPNEvaluator()

    def eval_expr(self, expr: str) -> float:
        postfix = self.shunting.to_postfix(expr)
        return self.evaluator.evaluate(postfix)

    def test_basic_addition(self) -> None:
        self.assertEqual(self.eval_expr("2 + 3"), 5.0)

    def test_basic_subtraction(self) -> None:
        self.assertEqual(self.eval_expr("5 - 3"), 2.0)

    def test_basic_multiplication(self) -> None:
        self.assertEqual(self.eval_expr("4 * 3"), 12.0)

    def test_basic_division(self) -> None:
        self.assertEqual(self.eval_expr("10 / 2"), 5.0)

    def test_operator_precedence(self) -> None:
        self.assertEqual(self.eval_expr("2 + 3 * 4"), 14.0)
        self.assertEqual(self.eval_expr("2 * 3 + 4"), 10.0)

    def test_parentheses(self) -> None:
        self.assertEqual(self.eval_expr("(2 + 3) * 4"), 20.0)

    def test_nested_parentheses(self) -> None:
        self.assertEqual(self.eval_expr("((2 + 3) * 2)"), 10.0)

    def test_exponentiation(self) -> None:
        self.assertEqual(self.eval_expr("2 ^ 3"), 8.0)
        self.assertEqual(self.eval_expr("2 ^ 3 ^ 2"), 512.0)

    def test_modulo(self) -> None:
        self.assertEqual(self.eval_expr("10 % 3"), 1.0)

    def test_unary_minus(self) -> None:
        self.assertEqual(self.eval_expr("-5"), -5.0)
        self.assertEqual(self.eval_expr("-(3 + 2)"), -5.0)
        self.assertEqual(self.eval_expr("-3 * 2"), -6.0)
        self.assertEqual(self.eval_expr("2 * -3"), -6.0)

    def test_float_values(self) -> None:
        self.assertAlmostEqual(self.eval_expr("3.5 + 2.5"), 6.0)
        self.assertAlmostEqual(self.eval_expr("1.5 * 2"), 3.0)

    def test_division_by_zero(self) -> None:
        with self.assertRaises(EvaluationError):
            self.eval_expr("10 / 0")

    def test_modulo_by_zero(self) -> None:
        with self.assertRaises(EvaluationError):
            self.eval_expr("10 % 0")

    def test_mismatched_parentheses(self) -> None:
        with self.assertRaises(ValueError):
            self.shunting.to_postfix("(2 + 3")

    def test_double_negative(self) -> None:
        self.assertEqual(self.eval_expr("--5"), 5.0)

    def test_complex_expression(self) -> None:
        result = self.eval_expr("(2 + 3) * 4 ^ 2 / (1 + 1)")
        self.assertAlmostEqual(result, 40.0)

    def test_to_postfix_order(self) -> None:
        postfix = self.shunting.to_postfix("2 + 3 * 4")
        self.assertEqual(postfix, ["2", "3", "4", "*", "+"])

    def test_to_postfix_parentheses(self) -> None:
        postfix = self.shunting.to_postfix("(2 + 3) * 4")
        self.assertEqual(postfix, ["2", "3", "+", "4", "*"])

    def test_to_postfix_unary(self) -> None:
        postfix = self.shunting.to_postfix("-5")
        self.assertEqual(postfix, ["5", "_u"])

    def test_to_postfix_double_unary(self) -> None:
        postfix = self.shunting.to_postfix("--5")
        self.assertEqual(postfix, ["5", "_u", "_u"])

    def test_unknown_token(self) -> None:
        with self.assertRaises(ValueError):
            self.shunting.to_postfix("2 @ 3")

    def test_functions(self) -> None:
        postfix = self.shunting.to_postfix("sin 0")
        result = self.evaluator.evaluate(postfix)
        self.assertAlmostEqual(result, 0.0)

        postfix = self.shunting.to_postfix("cos 0")
        result = self.evaluator.evaluate(postfix)
        self.assertAlmostEqual(result, 1.0)

        postfix = self.shunting.to_postfix("sqrt 9")
        result = self.evaluator.evaluate(postfix)
        self.assertAlmostEqual(result, 3.0)


if __name__ == "__main__":
    unittest.main()
