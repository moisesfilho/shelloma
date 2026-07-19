package ui

import (
	"testing"
)

func TestEncodeBase64(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "aGVsbG8="},
		{"ls -la", "bHMgLWxh"},
		{"shelloma", "c2hlbGxvbWE="},
	}

	for _, tt := range tests {
		got := encodeBase64(tt.input)
		if got != tt.expected {
			t.Errorf("encodeBase64(%q) = %q; esperava %q", tt.input, got, tt.expected)
		}
	}
}
