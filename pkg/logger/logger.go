package logger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type LogEntry struct {
	Timestamp              string  `json:"timestamp"`
	UserQuery              string  `json:"user_query"`
	SuggestedCommand       string  `json:"suggested_command"`
	UserAction             string  `json:"user_action"`
	ExitCode               int     `json:"exit_code"`
	CommandOutput          string  `json:"command_output"`
	WorkingDir             string  `json:"working_dir"`
	User                   string  `json:"user"`
	OS                     string  `json:"os"`
	OllamaURL              string  `json:"ollama_url"`
	Model                  string  `json:"model"`
	Temperature            float64 `json:"temperature"`
	AutoExecute            bool    `json:"auto_execute"`
	DangerousAlertShown    bool    `json:"dangerous_alert_shown"`
	MatchedDangerousWord   string  `json:"matched_dangerous_word"`
	DangerousCheckBypassed bool    `json:"dangerous_check_bypassed"`
}

func GetLogFilePath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		home, hErr := os.UserHomeDir()
		if hErr != nil {
			return "", err
		}
		return filepath.Join(home, ".cache", "shelloma", "shelloma.log"), nil
	}
	return filepath.Join(cacheDir, "shelloma", "shelloma.log"), nil
}

func WriteLogEntry(entry LogEntry) error {
	path, err := GetLogFilePath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if entry.Timestamp == "" {
		entry.Timestamp = time.Now().Format(time.RFC3339)
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(append(data, '\n')); err != nil {
		return err
	}

	return nil
}
