package com.qe.test.config;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.regex.Pattern;

public class IniParser {
    private static final Pattern SECTION_PATTERN = Pattern.compile("^\\[([^]]+)\\]\\s*$");
    private static final Pattern KEY_VALUE_PATTERN = Pattern.compile("^([^=#]+)=(.*)$");

    public Map<String, Map<String, String>> parse(Path path) throws IOException {
        if (path == null) {
            throw new NullPointerException("path must not be null");
        }
        return parse(Files.readString(path));
    }

    public Map<String, Map<String, String>> parse(String content) {
        if (content == null) {
            throw new NullPointerException("content must not be null");
        }

        Map<String, Map<String, String>> result = new LinkedHashMap<>();
        String currentSection = null;
        StringBuilder multiLineValue = null;
        String multiLineKey = null;

        for (String rawLine : content.split("\\R")) {
            String line = rawLine.strip();

            if (line.isEmpty() || line.startsWith(";") || line.startsWith("#")) {
                continue;
            }

            if (multiLineValue != null) {
                multiLineValue.append("\n");
                if (line.endsWith("\\")) {
                    multiLineValue.append(line, 0, line.length() - 1);
                    continue;
                } else {
                    multiLineValue.append(line);
                    result.get(currentSection).put(multiLineKey, multiLineValue.toString().strip());
                    multiLineValue = null;
                    multiLineKey = null;
                    continue;
                }
            }

            var sectionMatcher = SECTION_PATTERN.matcher(line);
            if (sectionMatcher.matches()) {
                currentSection = sectionMatcher.group(1);
                result.putIfAbsent(currentSection, new LinkedHashMap<>());
                continue;
            }

            var kvMatcher = KEY_VALUE_PATTERN.matcher(line);
            if (kvMatcher.matches()) {
                String key = kvMatcher.group(1).strip();
                String value = kvMatcher.group(2).strip();

                if (currentSection == null) {
                    currentSection = "";
                    result.putIfAbsent(currentSection, new LinkedHashMap<>());
                }

                if (value.endsWith("\\")) {
                    multiLineValue = new StringBuilder(value.substring(0, value.length() - 1));
                    multiLineKey = key;
                } else {
                    if ((value.startsWith("\"") && value.endsWith("\""))
                            || (value.startsWith("'") && value.endsWith("'"))) {
                        value = value.substring(1, value.length() - 1);
                    }
                    result.get(currentSection).put(key, value);
                }
            }
        }

        return result;
    }
}
