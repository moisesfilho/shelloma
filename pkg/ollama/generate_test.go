package ollama

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"shelloma/pkg/config"
	"shelloma/pkg/i18n"
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

func TestExplainCommand(t *testing.T) {
	var responseText string
	statusCode := http.StatusOK

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(statusCode)
		if statusCode == http.StatusOK {
			resp := GenerateResponse{
				Response: responseText,
				Done:     true,
			}
			_ = json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		Model:      "test-model",
		HTTPClient: server.Client(),
		Trans:      i18n.GetTranslations("en"),
	}

	// 1. Success case
	responseText = "Explicando o comando ls: lista arquivos."
	exp, err := client.ExplainCommand("ls")
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}
	if exp != "Explicando o comando ls: lista arquivos." {
		t.Errorf("Resposta incorreta, obteve %q", exp)
	}

	// 2. Error case (JSON Inválido)
	statusCode = http.StatusInternalServerError
	_, err = client.ExplainCommand("ls")
	if err == nil {
		t.Errorf("Esperava erro ao falhar a chamada HTTP")
	}
}

func TestGenerateFixCommand(t *testing.T) {
	var responseText string
	statusCode := http.StatusOK

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(statusCode)
		if statusCode == http.StatusOK {
			resp := GenerateResponse{
				Response: responseText,
				Done:     true,
			}
			_ = json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		Model:      "test-model",
		HTTPClient: server.Client(),
	}

	sysCtx := sysinfo.SystemContext{OS: "linux"}

	// 1. Success case
	responseText = "mkdir -p test"
	cmd, err := client.GenerateFixCommand(sysCtx, "cd test", "no such file")
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}
	if cmd != "mkdir -p test" {
		t.Errorf("Comando de correção incorreto: %q", cmd)
	}

	// 2. Invalid shell command returned (should return empty string)
	responseText = "Você deve criar o diretório."
	cmd, err = client.GenerateFixCommand(sysCtx, "cd test", "no such file")
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}
	if cmd != "" {
		t.Errorf("Esperava string vazia para comando inválido, obteve %q", cmd)
	}
}

func TestGenerateAlternativeCommand(t *testing.T) {
	var responseText string
	statusCode := http.StatusOK

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(statusCode)
		if statusCode == http.StatusOK {
			resp := GenerateResponse{
				Response: responseText,
				Done:     true,
			}
			_ = json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		Model:      "test-model",
		HTTPClient: server.Client(),
	}

	sysCtx := sysinfo.SystemContext{OS: "linux"}

	// 1. Success case
	responseText = "find . -name '*.txt'"
	cmd, err := client.GenerateAlternativeCommand(sysCtx, "listar arquivos txt", "ls *.txt", "no matches")
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}
	if cmd != "find . -name '*.txt'" {
		t.Errorf("Comando alternativo incorreto: %q", cmd)
	}

	// 2. "I_CANNOT_HELP" case (should return empty string)
	responseText = "I_CANNOT_HELP"
	cmd, err = client.GenerateAlternativeCommand(sysCtx, "...", "...", "...")
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}
	if cmd != "" {
		t.Errorf("Esperava string vazia para I_CANNOT_HELP, obteve %q", cmd)
	}

	// 3. Invalid shell command returned (should return empty string)
	responseText = "Você pode tentar criar o arquivo"
	cmd, err = client.GenerateAlternativeCommand(sysCtx, "...", "...", "...")
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}
	if cmd != "" {
		t.Errorf("Esperava string vazia para comando inválido, obteve %q", cmd)
	}
}

func TestGenerateRefinedCommand(t *testing.T) {
	var responseText string
	statusCode := http.StatusOK

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(statusCode)
		if statusCode == http.StatusOK {
			resp := GenerateResponse{
				Response: responseText,
				Done:     true,
			}
			_ = json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		Model:      "test-model",
		HTTPClient: server.Client(),
		Rules:      []string{"Use sudo"},
	}

	sysCtx := sysinfo.SystemContext{OS: "linux"}

	// 1. Success case
	responseText = "sudo ls"
	cmd, err := client.GenerateRefinedCommand(sysCtx, "listar", "ls", "com privilégios", 0.2)
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}
	if cmd != "sudo ls" {
		t.Errorf("Comando refinado incorreto: %q", cmd)
	}
}

