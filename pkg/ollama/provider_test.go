package ollama

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"shelloma/pkg/config"
)

func TestNewClient(t *testing.T) {
	// 1. Success case: Connection ok, models found, model pre-selected
	t.Run("Success Preselected Model", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/tags" {
				resp := TagsResponse{
					Models: []ModelInfo{
						{Name: "qwen2.5-coder:1.5b"},
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

		client, err := NewClient(cfg)
		if err != nil {
			t.Fatalf("Erro inesperado: %v", err)
		}

		if client.GetModel() != "qwen2.5-coder:1.5b" {
			t.Errorf("Modelo inesperado: %s", client.GetModel())
		}
	})

	// 2. Success case: Connection ok, models found, auto-select model
	t.Run("Success Auto Select Model", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/tags" {
				resp := TagsResponse{
					Models: []ModelInfo{
						{Name: "llama3:latest"},
					},
				}
				_ = json.NewEncoder(w).Encode(resp)
			}
		}))
		defer server.Close()

		cfg := config.Config{
			OllamaURL: server.URL,
			Model:     "", // empty triggers auto-select
		}

		client, err := NewClient(cfg)
		if err != nil {
			t.Fatalf("Erro inesperado: %v", err)
		}

		if client.GetModel() != "llama3:latest" {
			t.Errorf("Esperava llama3:latest, obteve %s", client.GetModel())
		}
	})

	// 3. Fail case: Offline / server error
	t.Run("Ollama Offline", func(t *testing.T) {
		// Use a non-routable port or closed port URL
		cfg := config.Config{
			OllamaURL: "http://127.0.0.1:9999",
			Model:     "qwen2.5-coder:1.5b",
		}

		_, err := NewClient(cfg)
		if err == nil {
			t.Errorf("Esperava erro ao tentar conectar com servidor offline")
		}
		if !strings.Contains(err.Error(), "Could not connect to Ollama service") {
			t.Errorf("Mensagem de erro inesperada: %v", err)
		}
	})

	// 4. Fail case: No models found (empty tags response)
	t.Run("No Models Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/tags" {
				resp := TagsResponse{
					Models: []ModelInfo{}, // Empty list
				}
				_ = json.NewEncoder(w).Encode(resp)
			}
		}))
		defer server.Close()

		cfg := config.Config{
			OllamaURL: server.URL,
			Model:     "qwen2.5-coder:1.5b",
		}

		_, err := NewClient(cfg)
		if err == nil {
			t.Errorf("Esperava erro por falta de modelos")
		}
		if !strings.Contains(err.Error(), "No models found in Ollama") {
			t.Errorf("Mensagem de erro inesperada: %v", err)
		}
	})
}
