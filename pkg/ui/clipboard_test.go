package ui

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"shelloma/pkg/i18n"
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

func TestQuoteCmdArg(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "'hello'"},
		{"hello'world", "'hello''world'"},
		{"ls -la", "'ls -la'"},
	}

	for _, tt := range tests {
		got := quoteCmdArg(tt.input)
		if got != tt.expected {
			t.Errorf("quoteCmdArg(%q) = %q; esperava %q", tt.input, got, tt.expected)
		}
	}
}

func TestIsCommandAvailable(t *testing.T) {
	// A maioria dos sistemas Unix e Windows tem o comando 'go' ou 'echo' disponível no PATH
	if !isCommandAvailable("go") && !isCommandAvailable("echo") {
		t.Log("Nem 'go' nem 'echo' foram encontrados no PATH, pulando parte do teste")
	}

	if isCommandAvailable("comando_absolutamente_inexistente_123_456") {
		t.Errorf("Esperava que o comando inexistente não estivesse disponível")
	}
}

func TestCopyToClipboard(t *testing.T) {
	trans := i18n.GetTranslations("en")

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Chamamos CopyToClipboard. Mesmo se falhar pela ausência de utilitários no ambiente de teste,
	// ele deve imprimir a sequência OSC 52 no stdout antes de tentar executar os comandos do SO.
	_ = CopyToClipboard("ls -la", trans)

	w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	os.Stdout = oldStdout
	out := buf.String()

	// OSC 52 para "ls -la" é: \033]52;c;bHMgLWxh\a
	expectedOSC := "\033]52;c;bHMgLWxh\a"
	if !strings.Contains(out, expectedOSC) {
		t.Errorf("Esperava que a saída contivesse a sequência OSC 52 %q, obteve %q", expectedOSC, out)
	}
}

