package ollama

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"shelloma/pkg/i18n"
	"shelloma/pkg/sysinfo"
)

func TestContainsErrorKeywords(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"Everything is fine", false},
		{"permission denied to access this folder", true},
		{"error: something went wrong", true},
		{"fatal: database connection failed", true},
		{"Comando não encontrado no sistema", true},
		{"comando normal", false},
		{"permissão negada ao executar o script", true},
	}

	for _, tt := range tests {
		got := containsErrorKeywords(tt.input)
		if got != tt.expected {
			t.Errorf("containsErrorKeywords(%q) = %v; esperava %v", tt.input, got, tt.expected)
		}
	}
}

func TestAnalyzeExecutionResult(t *testing.T) {
	sysCtx := sysinfo.SystemContext{
		OS:         "linux",
		DistroName: "Ubuntu",
		Shell:      "bash",
	}

	trans := i18n.GetTranslations("en")

	// 1. Success case without calling server (Exit code 0 and no error keywords)
	t.Run("Immediate Success", func(t *testing.T) {
		client := &Client{Trans: trans}
		res, err := client.AnalyzeExecutionResult("ls", 0, "file1.txt\nfile2.txt", sysCtx)
		if err != nil {
			t.Fatalf("Erro inesperado: %v", err)
		}
		if !res.Success {
			t.Errorf("Esperava Success=true")
		}
		if res.Reason != "Comando executado com sucesso" {
			t.Errorf("Esperava razão de sucesso, obteve %q", res.Reason)
		}
	})

	// 2. Server Mock for failures / analysis
	t.Run("Server Analysis Success and Failures", func(t *testing.T) {
		var returnedResponse string
		statusCode := http.StatusOK

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(statusCode)
			if statusCode == http.StatusOK {
				resp := GenerateResponse{
					Response: returnedResponse,
					Done:     true,
				}
				_ = json.NewEncoder(w).Encode(resp)
			} else {
				_, _ = w.Write([]byte("internal error"))
			}
		}))
		defer server.Close()

		client := &Client{
			BaseURL:    server.URL,
			Model:      "test-model",
			HTTPClient: server.Client(),
			Trans:      trans,
		}

		// Scenario A: Server returns valid JSON analysis showing failure
		returnedResponse = `{"success": false, "reason": "folder does not exist", "suggested_command": "mkdir -p test"}`
		res, err := client.AnalyzeExecutionResult("cd test", 1, "no such file or directory", sysCtx)
		if err != nil {
			t.Fatalf("Erro inesperado no cenário A: %v", err)
		}
		if res.Success {
			t.Errorf("Esperava Success=false")
		}
		if res.Reason != "folder does not exist" {
			t.Errorf("Razão incorreta: %q", res.Reason)
		}
		if res.SuggestedCommand != "mkdir -p test" {
			t.Errorf("Comando sugerido incorreto: %q", res.SuggestedCommand)
		}

		// Scenario B: Server returns invalid JSON
		returnedResponse = `invalid json output`
		res, err = client.AnalyzeExecutionResult("cd test", 1, "no such file or directory", sysCtx)
		if err != nil {
			t.Fatalf("Erro inesperado no cenário B: %v", err)
		}
		if res.Success {
			t.Errorf("Esperava Success=false")
		}
		if !strings.Contains(res.Reason, "Finished with exit code 1") {
			t.Errorf("Esperava fallback na razão, obteve %q", res.Reason)
		}

		// Scenario C: Server returns suggested command that is invalid
		returnedResponse = `{"success": false, "reason": "folder does not exist", "suggested_command": "Você precisa criar a pasta antes"}`
		res, err = client.AnalyzeExecutionResult("cd test", 1, "no such file or directory", sysCtx)
		if err != nil {
			t.Fatalf("Erro inesperado no cenário C: %v", err)
		}
		if res.SuggestedCommand != "" {
			t.Errorf("Esperava que o comando sugerido inválido fosse limpo, obteve %q", res.SuggestedCommand)
		}

		// Scenario D: Server returns non-200 status
		statusCode = http.StatusInternalServerError
		res, err = client.AnalyzeExecutionResult("cd test", 1, "no such file or directory", sysCtx)
		if err != nil {
			t.Fatalf("Erro inesperado no cenário D: %v", err)
		}
		if res.Success {
			t.Errorf("Esperava Success=false")
		}
		if !strings.Contains(res.Reason, "Execution error (exit code 1)") {
			t.Errorf("Esperava erro de execução, obteve %q", res.Reason)
		}
	})
}
