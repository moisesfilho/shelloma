package ollama

import (
	"testing"
)

func TestCleanCommandOutput(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"```bash\nls -l\n```", "ls -l"},
		{"`ls -la`", "ls -la"},
		{"  mkdir -p /tmp/test  \n", "mkdir -p /tmp/test"},
		{"```\necho hello\n```", "echo hello"},
	}

	for _, tt := range tests {
		got := cleanCommandOutput(tt.input)
		if got != tt.expected {
			t.Errorf("cleanCommandOutput(%q) = %q; esperava %q", tt.input, got, tt.expected)
		}
	}
}

func TestCleanJSONOutput(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"```json\n{\"success\": true}\n```", "{\"success\": true}"},
		{"```\n{\"success\": false}\n```", "{\"success\": false}"},
		{"  {\"reason\": \"ok\"}  ", "{\"reason\": \"ok\"}"},
	}

	for _, tt := range tests {
		got := cleanJSONOutput(tt.input)
		if got != tt.expected {
			t.Errorf("cleanJSONOutput(%q) = %q; esperava %q", tt.input, got, tt.expected)
		}
	}
}

func TestIsValidShellCommand(t *testing.T) {
	validCmds := []string{
		"ls -la /home",
		"mkdir -p /tmp/foo",
		"cat /etc/os-release",
		"touch file.txt",
		"sudo apt update",
	}

	for _, cmd := range validCmds {
		if !isValidShellCommand(cmd) {
			t.Errorf("Esperava que %q fosse considerado um comando válido", cmd)
		}
	}

	invalidCmds := []string{
		"",
		"Você pode tentar criar o arquivo",
		"Verifique se o diretório existe",
		"Caso este arquivo não exista, crie",
	}

	for _, cmd := range invalidCmds {
		if isValidShellCommand(cmd) {
			t.Errorf("Esperava que %q fosse considerado um comando INVÁLIDO", cmd)
		}
	}
}
