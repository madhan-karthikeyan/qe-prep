package config_parser

import (
	"fmt"
	"strings"
)

// ParseEnv parses a .env formatted string and returns a map of key-value
// pairs. Supports comments (#), quoted values, escape sequences, and variable
// substitution ($VAR and ${VAR}).
func ParseEnv(data string) (map[string]string, error) {
	result := make(map[string]string)
	lines := strings.Split(data, "\n")

	for i, line := range lines {
		line = strings.TrimRight(line, "\r")
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		eq := strings.Index(line, "=")
		if eq < 0 {
			continue
		}

		key := strings.TrimSpace(line[:eq])
		if key == "" {
			return nil, fmt.Errorf("line %d: empty key", i+1)
		}

		rawValue := line[eq+1:]
		value := strings.TrimSpace(rawValue)

		if len(value) >= 2 {
			quote := value[0]
			if quote == '"' || quote == '\'' {
				if value[len(value)-1] != quote {
					return nil, fmt.Errorf("line %d: unterminated string", i+1)
				}
				value = value[1 : len(value)-1]
				if quote == '"' {
					value = unescapeDouble(value)
				}
			}
		}

		value = expandVars(value, result)
		result[key] = value
	}

	return result, nil
}

func unescapeDouble(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case '\\':
				b.WriteByte('\\')
			case '"':
				b.WriteByte('"')
			case '\'':
				b.WriteByte('\'')
			default:
				b.WriteByte(s[i+1])
			}
			i++
		} else {
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

func expandVars(s string, vars map[string]string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) && s[i+1] == '$' {
			b.WriteByte('$')
			i++
			continue
		}
		if s[i] == '$' && i+1 < len(s) {
			if s[i+1] == '{' {
				end := strings.IndexByte(s[i+2:], '}')
				if end >= 0 {
					name := s[i+2 : i+2+end]
					if val, ok := vars[name]; ok {
						b.WriteString(val)
					}
					i += 2 + end
					continue
				}
			} else if isIdentStart(s[i+1]) {
				j := i + 1
				for j < len(s) && isIdentChar(s[j]) {
					j++
				}
				name := s[i+1 : j]
				if val, ok := vars[name]; ok {
					b.WriteString(val)
				}
				i = j - 1
				continue
			}
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func isIdentStart(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isIdentChar(c byte) bool {
	return isIdentStart(c) || (c >= '0' && c <= '9')
}
