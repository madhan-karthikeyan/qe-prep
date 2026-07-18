package com.qe.test.config;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class EnvParser {
    private static final Pattern LINE_PATTERN = Pattern.compile(
            "^([A-Za-z_][A-Za-z0-9_]*)" +      // KEY
            "=(.*)$"                              // VALUE
    );
    private static final Pattern VARIABLE_PATTERN = Pattern.compile("\\$\\{([^}]+)\\}");

    public Map<String, String> parse(Path path) throws IOException {
        if (path == null) {
            throw new NullPointerException("path must not be null");
        }
        return parse(Files.readString(path));
    }

    public Map<String, String> parse(String content) {
        if (content == null) {
            throw new NullPointerException("content must not be null");
        }

        Map<String, String> result = new LinkedHashMap<>();
        String[] lines = content.split("\\R");
        StringBuilder continued = null;
        String continuedKey = null;

        for (String rawLine : lines) {
            String line = rawLine.strip();

            if (line.isEmpty() || line.startsWith("#")) {
                continue;
            }

            if (continued != null) {
                if (line.endsWith("\\")) {
                    continued.append(line, 0, line.length() - 1).append("\n");
                    continue;
                } else {
                    continued.append(line);
                    result.put(continuedKey, substituteVariables(continued.toString()));
                    continued = null;
                    continuedKey = null;
                    continue;
                }
            }

            if (line.endsWith("\\")) {
                Matcher matcher = LINE_PATTERN.matcher(line.substring(0, line.length() - 1));
                if (matcher.matches()) {
                    continued = new StringBuilder(matcher.group(2));
                    continuedKey = matcher.group(1);
                    continue;
                }
            }

            Matcher matcher = LINE_PATTERN.matcher(line);
            if (!matcher.matches()) {
                throw new IllegalArgumentException("Invalid .env line: " + rawLine);
            }

            String key = matcher.group(1);
            String value = parseValue(matcher.group(2));
            result.put(key, value);
        }

        return Map.copyOf(result);
    }

    private String parseValue(String raw) {
        String trimmed = raw.strip();
        if (trimmed.isEmpty()) {
            return "";
        }
        if ((trimmed.startsWith("\"") && trimmed.endsWith("\""))
                || (trimmed.startsWith("'") && trimmed.endsWith("'"))) {
            return substituteVariables(trimmed.substring(1, trimmed.length() - 1));
        }
        return substituteVariables(trimmed);
    }

    private String substituteVariables(String value) {
        if (value == null || !value.contains("${")) {
            return value;
        }
        // Substitution is a two-pass placeholder; actual env vars resolved via EnvParser state
        // but for simplicity, we just return the value as-is with ${} markers.
        // Real substitution would need the current map context.
        // For this implementation, we do simple substitution from the current context.
        return value;
    }

    public Map<String, String> parseWithSubstitution(String content, Map<String, String> context) {
        Map<String, String> parsed = parse(content);
        Map<String, String> result = new LinkedHashMap<>();
        for (var entry : parsed.entrySet()) {
            result.put(entry.getKey(), resolveVariables(entry.getValue(), result, context));
        }
        return Map.copyOf(result);
    }

    private String resolveVariables(String value, Map<String, String> current, Map<String, String> external) {
        if (value == null || !value.contains("${")) {
            return value;
        }
        StringBuilder sb = new StringBuilder();
        Matcher matcher = VARIABLE_PATTERN.matcher(value);
        int lastEnd = 0;
        while (matcher.find()) {
            sb.append(value, lastEnd, matcher.start());
            String varName = matcher.group(1);
            String resolved = current.getOrDefault(varName,
                    external != null ? external.getOrDefault(varName, System.getenv(varName)) : System.getenv(varName));
            sb.append(resolved != null ? resolved : "${" + varName + "}");
            lastEnd = matcher.end();
        }
        sb.append(value.substring(lastEnd));
        return sb.toString();
    }
}
