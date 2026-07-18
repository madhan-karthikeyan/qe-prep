

_operators: dict[str, tuple[int, str]] = {
    "+": (2, "L"),
    "-": (2, "L"),
    "*": (3, "L"),
    "/": (3, "L"),
    "%": (3, "L"),
    "^": (4, "R"),
    "_u": (6, "R"),
}

_functions: set[str] = {
    "sin", "cos", "sqrt", "abs", "max", "min",
}


class ShuntingYard:
    def __init__(
        self,
        custom_precedence: dict[str, tuple[int, str]] | None = None,
    ) -> None:
        self._ops = dict(_operators)
        if custom_precedence:
            self._ops.update(custom_precedence)

    def tokenize(self, expression: str) -> list[str]:
        tokens: list[str] = []
        i = 0
        while i < len(expression):
            ch = expression[i]
            if ch.isspace():
                i += 1
                continue
            if ch.isdigit() or ch == ".":
                num = ch
                i += 1
                while i < len(expression) and (
                    expression[i].isdigit() or expression[i] == "."
                ):
                    num += expression[i]
                    i += 1
                tokens.append(num)
                continue
            if ch.isalpha():
                token = ch
                i += 1
                while i < len(expression) and (
                    expression[i].isalpha() or expression[i].isdigit()
                ):
                    token += expression[i]
                    i += 1
                tokens.append(token)
                continue
            if ch in self._ops or ch in "(),":
                tokens.append(ch)
                i += 1
                continue
            raise ValueError(f"Unexpected character: {ch!r}")
        return tokens

    def to_postfix(self, expression: str) -> list[str]:
        tokens = self.tokenize(expression)
        output: list[str] = []
        op_stack: list[str] = []
        last_token: str | None = None

        for token in tokens:
            if self._is_number(token):
                output.append(token)
            elif token.isalpha() and token not in _functions:
                output.append(token)
            elif token in _functions:
                op_stack.append(token)
            elif token == ",":
                while op_stack and op_stack[-1] != "(":
                    output.append(op_stack.pop())
                if not op_stack or op_stack[-1] == "(":
                    pass
            elif token == "(":
                op_stack.append(token)
            elif token == ")":
                while op_stack and op_stack[-1] != "(":
                    output.append(op_stack.pop())
                if not op_stack:
                    raise ValueError("Mismatched parentheses")
                op_stack.pop()
                if op_stack and op_stack[-1] in _functions:
                    output.append(op_stack.pop())
            elif token in self._ops:
                if (
                    token == "-"
                    and (
                        last_token is None
                        or last_token in self._ops
                        or last_token == "("
                    )
                ):
                    op_stack.append("_u")
                else:
                    while (
                        op_stack
                        and op_stack[-1] != "("
                        and op_stack[-1] not in _functions
                        and self._precedence(op_stack[-1], token)
                    ):
                        output.append(op_stack.pop())
                    op_stack.append(token)
            else:
                raise ValueError(f"Unknown token: {token!r}")
            last_token = token

        while op_stack:
            top = op_stack.pop()
            if top == "(":
                raise ValueError("Mismatched parentheses")
            output.append(top)

        return output

    def _is_number(self, token: str) -> bool:
        try:
            float(token)
            return True
        except ValueError:
            return False

    def _precedence(self, op1: str, op2: str) -> bool:
        p1, a1 = self._ops[op1]
        p2, a2 = self._ops[op2]
        if a1 == "L":
            return p1 >= p2
        return p1 > p2
