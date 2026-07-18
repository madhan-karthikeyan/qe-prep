package com.qe.test.expr;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;

import static org.junit.jupiter.api.Assertions.*;

@DisplayName("Evaluator")
class EvaluatorTest {

    private final Evaluator evaluator = new Evaluator();

    @ParameterizedTest
    @CsvSource({
            "1+2, 3.0",
            "3-1, 2.0",
            "2*3, 6.0",
            "10/2, 5.0",
            "2^3, 8.0",
            "10%3, 1.0",
            "1+2*3, 7.0",
            "(1+2)*3, 9.0",
            "2^3^2, 512.0",
            "3+4*2/(1-5), 1.0",
            "-5+3, -2.0",
            "--5, 5.0"
    })
    @DisplayName("evaluates expressions correctly")
    void evaluatesCorrectly(String expr, double expected) {
        assertEquals(expected, evaluator.evaluate(expr), 1e-10);
    }

    @Test
    @DisplayName("throws on division by zero")
    void divisionByZero() {
        assertThrows(ArithmeticException.class, () -> evaluator.evaluate("1/0"));
    }

    @Test
    @DisplayName("throws on modulo by zero")
    void moduloByZero() {
        assertThrows(ArithmeticException.class, () -> evaluator.evaluate("5%0"));
    }

    @Test
    @DisplayName("handles decimal numbers")
    void decimalNumbers() {
        assertEquals(3.5, evaluator.evaluate("1.5+2.0"), 1e-10);
        assertEquals(0.1, evaluator.evaluate("1.0/10"), 1e-10);
    }

    @Test
    @DisplayName("handles complex nested expressions")
    void complexNested() {
        assertEquals(10.0, evaluator.evaluate("(2+3)*(8-6)"), 1e-10);
        assertEquals(5.0, evaluator.evaluate("((1+2)*3-4)/1"), 1e-10);
    }
}
