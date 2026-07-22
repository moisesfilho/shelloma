package ollama

import (
	"fmt"
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
	Rules      []string
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
		Rules:   cfg.Rules,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}

	models, err := client.ListModels()
	if err != nil {
		offlineMsg := client.Trans.OllamaOfflineError
		if offlineMsg == "" {
			offlineMsg = "Could not connect to Ollama service."
		}
		return nil, fmt.Errorf("%s (%s): %w", offlineMsg, client.BaseURL, err)
	}

	if len(models) == 0 {
		noModelsMsg := client.Trans.OllamaNoModelsFound
		if noModelsMsg == "" {
			noModelsMsg = "No models found in Ollama. Download a model with e.g. `ollama pull qwen2.5-coder:1.5b`"
		}
		return nil, fmt.Errorf("%s", noModelsMsg)
	}

	if client.Model == "" {
		client.Model = selectBestModel(models)
	}

	return client, nil
}

func (c *Client) GetModel() string {
	return c.Model
}
