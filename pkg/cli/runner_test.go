package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"shelloma/pkg/config"
	"shelloma/pkg/i18n"
	"shelloma/pkg/logger"
	"shelloma/pkg/ollama"
	"shelloma/pkg/sysinfo"
	"shelloma/pkg/ui"
)

// Mock para LLMProvider usado no Runner
type runnerMockLLM struct {
	model             string
	explainResp       string
	fixCmd            string
	altCmd            string
	refinedCmd        string
	analysisSuccess   bool
	analysisReason    string
	analysisSuggest   string
}

func (m *runnerMockLLM) ListModels() ([]string, error) {
	return []string{m.model}, nil
}
func (m *runnerMockLLM) GenerateCommand(_ sysinfo.SystemContext, _ string, _ float64) (string, error) {
	return "ls -la", nil
}
func (m *runnerMockLLM) ExplainCommand(_ string) (string, error) {
	return m.explainResp, nil
}
func (m *runnerMockLLM) AnalyzeExecutionResult(_ string, _ int, _ string, _ sysinfo.SystemContext) (ollama.AnalysisResult, error) {
	return ollama.AnalysisResult{
		Success:          m.analysisSuccess,
		Reason:           m.analysisReason,
		SuggestedCommand: m.analysisSuggest,
	}, nil
}
func (m *runnerMockLLM) GenerateFixCommand(_ sysinfo.SystemContext, _ string, _ string) (string, error) {
	return m.fixCmd, nil
}
func (m *runnerMockLLM) GenerateRefinedCommand(_ sysinfo.SystemContext, _ string, _ string, _ string, _ float64) (string, error) {
	return m.refinedCmd, nil
}
func (m *runnerMockLLM) GenerateAlternativeCommand(_ sysinfo.SystemContext, _ string, _ string, _ string) (string, error) {
	return m.altCmd, nil
}
func (m *runnerMockLLM) GetModel() string {
	return m.model
}



func TestLogExecution(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "shelloma_log_exec_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	t.Setenv("XDG_CACHE_HOME", tempDir)
	t.Setenv("HOME", tempDir)
	t.Setenv("LocalAppData", tempDir)

	sysCtx := sysinfo.SystemContext{
		WorkingDir: "/tmp",
		User:       "test-user",
		OS:         "linux",
	}

	cfg := config.Config{
		OllamaURL:             "http://localhost:11434",
		DisableDangerousCheck: false,
		DangerousCommands:     []string{"rm"},
	}

	mockLLM := &runnerMockLLM{model: "qwen2.5-coder:1.5b"}

	LogExecution("listar", "ls", "Execute", 0, "file.txt", sysCtx, cfg, mockLLM)

	logPath, err := logger.GetLogFilePath()
	if err != nil {
		t.Fatalf("falhou ao obter caminho do log: %v", err)
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("falhou ao ler arquivo de log: %v", err)
	}

	var entry logger.LogEntry
	err = json.Unmarshal(data, &entry)
	if err != nil {
		t.Fatalf("falhou ao deserializar entrada de log: %v", err)
	}

	if entry.UserQuery != "listar" || entry.SuggestedCommand != "ls" || entry.UserAction != "Execute" {
		t.Errorf("Campos de log inconsistentes: %+v", entry)
	}
}

type mockLineReader struct {
	lines []string
	index int
}

func (r *mockLineReader) Read(p []byte) (n int, err error) {
	if r.index >= len(r.lines) {
		return 0, io.EOF
	}
	line := r.lines[r.index]
	r.index++
	copy(p, line)
	return len(line), io.EOF
}

func TestHandleUserAction(t *testing.T) {
	trans := i18n.GetTranslations("en")
	cfg := config.Config{}
	mockLLM := &runnerMockLLM{
		explainResp: "Explain response content",
	}

	// 1. Action Explain
	{
		oldStdin := ui.StdinReader
		ui.StdinReader = &mockLineReader{lines: []string{"e\n"}} // "e" -> ActionExplain
		defer func() { ui.StdinReader = oldStdin }()

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		cmd := "ls"
		action := HandleUserAction(mockLLM, sysinfo.SystemContext{}, &cmd, cfg, trans)

		w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		os.Stdout = oldStdout
		out := buf.String()

		if action != ui.ActionExplain {
			t.Errorf("Esperava ação Explain, obteve %d", action)
		}
		if !strings.Contains(out, "Explain response content") {
			t.Errorf("Esperava explicação na saída, obteve %q", out)
		}
	}

	// 2. Action Edit
	{
		oldStdin := ui.StdinReader
		ui.StdinReader = &mockLineReader{lines: []string{"m\n", "new_command\n"}} // "m" -> ActionEdit, seguido de "new_command"
		defer func() { ui.StdinReader = oldStdin }()

		cmd := "old_command"
		action := HandleUserAction(mockLLM, sysinfo.SystemContext{}, &cmd, cfg, trans)
		if action != ui.ActionEdit {
			t.Errorf("Esperava ação Edit, obteve %d", action)
		}
		if cmd != "new_command" {
			t.Errorf("Esperava comando editado para 'new_command', obteve %q", cmd)
		}
	}

	// 3. Action Quit
	{
		oldStdin := ui.StdinReader
		ui.StdinReader = &mockLineReader{lines: []string{"q\n"}} // "q" -> ActionQuit
		defer func() { ui.StdinReader = oldStdin }()

		cmd := "ls"
		action := HandleUserAction(mockLLM, sysinfo.SystemContext{}, &cmd, cfg, trans)
		if action != ui.ActionQuit {
			t.Errorf("Esperava ação Quit, obteve %d", action)
		}
	}
}



func TestConnectOrRecoverOllama(t *testing.T) {
	trans := i18n.GetTranslations("en")

	// Mock HTTP Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/tags" {
			resp := ollama.TagsResponse{
				Models: []ollama.ModelInfo{
					{Name: "qwen2.5-coder:1.5b"},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	cfg := config.Config{
		OllamaURL: server.URL,
		Model:     "qwen2.5-coder:1.5b",
	}

	client := ConnectOrRecoverOllama(cfg, trans)
	if client == nil {
		t.Fatalf("Esperava que o cliente fosse conectado")
	}
	if client.GetModel() != "qwen2.5-coder:1.5b" {
		t.Errorf("Modelo do cliente inesperado: %s", client.GetModel())
	}
}
