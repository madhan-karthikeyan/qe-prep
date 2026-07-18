package com.qe.test.expr;

import java.util.ArrayDeque;
import java.util.Deque;
import java.util.List;

public class Evaluator {
    private final ShuntingYard shuntingYard = new ShuntingYard();

    public double evaluate(String expression) {
        List<String> rpn = shuntingYard.toRpn(expression);
        return evaluateRpn(rpn);
    }

    public double evaluateRpn(List<String> rpn) {
        Deque<Double> stack = new ArrayDeque<>();

        for (String token : rpn) {
            switch (token) {
                case "+" -> {
                    double b = stack.pop();
                    double a = stack.pop();
                    stack.push(a + b);
                }
                case "-" -> {
                    double b = stack.pop();
                    double a = stack.pop();
                    stack.push(a - b);
                }
                case "*" -> {
                    double b = stack.pop();
                    double a = stack.pop();
                    stack.push(a * b);
                }
                case "/" -> {
                    double b = stack.pop();
                    double a = stack.pop();
                    if (b == 0.0) {
                        throw new ArithmeticException("Division by zero");
                    }
                    stack.push(a / b);
                }
                case "^" -> {
                    double b = stack.pop();
                    double a = stack.pop();
                    stack.push(Math.pow(a, b));
                }
                case "%" -> {
                    double b = stack.pop();
                    double a = stack.pop();
                    if (b == 0.0) {
                        throw new ArithmeticException("Modulo by zero");
                    }
                    stack.push(a % b);
                }
                case "_" -> stack.push(-stack.pop());
                default -> {
                    try {
                        stack.push(Double.parseDouble(token));
                    } catch (NumberFormatException e) {
                        throw new IllegalArgumentException("Unknown token: " + token, e);
                    }
                }
            }
        }

        if (stack.size() != 1) {
            throw new IllegalStateException("Invalid expression: stack has " + stack.size() + " elements");
        }

        return stack.pop();
    }
}
