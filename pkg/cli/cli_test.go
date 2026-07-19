package cli

import (
	"testing"

	"shelloma/pkg/config"
	"shelloma/pkg/i18n"
)

func TestApplyFlagOverrides(t *testing.T) {
	cfg := config.Config{
		Model:       "default-model",
		OllamaURL:   "http://localhost:11434",
		Language:    "en",
		AutoExecute: false,
	}

	trans := i18n.GetTranslations("en")

	ApplyFlagOverrides(&cfg, "custom-model", "http://custom:11434", "pt", true, &trans)

	if cfg.Model != "custom-model" {
		t.Errorf("Esperava custom-model, obteve %s", cfg.Model)
	}

	if cfg.OllamaURL != "http://custom:11434" {
		t.Errorf("Esperava http://custom:11434, obteve %s", cfg.OllamaURL)
	}

	if cfg.Language != "pt" {
		t.Errorf("Esperava pt, obteve %s", cfg.Language)
	}

	if !cfg.AutoExecute {
		t.Errorf("Esperava AutoExecute true, obteve false")
	}

	if trans.LanguageName != "Português Brasileiro" {
		t.Errorf("Esperava Português Brasileiro, obteve %s", trans.LanguageName)
	}
}
