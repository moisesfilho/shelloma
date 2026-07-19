package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.OllamaURL != "http://localhost:11434" {
		t.Errorf("URL padrão incorreta: %s", cfg.OllamaURL)
	}

	if cfg.Temperature != 0.1 {
		t.Errorf("Temperatura padrão incorreta: %f", cfg.Temperature)
	}

	if cfg.AutoExecute != false {
		t.Error("AutoExecute padrão deveria ser false")
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "shelloma_test_config")
	if err != nil {
		t.Fatalf("Erro ao criar diretório temporário: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Sobrescrever temporariamente os diretórios de configuração do usuário
	t.Setenv("HOME", tempDir)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tempDir, ".config"))
	t.Setenv("AppData", filepath.Join(tempDir, "AppData", "Roaming"))

	cfg := Config{
		OllamaURL:   "http://127.0.0.1:11434",
		Model:       "test-coder-model",
		Temperature: 0.2,
		AutoExecute: true,
	}

	err = SaveConfig(cfg)
	if err != nil {
		t.Fatalf("Erro ao salvar configuração: %v", err)
	}

	configFilePath := filepath.Join(tempDir, ".config", "shelloma", "config.json")
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		t.Fatalf("Arquivo de configuração não foi criado em %s", configFilePath)
	}

	loadedCfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Erro ao carregar configuração: %v", err)
	}

	if loadedCfg.OllamaURL != cfg.OllamaURL {
		t.Errorf("Esperava OllamaURL %s, obteve %s", cfg.OllamaURL, loadedCfg.OllamaURL)
	}

	if loadedCfg.Model != cfg.Model {
		t.Errorf("Esperava Model %s, obteve %s", cfg.Model, loadedCfg.Model)
	}

	if loadedCfg.Temperature != cfg.Temperature {
		t.Errorf("Esperava Temperature %f, obteve %f", cfg.Temperature, loadedCfg.Temperature)
	}

	if loadedCfg.AutoExecute != cfg.AutoExecute {
		t.Errorf("Esperava AutoExecute %t, obteve %t", cfg.AutoExecute, loadedCfg.AutoExecute)
	}
}
