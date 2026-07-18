package config_parser

import (
	"fmt"
	"strings"
)

// ParseINI parses a .ini formatted string and returns a map of sections to
// key-value pairs. Supports comments (# and ;), blank lines, and [section]
// headers.
func ParseINI(data string) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)
	lines := strings.Split(data, "\n")
	var currentSection string

	for i, line := range lines {
		line = strings.TrimRight(line, "\r")
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") {
			if !strings.HasSuffix(line, "]") {
				return nil, fmt.Errorf("line %d: malformed section header", i+1)
			}
			currentSection = line[1 : len(line)-1]
			if _, ok := result[currentSection]; !ok {
				result[currentSection] = make(map[string]string)
			}
			continue
		}

		eq := strings.Index(line, "=")
		if eq < 0 {
			return nil, fmt.Errorf("line %d: missing '='", i+1)
		}

		key := strings.TrimSpace(line[:eq])
		value := strings.TrimSpace(line[eq+1:])

		if key == "" {
			return nil, fmt.Errorf("line %d: empty key", i+1)
		}

		if currentSection == "" {
			return nil, fmt.Errorf("line %d: key %q outside of section", i+1, key)
		}

		result[currentSection][key] = value
	}

	return result, nil
}
