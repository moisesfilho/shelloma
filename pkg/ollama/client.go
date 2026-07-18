package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"shelloma/pkg/config"
	"shelloma/pkg/i18n"
	"shelloma/pkg/sysinfo"
)

// LLMProvider define o contrato abstrato para provedores de modelos de linguagem (DIP - SOLID)
type LLMProvider interface {
	ListModels() ([]string, error)
	GenerateCommand(sysCtx sysinfo.SystemContext, userPrompt string, temp float64) (string, error)
	ExplainCommand(command string) (string, error)
	AnalyzeExecutionResult(cmdStr string, exitCode int, output string, sysCtx sysinfo.SystemContext) (AnalysisResult, error)
	GenerateFixCommand(sysCtx sysinfo.SystemContext, failedCmd string, errorOutput string) (string, error)
	GetModel() string
}

type Client struct {
	BaseURL    string
	Model      string
	HTTPClient *http.Client
	Trans      i18n.Translations
}

type ModelInfo struct {
	Name string `json:"name"`
}

type TagsResponse struct {
	Models []ModelInfo `json:"models"`
}

type GenerateRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	System  string                 `json:"system,omitempty"`
	Format  string                 `json:"format,omitempty"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type GenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type AnalysisResult struct {
	Success          bool   `json:"success"`
	Reason           string `json:"reason"`
	SuggestedCommand string `json:"suggested_command"`
}

func NewClient(cfg config.Config) (LLMProvider, error) {
	client := &Client{
		BaseURL: strings.TrimRight(cfg.OllamaURL, "/"),
		Model:   cfg.Model,
		Trans:   i18n.GetTranslations(cfg.Language),
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}

	models, err := client.ListModels()
	if err != nil {
		return nil, fmt.Errorf("could not connect to Ollama at %s: %w\nMake sure Ollama is running (`ollama serve`)", client.BaseURL, err)
	}

	if len(models) == 0 {
		return nil, fmt.Errorf("no models found in Ollama. Download a model with e.g. `ollama pull qwen2.5-coder:1.5b`")
	}

	if client.Model == "" {
		client.Model = selectBestModel(models)
	}

	return client, nil
}

func (c *Client) GetModel() string {
	return c.Model
}

func (c *Client) ListModels() ([]string, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/api/tags")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama status %d", resp.StatusCode)
	}

	var tags TagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return nil, err
	}

	var names []string
	for _, m := range tags.Models {
		names = append(names, m.Name)
	}
	return names, nil
}

func selectBestModel(models []string) string {
	preferences := []string{
		"qwen2.5-coder",
		"deepseek-coder",
		"llama3.2",
		"deepseek-r1",
		"codellama",
		"mistral",
	}

	for _, pref := range preferences {
		for _, m := range models {
			if strings.Contains(strings.ToLower(m), pref) {
				return m
			}
		}
	}

	return models[0]
}

func (c *Client) GenerateCommand(sysCtx sysinfo.SystemContext, userPrompt string, temp float64) (string, error) {
	systemPrompt := fmt.Sprintf(`You are an expert Linux Shell CLI assistant. Your sole job is to translate the user request into a precise, valid, executable shell command for their system.

User System Info:
- OS: %s (%s %s)
- Shell: %s
- Working Directory ($PWD): %s
- Current User ($USER): %s (IsRoot: %t)
- Architecture: %s

STRICT RESPONSE RULES:
1. Respond ONLY with the raw executable shell command.
2. DO NOT use markdown code block wrappers like `+"```bash"+` or `+"```"+`.
3. DO NOT add introductions, explanations, greetings, or comments.
4. If multiple commands are needed, chain them with && or ;.
5. Adapt commands specifically for Linux distro (%s).`,
		sysCtx.OS, sysCtx.DistroName, sysCtx.DistroVer,
		sysCtx.Shell,
		sysCtx.WorkingDir,
		sysCtx.User, sysCtx.IsRoot,
		sysCtx.Arch,
		sysCtx.DistroName,
	)

	reqBody := GenerateRequest{
		Model:  c.Model,
		Prompt: userPrompt,
		System: systemPrompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": temp,
		},
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/generate", "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", fmt.Errorf("error calling Ollama API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var genResp GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return "", err
	}

	cmd := cleanCommandOutput(genResp.Response)
	return cmd, nil
}

func (c *Client) ExplainCommand(command string) (string, error) {
	prompt := fmt.Sprintf("%s\n\n%s", c.Trans.OllamaLangInstruction, command)

	reqBody := GenerateRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.2,
		},
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/generate", "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var genResp GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return "", err
	}

	return strings.TrimSpace(genResp.Response), nil
}

func (c *Client) AnalyzeExecutionResult(cmdStr string, exitCode int, output string, sysCtx sysinfo.SystemContext) (AnalysisResult, error) {
	if exitCode == 0 && !containsErrorKeywords(output) {
		return AnalysisResult{
			Success: true,
			Reason:  c.Trans.Success,
		}, nil
	}

	systemPrompt := fmt.Sprintf(`You are a Linux Shell execution output analyzer.
Analyze the executed command, exit code, and terminal return output.
Determine if the command achieved SUCCESS or if it FAILED.

%s

Respond STRICTLY in JSON format matching this schema:
{
  "success": true or false,
  "reason": "Concise explanation of success or failure reason in the target language",
  "suggested_command": "Valid Linux executable shell command to fix/investigate the issue (e.g. touch file.txt, mkdir -p dir, ls -la path, sudo apt update)"
}

STRICT RULES:
1. If exitCode != 0 or output contains explicit error messages (e.g. 'No such file or directory', 'Permission denied', 'command not found'), mark success: false.
2. When success is false, you MUST provide a valid executable shell command in 'suggested_command'.
3. The 'suggested_command' field MUST contain ONLY a raw executable shell command (NO explanations, NO quotes, NO code fences).
4. %s
5. Do NOT include any text before or after the JSON.`, c.Trans.OllamaLangInstruction, c.Trans.OllamaReasonInstruction)

	prompt := fmt.Sprintf("Executed Command: %s\nExit Code: %d\nTerminal Output:\n%s", cmdStr, exitCode, output)

	reqBody := GenerateRequest{
		Model:  c.Model,
		Prompt: prompt,
		System: systemPrompt,
		Format: "json",
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.1,
		},
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return AnalysisResult{Success: exitCode == 0}, err
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/generate", "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return AnalysisResult{
			Success: exitCode == 0,
			Reason:  fmt.Sprintf("Execution error (exit code %d)", exitCode),
		}, nil
	}
	defer resp.Body.Close()

	var genResp GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return AnalysisResult{Success: exitCode == 0}, nil
	}

	jsonStr := cleanJSONOutput(genResp.Response)
	var result AnalysisResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return AnalysisResult{
			Success: exitCode == 0,
			Reason:  fmt.Sprintf("Finished with exit code %d", exitCode),
		}, nil
	}

	result.SuggestedCommand = cleanCommandOutput(result.SuggestedCommand)
	if !isValidShellCommand(result.SuggestedCommand) {
		result.SuggestedCommand = ""
	}

	return result, nil
}

func containsErrorKeywords(output string) bool {
	lowered := strings.ToLower(output)
	errorKeywords := []string{
		"no such file or directory",
		"permission denied",
		"command not found",
		"syntax error",
		"fatal error",
		"failed to",
		"cannot access",
		"error:",
	}
	for _, kw := range errorKeywords {
		if strings.Contains(lowered, kw) {
			return true
		}
	}
	return false
}

func (c *Client) GenerateFixCommand(sysCtx sysinfo.SystemContext, failedCmd string, errorOutput string) (string, error) {
	systemPrompt := fmt.Sprintf(`You are a Linux Shell expert assistant.
A command failed to execute on (%s %s).
Your goal is to generate a NEW valid executable shell command to fix the issue, create missing files/folders, or investigate the cause.

RULES:
1. Respond ONLY with the raw executable shell command.
2. DO NOT use markdown code block formatting.
3. DO NOT add explanatory text or sentences.`, sysCtx.DistroName, sysCtx.DistroVer)

	prompt := fmt.Sprintf("Failed Command: %s\nError Output:\n%s", failedCmd, errorOutput)

	reqBody := GenerateRequest{
		Model:  c.Model,
		Prompt: prompt,
		System: systemPrompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.1,
		},
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/generate", "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var genResp GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return "", err
	}

	cmd := cleanCommandOutput(genResp.Response)
	if !isValidShellCommand(cmd) {
		return "", nil
	}
	return cmd, nil
}

func isValidShellCommand(cmd string) bool {
	if cmd == "" {
		return false
	}
	lowered := strings.ToLower(cmd)
	badWords := []string{"você", "pode", "tente", "criar", "verifique", "arquivo", "diretório", "caso", "este", "se ", "you ", "can ", "try "}
	for _, bw := range badWords {
		if strings.Contains(lowered, bw+" ") || strings.HasSuffix(lowered, bw) {
			return false
		}
	}
	return true
}

func cleanJSONOutput(output string) string {
	cleaned := strings.TrimSpace(output)
	if strings.HasPrefix(cleaned, "```json") {
		cleaned = strings.TrimPrefix(cleaned, "```json")
		cleaned = strings.TrimSuffix(cleaned, "```")
	} else if strings.HasPrefix(cleaned, "```") {
		cleaned = strings.TrimPrefix(cleaned, "```")
		cleaned = strings.TrimSuffix(cleaned, "```")
	}
	return strings.TrimSpace(cleaned)
}

func cleanCommandOutput(output string) string {
	cleaned := strings.TrimSpace(output)

	if strings.HasPrefix(cleaned, "```") {
		lines := strings.Split(cleaned, "\n")
		if len(lines) > 1 {
			if strings.HasPrefix(lines[0], "```") {
				lines = lines[1:]
			}
			if len(lines) > 0 && strings.HasPrefix(lines[len(lines)-1], "```") {
				lines = lines[:len(lines)-1]
			}
			cleaned = strings.Join(lines, "\n")
		}
	}

	cleaned = strings.TrimPrefix(cleaned, "`")
	cleaned = strings.TrimSuffix(cleaned, "`")
	return strings.TrimSpace(cleaned)
}
