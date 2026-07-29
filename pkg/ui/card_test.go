package ui

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"shelloma/pkg/i18n"
)

func TestPrintBanner(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintBanner("qwen2.5-coder:1.5b", "pt")

	w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	os.Stdout = oldStdout
	out := buf.String()

	if !strings.Contains(out, "[Shelloma]") {
		t.Errorf("Esperava que o banner contivesse '[Shelloma]', obteve: %q", out)
	}
	if !strings.Contains(out, "Model: qwen2.5-coder:1.5b") {
		t.Errorf("Esperava que o banner contivesse o modelo, obteve: %q", out)
	}
	if !strings.Contains(out, "Lang: PT") {
		t.Errorf("Esperava que o banner contivesse a linguagem, obteve: %q", out)
	}
}

func TestPrintCommandCard(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintCommandCard("ls -la\nrm -rf /tmp/test")

	w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	os.Stdout = oldStdout
	out := buf.String()

	if !strings.Contains(out, "┌") || !strings.Contains(out, "└") {
		t.Errorf("Esperava que o card tivesse bordas de caixa, obteve: %q", out)
	}
	if !strings.Contains(out, "ls -la") || !strings.Contains(out, "rm -rf /tmp/test") {
		t.Errorf("Esperava que o card contivesse os comandos impressos, obteve: %q", out)
	}
}

func TestPrintDangerousWarning(t *testing.T) {
	trans := i18n.GetTranslations("pt")
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintDangerousWarning("rm -rf", trans)

	w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	os.Stdout = oldStdout
	out := buf.String()

	if !strings.Contains(out, "⚠️") {
		t.Errorf("Esperava ícone de aviso de perigo, obteve: %q", out)
	}
	if !strings.Contains(out, "rm -rf") {
		t.Errorf("Esperava que o aviso de perigo mencionasse o comando, obteve: %q", out)
	}
}
