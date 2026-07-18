package url_parser

import (
	"fmt"
	"strconv"
	"strings"
)

// URL represents a parsed uniform resource locator.
type URL struct {
	Scheme   string
	Userinfo string
	Host     string
	Port     string
	Path     string
	Query    string
	Fragment string
}

// Parse parses rawURL into a URL struct. It does not use net/url.
// Supported forms:
//
//	scheme://userinfo@host:port/path?query#fragment
func Parse(rawURL string) (*URL, error) {
	if rawURL == "" {
		return nil, fmt.Errorf("empty URL")
	}
	u := &URL{}
	s := rawURL

	if i := strings.LastIndex(s, "#"); i >= 0 {
		u.Fragment = s[i+1:]
		s = s[:i]
	}

	if i := strings.LastIndex(s, "?"); i >= 0 {
		u.Query = s[i+1:]
		s = s[:i]
	}

	if i := strings.Index(s, "://"); i >= 0 {
		u.Scheme = s[:i]
		s = s[i+3:]
		if !isValidScheme(u.Scheme) {
			return nil, fmt.Errorf("invalid scheme %q", u.Scheme)
		}
	}

	var authority, path string
	if i := strings.Index(s, "/"); i >= 0 {
		authority = s[:i]
		path = s[i:]
	} else {
		authority = s
		path = ""
	}

	if authority != "" {
		parseAuthority(u, authority)
	}

	u.Path = path

	if u.Port == "" {
		switch u.Scheme {
		case "http":
			u.Port = "80"
		case "https":
			u.Port = "443"
		}
	}

	return u, nil
}

func parseAuthority(u *URL, authority string) {
	if i := strings.LastIndex(authority, "@"); i >= 0 {
		u.Userinfo = authority[:i]
		authority = authority[i+1:]
	}

	if len(authority) > 0 && authority[0] == '[' {
		if i := strings.Index(authority, "]"); i >= 0 {
			u.Host = authority[1:i]
			if len(authority) > i+1 && authority[i+1] == ':' {
				u.Port = authority[i+2:]
			}
			return
		}
	}

	if i := strings.LastIndex(authority, ":"); i >= 0 {
		u.Host = authority[:i]
		u.Port = authority[i+1:]
	} else {
		u.Host = authority
	}
}

func isValidScheme(scheme string) bool {
	if len(scheme) == 0 {
		return false
	}
	for i, c := range scheme {
		switch {
		case i == 0:
			if !isAlpha(c) {
				return false
			}
		case !isAlpha(c) && !isDigit(c) && c != '+' && c != '-' && c != '.':
			return false
		}
	}
	return true
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

// String reconstructs the URL string. Default ports are omitted for http/https.
func (u *URL) String() string {
	var b strings.Builder
	if u.Scheme != "" {
		b.WriteString(u.Scheme)
		b.WriteString("://")
	}
	if u.Userinfo != "" {
		b.WriteString(u.Userinfo)
		b.WriteString("@")
	}
	if strings.Contains(u.Host, ":") {
		b.WriteString("[")
		b.WriteString(u.Host)
		b.WriteString("]")
	} else {
		b.WriteString(u.Host)
	}
	if u.Port != "" && !isDefaultPort(u.Scheme, u.Port) {
		b.WriteString(":")
		b.WriteString(u.Port)
	}
	b.WriteString(u.Path)
	if u.Query != "" {
		b.WriteString("?")
		b.WriteString(u.Query)
	}
	if u.Fragment != "" {
		b.WriteString("#")
		b.WriteString(u.Fragment)
	}
	return b.String()
}

func isDefaultPort(scheme, port string) bool {
	switch scheme {
	case "http":
		return port == "80"
	case "https":
		return port == "443"
	}
	return false
}

// IsValid validates the URL components.
func (u *URL) IsValid() bool {
	if u.Scheme != "" && !isValidScheme(u.Scheme) {
		return false
	}
	if strings.Contains(u.Host, " ") {
		return false
	}
	if u.Port != "" {
		p, err := strconv.Atoi(u.Port)
		if err != nil || p < 0 || p > 65535 {
			return false
		}
	}
	return true
}

// DecodePercent decodes percent-encoded sequences in s.
func DecodePercent(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '%' && i+2 < len(s) {
			hi := unhex(s[i+1])
			lo := unhex(s[i+2])
			if hi >= 0 && lo >= 0 {
				b.WriteByte(byte(hi<<4 | lo))
				i += 2
				continue
			}
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func unhex(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c - 'a' + 10)
	case c >= 'A' && c <= 'F':
		return int(c - 'A' + 10)
	}
	return -1
}
