package ui

import (
	"strings"
	"testing"

	"shelloma/pkg/i18n"
)

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
