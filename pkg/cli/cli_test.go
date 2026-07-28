package cli

import (
	"os"
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

func TestParseLanguageOverride(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	cfg := config.Config{Language: "en"}
	os.Args = []string{"shelloma", "-l", "pt"}
	ParseLanguageOverride(&cfg)
	if cfg.Language != "pt" {
		t.Errorf("Esperava pt, obteve %s", cfg.Language)
	}

	cfg = config.Config{Language: "en"}
	os.Args = []string{"shelloma", "--lang=es", "unused"}
	ParseLanguageOverride(&cfg)
	if cfg.Language != "es" {
		t.Errorf("Esperava es, obteve %s", cfg.Language)
	}
}

func TestSetupFlags(t *testing.T) {
	t.Log("Testing SetupFlags")
	var modelFlag, urlFlag, langFlag string
	var yesFlag, verFlag, desktopFlag bool
	trans := i18n.GetTranslations("en")
	SetupFlags(&modelFlag, &urlFlag, &langFlag, &yesFlag, &verFlag, &desktopFlag, trans, "1.2.2")
}
