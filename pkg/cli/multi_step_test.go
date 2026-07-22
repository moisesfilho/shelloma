package cli

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"shelloma/pkg/config"
	"shelloma/pkg/i18n"
	"shelloma/pkg/logger"
	"shelloma/pkg/sysinfo"
	"shelloma/pkg/ui"
)

func TestSplitCommandSteps(t *testing.T) {
	input := "mkdir project\n\ncd project\n  \ngit init\n"
	expected := []string{"mkdir project", "cd project", "git init"}

	got := SplitCommandSteps(input)
	if len(got) != len(expected) {
		t.Fatalf("expected %d steps, got %d", len(expected), len(got))
	}

	for i, step := range got {
		if step != expected[i] {
			t.Errorf("step %d: expected %q, got %q", i, expected[i], step)
		}
	}
}

func TestHandleCdCommand(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "shelloma_cd_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working dir: %v", err)
	}
	defer func() { _ = os.Chdir(originalWd) }()

	sysCtx := &sysinfo.SystemContext{
		WorkingDir: originalWd,
	}

	err = handleCdCommand("cd "+tempDir, sysCtx)
	if err != nil {
		t.Fatalf("handleCdCommand failed: %v", err)
	}

	newWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get new working dir: %v", err)
	}

	// Paths might have symlink resolution, so evaluate real paths
	realTemp, _ := filepath.EvalSymlinks(tempDir)
	realNewWd, _ := filepath.EvalSymlinks(newWd)
	if realNewWd != realTemp {
		t.Errorf("expected process to change dir to %q, but got %q", realTemp, realNewWd)
	}

	if sysCtx.WorkingDir != newWd {
		t.Errorf("expected sysCtx.WorkingDir to be updated to %q, but got %q", newWd, sysCtx.WorkingDir)
	}
}

func TestExecuteMultiStep(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "shelloma_multistep_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Overwrite XDG_CACHE_HOME so tests write logs inside the temp directory
	t.Setenv("XDG_CACHE_HOME", tempDir)
	t.Setenv("HOME", tempDir)
	t.Setenv("LocalAppData", tempDir)

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working dir: %v", err)
	}
	defer func() { _ = os.Chdir(originalWd) }()

	sysCtx := &sysinfo.SystemContext{
		WorkingDir: originalWd,
		OS:         "linux",
		User:       "test-user",
	}

	cfg := config.Config{
		Language:              "en",
		DangerousCommands:     []string{"rm", "dd"},
		DisableDangerousCheck: false,
		Temperature:           0.2,
		OllamaURL:             "http://localhost:11434",
	}
	tTrans := i18n.GetTranslations("en")
	client := &mockLLM{suggestedCmd: "echo 'success'"}

	// Simulate user pressing Enter (empty input defaults to "y" inside choice parsing)
	oldStdin := ui.StdinReader
	ui.StdinReader = strings.NewReader("\n\n\n")
	defer func() { ui.StdinReader = oldStdin }()

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	multiStepCmd := "cd " + tempDir + "\nmkdir test_sub_dir\ncd test_sub_dir\n"
	success, lastEc, lastOut := ExecuteMultiStep(client, sysCtx, multiStepCmd, cfg, tTrans, "run multi steps")

	w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	os.Stdout = oldStdout
	output := buf.String()

	if !success {
		t.Fatalf("ExecuteMultiStep failed: last exit code: %d, output: %s", lastEc, lastOut)
	}

	// Verify terminal output contains prompt details
	if !strings.Contains(output, "Executing multi-step command") {
		t.Errorf("expected output to contain multi-step execution message, got %q", output)
	}
	if !strings.Contains(output, "[1/3]") || !strings.Contains(output, "[2/3]") || !strings.Contains(output, "[3/3]") {
		t.Errorf("expected output to show step progression (1/3, 2/3, 3/3), got %q", output)
	}

	// Verify directory traversal succeeded
	newWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get process working dir: %v", err)
	}
	expectedSubDir := filepath.Join(tempDir, "test_sub_dir")
	realExpected, _ := filepath.EvalSymlinks(expectedSubDir)
	realNewWd, _ := filepath.EvalSymlinks(newWd)
	if realNewWd != realExpected {
		t.Errorf("expected process working dir to be %q, but got %q", realExpected, realNewWd)
	}

	// Verify that logs were written for each step
	logPath, err := logger.GetLogFilePath()
	if err != nil {
		t.Fatalf("failed to get log file path: %v", err)
	}

	logData, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(logData)), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 log lines written (one for each executed step), but got %d", len(lines))
	}
}
