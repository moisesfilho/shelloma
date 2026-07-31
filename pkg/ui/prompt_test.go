package ui

import (
	"fmt"
	"strings"
	"testing"

	"shelloma/pkg/i18n"
	"shelloma/pkg/sysinfo"
)

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
		{"p\n", ActionNewPrompt},
		{"prompt\n", ActionNewPrompt},
		{"new\n", ActionNewPrompt},
		{"novo\n", ActionNewPrompt},
		{"nuevo\n", ActionNewPrompt},
		{"r\n", ActionRefine},
		{"refine\n", ActionRefine},
		{"refinar\n", ActionRefine},
		{"complement\n", ActionRefine},
		{"complementar\n", ActionRefine},
		{"a\n", ActionAdjustPrompt},
		{"adjust\n", ActionAdjustPrompt},
		{"ajustar\n", ActionAdjustPrompt},
		{"alterar\n", ActionAdjustPrompt},
	}

	for _, tt := range tests {
		reader := strings.NewReader(tt.input)
		action := PromptActionWithReader(reader, trans)
		if action != tt.expected {
			t.Errorf("Para a entrada %q, esperava ação %d, obteve %d", tt.input, tt.expected, action)
		}
	}
}

func TestStripANSI(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"\x1b[1mHello\x1b[0m", "Hello"},
		{"Plain Text", "Plain Text"},
		{"\x1b[31mRed\x1b[32mGreen\x1b[0m", "RedGreen"},
		{"", ""},
	}
	for _, tt := range tests {
		got := stripANSI(tt.input)
		if got != tt.expected {
			t.Errorf("stripANSI(%q) = %q, expected %q", tt.input, got, tt.expected)
		}
	}
}

func TestGetTerminalSize(t *testing.T) {
	rows, cols := getTerminalSize()
	if rows <= 0 || cols <= 0 {
		t.Errorf("getTerminalSize retornou valores inválidos: rows=%d, cols=%d", rows, cols)
	}
}

func TestTerminalSetupResetAndDraw(t *testing.T) {
	t.Log("Testing Terminal Setup, Reset, and Draw functions")
	sysCtx := sysinfo.SystemContext{
		WorkingDir: "/tmp",
		OS:         "linux",
		DistroName: "Ubuntu",
		Shell:      "bash",
	}
	trans := i18n.GetTranslations("pt")

	SetupTerminal(sysCtx, "qwen2.5-coder:1.5b", "1.3.0", trans)
	DrawLegendAtBottom("Test Legend")
	ClearLegendAtBottom()
	ResetTerminal()
}

func TestEditCommand(t *testing.T) {
	trans := i18n.GetTranslations("en")

	// 1. User enters new command
	{
		oldStdin := StdinReader
		StdinReader = strings.NewReader("mkdir -p /tmp/new_test_dir\n")
		defer func() { StdinReader = oldStdin }()

		cmd := EditCommand("ls -la", trans)
		if cmd != "mkdir -p /tmp/new_test_dir" {
			t.Errorf("Esperava 'mkdir -p /tmp/new_test_dir', obteve %q", cmd)
		}
	}

	// 2. User enters empty line (should return current/original command)
	{
		oldStdin := StdinReader
		StdinReader = strings.NewReader("\n")
		defer func() { StdinReader = oldStdin }()

		cmd := EditCommand("ls -la", trans)
		if cmd != "ls -la" {
			t.Errorf("Esperava retorno do comando original 'ls -la', obteve %q", cmd)
		}
	}
}

func TestPromptSecurityWord(t *testing.T) {
	trans := i18n.GetTranslations("en")

	// 1. User enters correct word
	{
		oldStdin := StdinReader
		StdinReader = strings.NewReader("CONFIRM\n")
		defer func() { StdinReader = oldStdin }()

		success := PromptSecurityWord("CONFIRM", trans)
		if !success {
			t.Errorf("Esperava que a palavra de segurança CONFIRM fosse considerada correta")
		}
	}

	// 2. User enters incorrect word
	{
		oldStdin := StdinReader
		StdinReader = strings.NewReader("INCORRECT\n")
		defer func() { StdinReader = oldStdin }()

		success := PromptSecurityWord("CONFIRM", trans)
		if success {
			t.Errorf("Esperava que a palavra incorreta falhasse")
		}
	}
}

type errorReader struct{}

func (errorReader) Read(_ []byte) (n int, err error) {
	return 0, fmt.Errorf("simulated read error")
}

func TestReadFilteredInputError(t *testing.T) {
	reader := errorReader{}
	input, err := ReadFilteredInput(reader)
	if err == nil {
		t.Errorf("Esperava erro ao ler do reader com erro, obteve input: %q", input)
	}
	if !strings.Contains(err.Error(), "simulated read error") {
		t.Errorf("Erro inesperado: %v", err)
	}
}


