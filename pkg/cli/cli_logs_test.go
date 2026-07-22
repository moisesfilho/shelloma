package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"shelloma/pkg/config"
	"shelloma/pkg/i18n"
	"shelloma/pkg/logger"
	"shelloma/pkg/ui"
)

func TestHandleLogsCommandTerminalOption(t *testing.T) {
	// Set up temporary log directory
	tempDir, err := os.MkdirTemp("", "shelloma_cli_log_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	t.Setenv("XDG_CACHE_HOME", tempDir)
	t.Setenv("HOME", tempDir)
	t.Setenv("LocalAppData", tempDir)

	tTrans := i18n.GetTranslations("en")
	cfg := config.Config{Language: "en"}

	// 1. Check when no logs exist
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	HandleLogsCommand(cfg, tTrans)

	w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	os.Stdout = oldStdout
	outNoLogs := buf.String()

	if !strings.Contains(outNoLogs, "No logs found") {
		t.Errorf("expected output to indicate no logs found, but got: %q", outNoLogs)
	}

	// 2. Write a dummy log entry
	entry := logger.LogEntry{
		Timestamp:        "2026-07-22T13:38:48Z",
		UserQuery:        "test query prompt",
		SuggestedCommand: "echo 'hello'",
		UserAction:       "Execute",
		ExitCode:         0,
		User:             "test-user",
	}
	err = logger.WriteLogEntry(entry)
	if err != nil {
		t.Fatalf("failed to write dummy log: %v", err)
	}

	// Override StdinReader to simulate option "1" (Show in terminal)
	oldStdin := ui.StdinReader
	ui.StdinReader = strings.NewReader("1\n")
	defer func() { ui.StdinReader = oldStdin }()

	// Capture stdout
	r2, w2, _ := os.Pipe()
	os.Stdout = w2

	HandleLogsCommand(cfg, tTrans)

	w2.Close()
	var buf2 bytes.Buffer
	_, _ = io.Copy(&buf2, r2)
	os.Stdout = oldStdout
	outWithLogs := buf2.String()

	if !strings.Contains(outWithLogs, "test query prompt") {
		t.Errorf("expected output to contain log entry prompt, but got: %q", outWithLogs)
	}
	if !strings.Contains(outWithLogs, "echo 'hello'") {
		t.Errorf("expected output to contain log entry command, but got: %q", outWithLogs)
	}
}
