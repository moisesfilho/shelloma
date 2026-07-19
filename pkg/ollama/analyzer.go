package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"shelloma/pkg/sysinfo"
)

func (c *Client) AnalyzeExecutionResult(cmdStr string, exitCode int, output string, sysCtx sysinfo.SystemContext) (AnalysisResult, error) {
	if exitCode == 0 && !containsErrorKeywords(output) {
		return AnalysisResult{
			Success: true,
			Reason:  "Comando executado com sucesso",
		}, nil
	}

	systemPrompt := fmt.Sprintf(`You are an expert CLI execution analyzer.
You analyze terminal command outputs to determine if a command succeeded or failed.
%s

Input System Info:
- OS: %s (%s)
- Shell: %s

Required Output Format (STRICT JSON ONLY):
{
  "success": false,
  "reason": "Brief explanation of what failed or was missing (in the target language)",
  "suggested_command": "A valid shell command to fix the issue or create required resources, or empty string if no fix needed"
}`,
		c.Trans.OllamaReasonInstruction,
		sysCtx.OS, sysCtx.DistroName,
		sysCtx.Shell,
	)

	prompt := fmt.Sprintf("Command Executed: %s\nExit Code: %d\nOutput:\n%s", cmdStr, exitCode, output)

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

	if resp.StatusCode != http.StatusOK {
		return AnalysisResult{
			Success: exitCode == 0,
			Reason:  fmt.Sprintf("Execution error (exit code %d)", exitCode),
		}, nil
	}

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
		"permission denied",
		"permissão negada",
		"command not found",
		"comando não encontrado",
		"no such file or directory",
		"arquivo ou diretório não encontrado",
		"fatal:",
		"error:",
		"erro:",
		"failed to",
		"falha ao",
	}

	for _, kw := range errorKeywords {
		if strings.Contains(lowered, kw) {
			return true
		}
	}
	return false
}
