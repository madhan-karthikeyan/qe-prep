package com.qe.test.expr;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

@DisplayName("ShuntingYard")
class ShuntingYardTest {

    private final ShuntingYard sy = new ShuntingYard();

    @Test
    @DisplayName("converts simple addition to RPN")
    void simpleAddition() {
        assertEquals(List.of("1", "2", "+"), sy.toRpn("1+2"));
    }

    @Test
    @DisplayName("respects operator precedence")
    void precedence() {
        assertEquals(List.of("1", "2", "3", "*", "+"), sy.toRpn("1+2*3"));
        assertEquals(List.of("1", "2", "+", "3", "*"), sy.toRpn("(1+2)*3"));
    }

    @Test
    @DisplayName("handles parentheses")
    void parentheses() {
        assertEquals(List.of("3", "4", "+", "5", "+"), sy.toRpn("(3+4)+5"));
    }

    @Test
    @DisplayName("handles right-associative power operator")
    void rightAssociativePower() {
        assertEquals(List.of("2", "3", "^"), sy.toRpn("2^3"));
    }

    @Test
    @DisplayName("tokenizes correctly")
    void tokenize() {
        assertEquals(List.of("3", "+", "4.5", "*", "(", "2", "-", "1", ")"), sy.tokenize("3+4.5*(2-1)"));
    }

    @Test
    @DisplayName("rejects null expression")
    void rejectsNull() {
        assertThrows(IllegalArgumentException.class, () -> sy.toRpn(null));
    }

    @Test
    @DisplayName("rejects blank expression")
    void rejectsBlank() {
        assertThrows(IllegalArgumentException.class, () -> sy.toRpn("  "));
    }

    @Test
    @DisplayName("rejects mismatched parentheses")
    void rejectsMismatchedParentheses() {
        assertThrows(IllegalArgumentException.class, () -> sy.toRpn("((1+2)"));
        assertThrows(IllegalArgumentException.class, () -> sy.toRpn("(1+2))"));
    }

    @Test
    @DisplayName("handles unary minus")
    void unaryMinus() {
        var rpn = sy.toRpn("-5+3");
        assertTrue(rpn.contains("_"), "RPN should contain unary minus marker");
    }
}
