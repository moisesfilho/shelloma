package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Config struct {
	OllamaURL             string   `json:"ollama_url"`
	Model                 string   `json:"model"`
	Temperature           float64  `json:"temperature"`
	AutoExecute           bool     `json:"auto_execute"`
	Language              string   `json:"language"`
	DangerousCommands     []string `json:"dangerous_commands"`
	DisableDangerousCheck bool     `json:"disable_dangerous_check"`
}

func DefaultConfig() Config {
	return Config{
		OllamaURL:             "http://localhost:11434",
		Model:                 "", // Vazio para detecção automática
		Temperature:           0.1,
		AutoExecute:           false,
		Language:              "en", // Idioma padrão: Inglês
		DangerousCommands:     []string{"rm", "dd", "mkfs", "shred", "chmod", "chown", "Remove-Item", "del", "rd", "rmdir", "format", "Format-Volume"},
		DisableDangerousCheck: false,
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
				if len(sysCfg.DangerousCommands) > 0 {
					cfg.DangerousCommands = sysCfg.DangerousCommands
				}
				cfg.DisableDangerousCheck = sysCfg.DisableDangerousCheck
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
				if len(userCfg.DangerousCommands) > 0 {
					cfg.DangerousCommands = userCfg.DangerousCommands
				}
				cfg.DisableDangerousCheck = userCfg.DisableDangerousCheck
			}
		} else if os.IsNotExist(err) {
			// Se o arquivo do usuário ainda não existir, salva mesclando com /etc/shelloma/config.json
			_ = SaveConfig(cfg)
		}
	}

	if cfg.Language == "" {
		cfg.Language = "en"
	}

	if len(cfg.DangerousCommands) == 0 {
		cfg.DangerousCommands = DefaultConfig().DangerousCommands
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

// CheckDangerous checks if the command string uses any command from the dangerous list.
// It splits the command by shell syntax delimiters and checks for exact matches of words.
func CheckDangerous(cmd string, dangerousList []string) (bool, string) {
	words := strings.FieldsFunc(cmd, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == ';' || r == '|' || r == '&' || r == '`' || r == '(' || r == ')' || r == '$' || r == '<' || r == '>' || r == '"' || r == '\''
	})
	for _, word := range words {
		for _, dangerous := range dangerousList {
			if word == dangerous {
				return true, dangerous
			}
		}
	}
	return false, ""
}
