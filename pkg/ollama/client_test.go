package ollama

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"shelloma/pkg/config"
	"shelloma/pkg/sysinfo"
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

func TestSelectBestModel(t *testing.T) {
	models := []string{"llama2:latest", "qwen2.5-coder:1.5b", "mistral:latest"}
	best := selectBestModel(models)
	if best != "qwen2.5-coder:1.5b" {
		t.Errorf("Esperava qwen2.5-coder:1.5b, obteve %s", best)
	}

	fallbackModels := []string{"custom-model:v1"}
	bestFallback := selectBestModel(fallbackModels)
	if bestFallback != "custom-model:v1" {
		t.Errorf("Esperava custom-model:v1, obteve %s", bestFallback)
	}
}

func TestOllamaClientMock(t *testing.T) {
	// Mock HTTP Server para simular o Ollama
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/tags":
			resp := TagsResponse{
				Models: []ModelInfo{
					{Name: "qwen2.5-coder:1.5b"},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/api/generate":
			resp := GenerateResponse{
				Response: "ls -la",
				Done:     true,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	cfg := config.Config{
		OllamaURL:   server.URL,
		Model:       "qwen2.5-coder:1.5b",
		Temperature: 0.1,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Erro ao criar cliente Ollama: %v", err)
	}

	models, err := client.ListModels()
	if err != nil {
		t.Fatalf("Erro ao listar modelos: %v", err)
	}
	if len(models) != 1 || models[0] != "qwen2.5-coder:1.5b" {
		t.Errorf("ListModels retornou inesperado: %v", models)
	}

	sysCtx := sysinfo.SystemContext{
		OS:         "linux",
		DistroName: "Ubuntu",
		Shell:      "bash",
		WorkingDir: "/tmp",
	}

	cmd, err := client.GenerateCommand(sysCtx, "listar arquivos", 0.1)
	if err != nil {
		t.Fatalf("Erro ao gerar comando: %v", err)
	}

	if cmd != "ls -la" {
		t.Errorf("Esperava comando 'ls -la', obteve %q", cmd)
	}
}
