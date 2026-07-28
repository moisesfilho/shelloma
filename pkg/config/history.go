package config

import (
	"os"
	"path/filepath"
	"strings"
)

func GetHistoryPath() (string, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(configPath), "history.txt"), nil
}

func LoadHistory() ([]string, error) {
	path, err := GetHistoryPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	var history []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			history = append(history, trimmed)
		}
	}
	return history, nil
}

func SaveHistory(history []string) error {
	path, err := GetHistoryPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	if len(history) > 1000 {
		history = history[len(history)-1000:]
	}
	content := strings.Join(history, "\n") + "\n"
	return os.WriteFile(path, []byte(content), 0600)
}

func AddToHistory(prompt string) {
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return
	}
	history, err := LoadHistory()
	if err != nil {
		return
	}
	if len(history) > 0 && history[len(history)-1] == prompt {
		return
	}
	history = append(history, prompt)
	_ = SaveHistory(history)
}
