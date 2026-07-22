package logger

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestWriteLogEntryAndRead(t *testing.T) {
	// Create a temporary log file path by overriding GetLogFilePath behavior
	tempDir, err := os.MkdirTemp("", "shelloma_log_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Temporarily override UserCacheDir env variables to redirect logs to tempDir
	t.Setenv("XDG_CACHE_HOME", tempDir)
	t.Setenv("HOME", tempDir)
	t.Setenv("LocalAppData", tempDir)

	path, err := GetLogFilePath()
	if err != nil {
		t.Fatalf("failed to get log file path: %v", err)
	}

	if !strings.Contains(path, tempDir) {
		t.Errorf("expected log file path to be under tempDir %q, but got %q", tempDir, path)
	}

	entry := LogEntry{
		Timestamp:            "2026-07-22T13:38:48Z",
		UserQuery:            "list files",
		SuggestedCommand:     "ls -la",
		UserAction:           "Execute",
		ExitCode:             0,
		CommandOutput:        "file1.txt",
		WorkingDir:           "/tmp",
		User:                 "test-user",
		OS:                   "linux",
		OllamaURL:            "http://localhost:11434",
		Model:                "test-model",
		Temperature:          0.1,
		AutoExecute:          false,
		DangerousAlertShown:  false,
		MatchedDangerousWord: "",
	}

	err = WriteLogEntry(entry)
	if err != nil {
		t.Fatalf("failed to write log entry: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("log file was not created")
	}

	// Read and verify entry
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	var readEntry LogEntry
	err = json.Unmarshal([]byte(strings.TrimSpace(string(data))), &readEntry)
	if err != nil {
		t.Fatalf("failed to unmarshal log entry: %v", err)
	}

	if readEntry.UserQuery != entry.UserQuery {
		t.Errorf("expected UserQuery %q, got %q", entry.UserQuery, readEntry.UserQuery)
	}
	if readEntry.SuggestedCommand != entry.SuggestedCommand {
		t.Errorf("expected SuggestedCommand %q, got %q", entry.SuggestedCommand, readEntry.SuggestedCommand)
	}
	if readEntry.UserAction != entry.UserAction {
		t.Errorf("expected UserAction %q, got %q", entry.UserAction, readEntry.UserAction)
	}
}
