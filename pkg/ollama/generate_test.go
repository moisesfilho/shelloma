package ollama

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"shelloma/pkg/config"
	"shelloma/pkg/sysinfo"
)

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
			_ = json.NewEncoder(w).Encode(resp)
		case "/api/generate":
			resp := GenerateResponse{
				Response: "ls -la",
				Done:     true,
			}
			_ = json.NewEncoder(w).Encode(resp)
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

	winSysCtx := sysinfo.SystemContext{
		OS:         "windows",
		DistroName: "Windows",
		Shell:      "powershell.exe",
		WorkingDir: "C:\\Users\\Test",
	}

	winCmd, err := client.GenerateCommand(winSysCtx, "listar arquivos", 0.1)
	if err != nil {
		t.Fatalf("Erro ao gerar comando no Windows: %v", err)
	}
	if winCmd != "ls -la" {
		t.Errorf("Esperava resposta mockada 'ls -la', obteve %q", winCmd)
	}
}
