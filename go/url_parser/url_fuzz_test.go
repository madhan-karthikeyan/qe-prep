package url_parser

import (
	"testing"
)

func FuzzParse(f *testing.F) {
	seeds := []string{
		"http://example.com",
		"https://user:pass@host:8080/path?q=1#frag",
		"http://[::1]:8080/path",
		"ftp://files.example.com",
		"/path/to/resource",
		"http://example.com/path?q=a&r=b#section",
		"http://192.168.1.1:8080",
		"",
		"not-a-url",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		u, err := Parse(s)
		if err != nil {
			return
		}
		reconstructed := u.String()
		u2, err := Parse(reconstructed)
		if err != nil {
			t.Errorf("reparse of %q failed: %v", reconstructed, err)
			return
		}
		if u2.String() != reconstructed {
			t.Errorf("non-idempotent: %q -> %q", reconstructed, u2.String())
		}
	})
}
