package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	OllamaURL   string  `json:"ollama_url"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	AutoExecute bool    `json:"auto_execute"`
	Language    string  `json:"language"`
}

func DefaultConfig() Config {
	return Config{
		OllamaURL:   "http://localhost:11434",
		Model:       "",   // Vazio para detecção automática
		Temperature: 0.1,
		AutoExecute: false,
		Language:    "en", // Idioma padrão: Inglês
	}
}

func GetConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		home, hErr := os.UserHomeDir()
		if hErr != nil {
			return "", err
		}
		return filepath.Join(home, ".config", "shelloma", "config.json"), nil
	}
	return filepath.Join(dir, "shelloma", "config.json"), nil
}

func LoadConfig() (Config, error) {
	cfg := DefaultConfig()

	// 1. Ler a configuração do sistema em /etc/shelloma/config.json (se existir no Linux/macOS)
	if runtime.GOOS != "windows" {
		systemPath := "/etc/shelloma/config.json"
		if data, err := os.ReadFile(systemPath); err == nil {
			var sysCfg Config
			if err := json.Unmarshal(data, &sysCfg); err == nil {
				if sysCfg.OllamaURL != "" {
					cfg.OllamaURL = sysCfg.OllamaURL
				}
				if sysCfg.Model != "" {
					cfg.Model = sysCfg.Model
				}
				if sysCfg.Language != "" {
					cfg.Language = sysCfg.Language
				}
				if sysCfg.Temperature > 0 {
					cfg.Temperature = sysCfg.Temperature
				}
				cfg.AutoExecute = sysCfg.AutoExecute
			}
		}
	}

	// 2. Ler a configuração individual do usuário em ~/.config/shelloma/config.json (precedência maior)
	userPath, err := GetConfigPath()
	if err == nil {
		data, err := os.ReadFile(userPath)
		if err == nil {
			var userCfg Config
			if err := json.Unmarshal(data, &userCfg); err == nil {
				if userCfg.OllamaURL != "" {
					cfg.OllamaURL = userCfg.OllamaURL
				}
				if userCfg.Model != "" {
					cfg.Model = userCfg.Model
				}
				if userCfg.Language != "" {
					cfg.Language = userCfg.Language
				}
				if userCfg.Temperature > 0 {
					cfg.Temperature = userCfg.Temperature
				}
				cfg.AutoExecute = userCfg.AutoExecute
			}
		} else if os.IsNotExist(err) {
			// Se o arquivo do usuário ainda não existir, salva mesclando com /etc/shelloma/config.json
			_ = SaveConfig(cfg)
		}
	}

	if cfg.Language == "" {
		cfg.Language = "en"
	}

	return cfg, nil
}

func SaveConfig(cfg Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
