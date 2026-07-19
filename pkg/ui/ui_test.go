package ui

import (
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

func TestExecuteCommandSuccess(t *testing.T) {
	trans := i18n.GetTranslations("en")
	exitCode, output, err := ExecuteCommand("echo 'shelloma test'", trans)
	if err != nil {
		t.Fatalf("ExecuteCommand falhou com erro: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Esperava exitCode 0, obteve %d", exitCode)
	}

	if !strings.Contains(output, "shelloma test") {
		t.Errorf("Saída esperada 'shelloma test' não encontrada na saída: %q", output)
	}
}

func TestExecuteCommandFailure(t *testing.T) {
	trans := i18n.GetTranslations("en")
	exitCode, _, _ := ExecuteCommand("ls /caminho_inexistente_123_456_shelloma_test", trans)
	if exitCode == 0 {
		t.Errorf("Esperava exitCode != 0 para caminho inexistente, obteve %d", exitCode)
	}
}

func TestPromptActionCopyAndOtherInputs(t *testing.T) {
	trans := i18n.GetTranslations("pt")

	tests := []struct {
		input    string
		expected Action
	}{
		{"c\n", ActionCopy},
		{"copy\n", ActionCopy},
		{"copiar\n", ActionCopy},
		{"q\n", ActionQuit},
		{"sair\n", ActionQuit},
		{"y\n", ActionExecute},
		{"\n", ActionExecute},
		{"e\n", ActionExplain},
		{"m\n", ActionEdit},
	}

	for _, tt := range tests {
		reader := strings.NewReader(tt.input)
		action := PromptActionWithReader(reader, trans)
		if action != tt.expected {
			t.Errorf("Para a entrada %q, esperava ação %d, obteve %d", tt.input, tt.expected, action)
		}
	}
}
