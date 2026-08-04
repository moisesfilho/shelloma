package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LearnedCommand struct {
	Command   string    `json:"command"`
	Help      string    `json:"help"`
	Timestamp time.Time `json:"timestamp"`
}

func GetLearnedDir() (string, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(filepath.Dir(configPath), "learned")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func SaveLearnedCommand(cmdName string, helpText string) error {
	dir, err := GetLearnedDir()
	if err != nil {
		return err
	}
	filePath := filepath.Join(dir, strings.ToLower(cmdName)+".json")
	data := LearnedCommand{
		Command:   cmdName,
		Help:      helpText,
		Timestamp: time.Now(),
	}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, jsonData, 0600)
}

func LoadLearnedCommands() ([]LearnedCommand, error) {
	dir, err := GetLearnedDir()
	if err != nil {
		return nil, err
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var list []LearnedCommand
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
			filePath := filepath.Join(dir, f.Name())
			bytes, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}
			var cmd LearnedCommand
			if err := json.Unmarshal(bytes, &cmd); err == nil {
				list = append(list, cmd)
			}
		}
	}
	return list, nil
}
