package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"shelloma/pkg/config"
	"shelloma/pkg/i18n"
	"shelloma/pkg/ollama"
)

func TestHandleConfigCommand(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "shelloma_config_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	t.Setenv("XDG_CONFIG_HOME", tempDir)
	t.Setenv("HOME", tempDir)
	t.Setenv("LocalAppData", tempDir)

	trans := i18n.GetTranslations("en")
	cfg := config.Config{
		OllamaURL:             "http://localhost:11434",
		Model:                 "qwen2.5-coder:1.5b",
		Language:              "en",
		Temperature:           0.1,
		AutoExecute:           false,
		DangerousCommands:     []string{"rm", "dd"},
		DisableDangerousCheck: false,
	}

	// 1. Test Config Set Model
	t.Run("Config Set Model", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		HandleConfigCommand(cfg, []string{"set", "model", "llama3:latest"}, trans)

		w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		os.Stdout = oldStdout
		out := buf.String()

		if !strings.Contains(out, "Configuration saved successfully") {
			t.Errorf("Esperava mensagem de configuração salva, obteve %q", out)
		}

		// Valida se a configuração de fato salvou (precisa ler via arquivo agora)
		newCfg, err := config.LoadConfig()
		if err != nil {
			t.Fatalf("failed to load saved config: %v", err)
		}
		if newCfg.Model != "llama3:latest" {
			t.Errorf("Modelo salvo incorreto, esperava llama3:latest, obteve %s", newCfg.Model)
		}
	})

	// 2. Test Config Set URL e Outros Parâmetros
	t.Run("Config Set Language and URL", func(t *testing.T) {
		newCfg, _ := config.LoadConfig()
		HandleConfigCommand(newCfg, []string{"set", "url", "http://my-ollama:11434"}, trans)
		newCfg, _ = config.LoadConfig()
		if newCfg.OllamaURL != "http://my-ollama:11434" {
			t.Errorf("URL salva incorreta: %s", newCfg.OllamaURL)
		}

		HandleConfigCommand(newCfg, []string{"set", "lang", "pt"}, trans)
		newCfg, _ = config.LoadConfig()
		if newCfg.Language != "pt" {
			t.Errorf("Linguagem salva incorreta: %s", newCfg.Language)
		}

		HandleConfigCommand(newCfg, []string{"set", "dangerous", "rm,dd,mkfs"}, trans)
		newCfg, _ = config.LoadConfig()
		if len(newCfg.DangerousCommands) != 3 || newCfg.DangerousCommands[2] != "mkfs" {
			t.Errorf("Comandos perigosos incorretos: %v", newCfg.DangerousCommands)
		}

		HandleConfigCommand(newCfg, []string{"set", "disable_dangerous_check", "true"}, trans)
		newCfg, _ = config.LoadConfig()
		if !newCfg.DisableDangerousCheck {
			t.Errorf("DisableDangerousCheck deveria ser true")
		}
	})

	// 3. Test Show Config
	t.Run("Show Config Details", func(t *testing.T) {
		newCfg, _ := config.LoadConfig()
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		HandleConfigCommand(newCfg, []string{}, trans)

		w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		os.Stdout = oldStdout
		out := buf.String()

		if !strings.Contains(out, "ollama_url") || !strings.Contains(out, "model") {
			t.Errorf("Esperava detalhes da configuração impressos, obteve %q", out)
		}
	})
}

func TestHandleModelsCommand(t *testing.T) {
	trans := i18n.GetTranslations("en")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/tags" {
			resp := ollama.TagsResponse{
				Models: []ollama.ModelInfo{
					{Name: "qwen2.5-coder:1.5b"},
					{Name: "llama3:latest"},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	cfg := config.Config{
		OllamaURL: server.URL,
		Model:     "qwen2.5-coder:1.5b",
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	HandleModelsCommand(cfg, trans)

	w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	os.Stdout = oldStdout
	out := buf.String()

	if !strings.Contains(out, "Models installed") {
		t.Errorf("Esperava cabeçalho 'Models installed', obteve %q", out)
	}
	if !strings.Contains(out, "qwen2.5-coder:1.5b") || !strings.Contains(out, "llama3:latest") {
		t.Errorf("Esperava que os modelos estivessem listados, obteve %q", out)
	}
	if !strings.Contains(out, "(active)") {
		t.Errorf("Esperava que o modelo ativo estivesse indicado, obteve %q", out)
	}
}
