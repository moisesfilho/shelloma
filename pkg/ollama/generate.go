package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"shelloma/pkg/sysinfo"
)

func (c *Client) GenerateCommand(sysCtx sysinfo.SystemContext, userPrompt string, temp float64) (string, error) {
	osPersona := "Linux Shell CLI assistant"
	if sysCtx.OS == "windows" {
		osPersona = "Windows Command Line & PowerShell assistant"
	} else if sysCtx.OS == "darwin" {
		osPersona = "macOS Terminal & Shell assistant"
	}

	systemPrompt := fmt.Sprintf(`You are an expert %s. Your sole job is to translate the user request into a precise, valid, executable terminal command for their system.

User System Info:
- OS: %s (%s %s)
- Shell: %s
- Working Directory ($PWD): %s
- Current User ($USER): %s (IsRoot: %t)
- Architecture: %s

STRICT RESPONSE RULES:
1. Respond ONLY with the raw executable command.
2. DO NOT use markdown code block wrappers like `+"```"+`.
3. DO NOT add introductions, explanations, greetings, or comments.
4. If multiple commands are needed, write them as separate command lines (one command per line) so the tool can parse and execute them step-by-step with user confirmation.
5. Adapt commands specifically for system %s / shell %s.`,
		osPersona,
		sysCtx.OS, sysCtx.DistroName, sysCtx.DistroVer,
		sysCtx.Shell,
		sysCtx.WorkingDir,
		sysCtx.User, sysCtx.IsRoot,
		sysCtx.Arch,
		sysCtx.DistroName, sysCtx.Shell,
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

func (c *Client) GenerateFixCommand(sysCtx sysinfo.SystemContext, failedCmd string, errorOutput string) (string, error) {
	systemPrompt := fmt.Sprintf(`You are an expert CLI assistant.
A command failed to execute on (%s %s, shell %s).
Your goal is to generate a NEW valid executable command to fix the issue, create missing files/folders, or investigate the cause.

RULES:
1. Respond ONLY with the raw executable command.
2. DO NOT use markdown code block formatting.
3. DO NOT add explanatory text or sentences.`, sysCtx.DistroName, sysCtx.DistroVer, sysCtx.Shell)

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
