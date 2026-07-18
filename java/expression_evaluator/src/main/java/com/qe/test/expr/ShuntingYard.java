package com.qe.test.expr;

import java.util.ArrayDeque;
import java.util.ArrayList;
import java.util.Deque;
import java.util.List;
import java.util.Map;

public class ShuntingYard {
    private static final Map<String, Integer> PRECEDENCE = Map.of(
            "+", 2, "-", 2,
            "*", 3, "/", 3, "%", 3,
            "^", 4,
            "_", 5
    );

    private static final Map<String, Associativity> ASSOCIATIVITY = Map.of(
            "+", Associativity.LEFT, "-", Associativity.LEFT,
            "*", Associativity.LEFT, "/", Associativity.LEFT, "%", Associativity.LEFT,
            "^", Associativity.RIGHT,
            "_", Associativity.RIGHT
    );

    private enum Associativity { LEFT, RIGHT }

    public List<String> toRpn(String expression) {
        if (expression == null || expression.isBlank()) {
            throw new IllegalArgumentException("Expression must not be null or blank");
        }

        List<String> tokens = tokenize(expression);
        List<String> output = new ArrayList<>();
        Deque<String> operators = new ArrayDeque<>();

        for (int i = 0; i < tokens.size(); i++) {
            String token = tokens.get(i);

            if (isNumber(token)) {
                output.add(token);
            } else if (isFunction(token)) {
                operators.push(token);
            } else if (token.equals(",")) {
                while (!operators.isEmpty() && !operators.peek().equals("(")) {
                    output.add(operators.pop());
                }
                if (operators.isEmpty()) {
                    throw new IllegalArgumentException("Mismatched parentheses or function call");
                }
            } else if (isOperator(token)) {
                if (token.equals("-")) {
                    boolean isUnary = i == 0
                            || tokens.get(i - 1).equals("(")
                            || isOperator(tokens.get(i - 1))
                            || tokens.get(i - 1).equals(",");
                    if (isUnary) {
                        operators.push("_");
                        continue;
                    }
                }
                while (!operators.isEmpty() && isOperator(operators.peek())) {
                    String top = operators.peek();
                    if ((ASSOCIATIVITY.get(token) == Associativity.LEFT
                            && PRECEDENCE.get(token) <= PRECEDENCE.get(top))
                            || (ASSOCIATIVITY.get(token) == Associativity.RIGHT
                            && PRECEDENCE.get(token) < PRECEDENCE.get(top))) {
                        output.add(operators.pop());
                    } else {
                        break;
                    }
                }
                operators.push(token);
            } else if (token.equals("(")) {
                operators.push(token);
            } else if (token.equals(")")) {
                while (!operators.isEmpty() && !operators.peek().equals("(")) {
                    output.add(operators.pop());
                }
                if (operators.isEmpty()) {
                    throw new IllegalArgumentException("Mismatched parentheses");
                }
                operators.pop(); // discard "("
                if (!operators.isEmpty() && isFunction(operators.peek())) {
                    output.add(operators.pop());
                }
            }
        }

        while (!operators.isEmpty()) {
            String op = operators.pop();
            if (op.equals("(") || op.equals(")")) {
                throw new IllegalArgumentException("Mismatched parentheses");
            }
            output.add(op);
        }

        return output;
    }

    List<String> tokenize(String expr) {
        List<String> tokens = new ArrayList<>();
        StringBuilder current = new StringBuilder();

        for (int i = 0; i < expr.length(); i++) {
            char c = expr.charAt(i);

            if (Character.isWhitespace(c)) {
                continue;
            }

            if (Character.isDigit(c) || c == '.') {
                current.append(c);
            } else {
                if (!current.isEmpty()) {
                    tokens.add(current.toString());
                    current.setLength(0);
                }
                if (isOperatorChar(c) || c == '(' || c == ')' || c == ',') {
                    tokens.add(String.valueOf(c));
                } else if (Character.isLetter(c)) {
                    current.append(c);
                } else {
                    throw new IllegalArgumentException("Unexpected character: " + c);
                }
            }
        }

        if (!current.isEmpty()) {
            tokens.add(current.toString());
        }

        return tokens;
    }

    private boolean isNumber(String token) {
        try {
            Double.parseDouble(token);
            return true;
        } catch (NumberFormatException e) {
            return false;
        }
    }

    private boolean isOperator(String token) {
        return PRECEDENCE.containsKey(token);
    }

    private boolean isFunction(String token) {
        return Character.isLetter(token.charAt(0));
    }

    private boolean isOperatorChar(char c) {
        return "+-*/%^".indexOf(c) >= 0;
    }
}
