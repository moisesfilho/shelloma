package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"shelloma/pkg/config"
	"shelloma/pkg/i18n"
	"shelloma/pkg/ollama"
	"shelloma/pkg/sysinfo"
	"shelloma/pkg/ui"
)

// mockLLM is a helper for mocking LLMProvider in runner tests
type mockLLM struct {
	ollama.LLMProvider
	suggestedCmd string
}

func (m *mockLLM) GenerateCommand(_ sysinfo.SystemContext, _ string, _ float64) (string, error) {
	return m.suggestedCmd, nil
}

func (m *mockLLM) GetModel() string {
	return "mock-model"
}

func (m *mockLLM) AnalyzeExecutionResult(_ string, exitCode int, _ string, _ sysinfo.SystemContext) (ollama.AnalysisResult, error) {
	return ollama.AnalysisResult{
		Success: exitCode == 0,
		Reason:  "Execution completed",
	}, nil
}

func TestDangerousCommandFlow(t *testing.T) {
	// Save and restore StdinReader
	oldStdinReader := ui.StdinReader
	defer func() {
		ui.StdinReader = oldStdinReader
	}()

	// 1. Create a temp file
	tmpFile, err := os.CreateTemp("", "shelloma-danger-test-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	tmpFile.Close()

	// Verify temp file exists
	if _, err := os.Stat(tmpFile.Name()); os.IsNotExist(err) {
		t.Fatalf("temp file was not created")
	}

	// 2. Prepare configuration with "rm" as a dangerous command
	cfg := config.Config{
		OllamaURL:             "http://localhost:11434",
		Model:                 "mock-model",
		Temperature:           0.1,
		AutoExecute:           true,
		Language:              "en",
		DangerousCommands:     []string{"rm"},
		DisableDangerousCheck: false,
	}

	tTrans := i18n.GetTranslations("en")

	cmdToDelete := "rm " + tmpFile.Name()
	client := &mockLLM{suggestedCmd: cmdToDelete}
	sysCtx := sysinfo.SystemContext{OS: "linux"}

	// Scenario A: Incorrect security word
	// Simulate entering an incorrect safety word
	ui.StdinReader = strings.NewReader("WRONG\n")

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	success := ExecuteWithRecovery(client, sysCtx, cmdToDelete, cfg, tTrans)

	w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	os.Stdout = oldStdout
	outputA := buf.String()

	if success {
		t.Error("Expected execution to fail due to incorrect security word, but it succeeded")
	}

	// Verify file was NOT deleted
	if _, err := os.Stat(tmpFile.Name()); os.IsNotExist(err) {
		t.Error("File was deleted even though incorrect safety word was entered")
	}

	// Validate warning and abort messages were displayed
	if !strings.Contains(outputA, "Incorrect security word. Execution aborted.") {
		t.Errorf("Expected output to contain incorrect security word abort message, but got: %q", outputA)
	}

	// Scenario B: Correct security word
	ui.StdinReader = strings.NewReader("CONFIRM\n")

	// Capture stdout again
	r2, w2, _ := os.Pipe()
	os.Stdout = w2

	success = ExecuteWithRecovery(client, sysCtx, cmdToDelete, cfg, tTrans)

	w2.Close()
	var buf2 bytes.Buffer
	_, _ = io.Copy(&buf2, r2)
	os.Stdout = oldStdout
	outputB := buf2.String()

	if !success {
		t.Errorf("Expected execution to succeed with correct safety word, but it failed. Output: %s", outputB)
	}

	// Verify file WAS deleted
	if _, err := os.Stat(tmpFile.Name()); !os.IsNotExist(err) {
		t.Error("File was not deleted even though correct safety word was entered")
	}

	// Scenario C: Disable dangerous check configuration
	// Create another temp file to delete
	tmpFileC, err := os.CreateTemp("", "shelloma-danger-test-c-*")
	if err != nil {
		t.Fatalf("failed to create temp file for scenario C: %v", err)
	}
	defer func() { _ = os.Remove(tmpFileC.Name()) }()
	tmpFileC.Close()

	cfg.DisableDangerousCheck = true
	cmdToDeleteC := "rm " + tmpFileC.Name()

	// Feed dummy input that should be ignored because no prompt should happen
	ui.StdinReader = strings.NewReader("SHOULD_NOT_BE_READ\n")

	// Capture stdout again
	r3, w3, _ := os.Pipe()
	os.Stdout = w3

	successC := ExecuteWithRecovery(client, sysCtx, cmdToDeleteC, cfg, tTrans)

	w3.Close()
	var buf3 bytes.Buffer
	_, _ = io.Copy(&buf3, r3)
	os.Stdout = oldStdout
	outputC := buf3.String()

	if !successC {
		t.Errorf("Expected execution to succeed directly with DisableDangerousCheck: true, but it failed. Output: %s", outputC)
	}

	// Verify file WAS deleted
	if _, err := os.Stat(tmpFileC.Name()); !os.IsNotExist(err) {
		t.Error("File in scenario C was not deleted")
	}

	// Verify no security prompt was printed
	if strings.Contains(outputC, "security word") || strings.Contains(outputC, "palavra de segurança") {
		t.Error("Expected no safety confirmation prompt, but found one in output")
	}
}

func TestDangerousWindowsCommandFlow(t *testing.T) {
	oldStdinReader := ui.StdinReader
	defer func() {
		ui.StdinReader = oldStdinReader
	}()

	cfg := config.Config{
		OllamaURL:             "http://localhost:11434",
		Model:                 "mock-model",
		Temperature:           0.1,
		AutoExecute:           true,
		Language:              "en",
		DangerousCommands:     []string{"Remove-Item", "del", "Format-Volume"},
		DisableDangerousCheck: false,
	}

	tTrans := i18n.GetTranslations("en")
	client := &mockLLM{suggestedCmd: "Remove-Item -Path C:\\test -Recurse"}
	sysCtx := sysinfo.SystemContext{OS: "windows"}

	// Simulate incorrect security word
	ui.StdinReader = strings.NewReader("WRONG\n")

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	success := ExecuteWithRecovery(client, sysCtx, "Remove-Item -Path C:\\test -Recurse", cfg, tTrans)

	w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	os.Stdout = oldStdout
	output := buf.String()

	if success {
		t.Error("Expected execution to fail due to incorrect security word, but it succeeded")
	}

	if !strings.Contains(output, "Incorrect security word. Execution aborted.") {
		t.Errorf("Expected output to contain incorrect security word abort message, but got: %q", output)
	}
}
