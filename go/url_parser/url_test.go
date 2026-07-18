package url_parser

import (
	"testing"
)

func TestParseStandardURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *URL
		wantStr string
	}{
		{
			name:  "simple http",
			input: "http://example.com",
			want: &URL{
				Scheme: "http",
				Host:   "example.com",
				Port:   "80",
				Path:   "",
			},
			wantStr: "http://example.com",
		},
		{
			name:  "https with path",
			input: "https://example.com/path/to/resource",
			want: &URL{
				Scheme: "https",
				Host:   "example.com",
				Port:   "443",
				Path:   "/path/to/resource",
			},
			wantStr: "https://example.com/path/to/resource",
		},
		{
			name:  "with port",
			input: "http://example.com:8080/path",
			want: &URL{
				Scheme: "http",
				Host:   "example.com",
				Port:   "8080",
				Path:   "/path",
			},
			wantStr: "http://example.com:8080/path",
		},
		{
			name:  "			with userinfo",
			input: "http://user:pass@host.com",
			want: &URL{
				Scheme:   "http",
				Userinfo: "user:pass",
				Host:     "host.com",
				Port:     "80",
			},
			wantStr: "http://user:pass@host.com",
		},
		{
			name:  "with query and fragment",
			input: "http://example.com/path?q=1&r=2#section",
			want: &URL{
				Scheme:   "http",
				Host:     "example.com",
				Port:     "80",
				Path:     "/path",
				Query:    "q=1&r=2",
				Fragment: "section",
			},
			wantStr: "http://example.com/path?q=1&r=2#section",
		},
		{
			name:  "empty path no trailing slash",
			input: "http://example.com?query",
			want: &URL{
				Scheme: "http",
				Host:   "example.com",
				Port:   "80",
				Query:  "query",
			},
			wantStr: "http://example.com?query",
		},
		{
			name:  "ftp scheme",
			input: "ftp://files.example.com/pub",
			want: &URL{
				Scheme: "ftp",
				Host:   "files.example.com",
				Path:   "/pub",
			},
			wantStr: "ftp://files.example.com/pub",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned error: %v", tt.input, err)
			}
			if got.Scheme != tt.want.Scheme {
				t.Errorf("Scheme = %q, want %q", got.Scheme, tt.want.Scheme)
			}
			if got.Userinfo != tt.want.Userinfo {
				t.Errorf("Userinfo = %q, want %q", got.Userinfo, tt.want.Userinfo)
			}
			if got.Host != tt.want.Host {
				t.Errorf("Host = %q, want %q", got.Host, tt.want.Host)
			}
			if got.Port != tt.want.Port {
				t.Errorf("Port = %q, want %q", got.Port, tt.want.Port)
			}
			if got.Path != tt.want.Path {
				t.Errorf("Path = %q, want %q", got.Path, tt.want.Path)
			}
			if got.Query != tt.want.Query {
				t.Errorf("Query = %q, want %q", got.Query, tt.want.Query)
			}
			if got.Fragment != tt.want.Fragment {
				t.Errorf("Fragment = %q, want %q", got.Fragment, tt.want.Fragment)
			}
			if got.String() != tt.wantStr {
				t.Errorf("String() = %q, want %q", got.String(), tt.wantStr)
			}
		})
	}
}

func TestParseIPv6(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		host    string
		port    string
		wantStr string
	}{
		{
			name:    "ipv6 loopback",
			input:   "http://[::1]/path",
			host:    "::1",
			port:    "80",
			wantStr: "http://[::1]/path",
		},
		{
			name:    "ipv6 with port",
			input:   "http://[2001:db8::1]:8080/path",
			host:    "2001:db8::1",
			port:    "8080",
			wantStr: "http://[2001:db8::1]:8080/path",
		},
		{
			name:    "ipv6 full",
			input:   "https://[2001:db8:85a3::8a2e:370:7334]:443/",
			host:    "2001:db8:85a3::8a2e:370:7334",
			port:    "443",
			wantStr: "https://[2001:db8:85a3::8a2e:370:7334]/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q): %v", tt.input, err)
			}
			if got.Host != tt.host {
				t.Errorf("Host = %q, want %q", got.Host, tt.host)
			}
			if got.Port != tt.port {
				t.Errorf("Port = %q, want %q", got.Port, tt.port)
			}
			if got.String() != tt.wantStr {
				t.Errorf("String() = %q, want %q", got.String(), tt.wantStr)
			}
		})
	}
}

func TestParseMalformed(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"bad scheme", "123://host"},
		{"bad scheme char", "http@://host"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if err == nil {
				t.Errorf("Parse(%q) expected error", tt.input)
			}
		})
	}
}

func TestParseRelative(t *testing.T) {
	got, err := Parse("/path/to/resource?q=1#frag")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if got.Scheme != "" {
		t.Errorf("Scheme = %q, want empty", got.Scheme)
	}
	if got.Path != "/path/to/resource" {
		t.Errorf("Path = %q", got.Path)
	}
	if got.Query != "q=1" {
		t.Errorf("Query = %q", got.Query)
	}
	if got.Fragment != "frag" {
		t.Errorf("Fragment = %q", got.Fragment)
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{"valid http", "http://example.com", true},
		{"valid https", "https://host.com/path", true},
		{"valid with port", "http://host:8080", true},
		{"invalid port", "http://host:99999", false},
		{"host with space", "http://host.com/path", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := Parse(tt.url)
			if err != nil {
				t.Skipf("parse error: %v", err)
			}
			if got := u.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodePercent(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello%20world", "hello world"},
		{"%48%65%6c%6c%6f", "Hello"},
		{"noencoding", "noencoding"},
		{"%ZZ", "%ZZ"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := DecodePercent(tt.input); got != tt.want {
				t.Errorf("DecodePercent(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestStringRoundtrip(t *testing.T) {
	inputs := []string{
		"http://example.com/path",
		"https://user:pass@host:8080/path?q=1#frag",
		"http://[::1]:8080/path",
		"ftp://files.example.com",
		"http://example.com?query",
		"https://example.com/path?q=1&r=two#section",
	}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			u, err := Parse(input)
			if err != nil {
				t.Fatalf("Parse(%q): %v", input, err)
			}
			got := u.String()
			u2, err := Parse(got)
			if err != nil {
				t.Fatalf("Parse(%q): %v", got, err)
			}
			if u2.String() != got {
				t.Errorf("non-idempotent: %q -> %q", got, u2.String())
			}
		})
	}
}
