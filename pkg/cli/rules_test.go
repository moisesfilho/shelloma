package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"shelloma/pkg/config"
	"shelloma/pkg/i18n"
	"shelloma/pkg/ui"
)

func TestRulesCommandFlow(t *testing.T) {
	// Set up temporary configuration path envs
	tempDir, err := os.MkdirTemp("", "shelloma_rules_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	t.Setenv("XDG_CACHE_HOME", tempDir)
	t.Setenv("HOME", tempDir)
	t.Setenv("LocalAppData", tempDir)

	tTrans := i18n.GetTranslations("en")
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("failed to load default config: %v", err)
	}

	// Stdin hijacking
	oldStdin := ui.StdinReader
	defer func() { ui.StdinReader = oldStdin }()

	// 1. List rules when empty
	{
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		HandleRulesCommand(cfg, []string{"list"}, tTrans)

		w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		os.Stdout = oldStdout
		out := buf.String()

		if !strings.Contains(out, "No custom rules defined") {
			t.Errorf("expected empty rules message, got: %q", out)
		}
	}

	// 2. Add rule via arguments
	{
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		HandleRulesCommand(cfg, []string{"add", "Always use nano for editing"}, tTrans)

		w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		os.Stdout = oldStdout
		out := buf.String()

		if !strings.Contains(out, "Rule added successfully") {
			t.Errorf("expected rule added message, got: %q", out)
		}

		// Reload config and verify
		cfg, _ = config.LoadConfig()
		if len(cfg.Rules) != 1 || cfg.Rules[0] != "Always use nano for editing" {
			t.Errorf("expected 1 rule in config, got %d: %v", len(cfg.Rules), cfg.Rules)
		}
	}

	// 3. Add rule via Stdin prompt
	{
		ui.StdinReader = strings.NewReader("Always log to custom paths\n")

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		HandleRulesCommand(cfg, []string{"add"}, tTrans)

		w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		os.Stdout = oldStdout
		out := buf.String()

		if !strings.Contains(out, "Rule added successfully") {
			t.Errorf("expected rule added message, got: %q", out)
		}

		// Reload config and verify
		cfg, _ = config.LoadConfig()
		if len(cfg.Rules) != 2 || cfg.Rules[1] != "Always log to custom paths" {
			t.Errorf("expected 2 rules in config, got %d: %v", len(cfg.Rules), cfg.Rules)
		}
	}

	// 4. Edit rule via arguments (Edit index 1: index 0 internally)
	{
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		HandleRulesCommand(cfg, []string{"edit", "1", "Always use vim for editing"}, tTrans)

		w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		os.Stdout = oldStdout
		out := buf.String()

		if !strings.Contains(out, "Rule updated successfully") {
			t.Errorf("expected rule updated message, got: %q", out)
		}

		// Reload config and verify
		cfg, _ = config.LoadConfig()
		if cfg.Rules[0] != "Always use vim for editing" {
			t.Errorf("expected rule to be updated to vim, got: %q", cfg.Rules[0])
		}
	}

	// 5. Edit rule via Stdin selection & input (Edit index 2)
	{
		ui.StdinReader = strings.NewReader("2\nAlways log to stdout\n")

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		HandleRulesCommand(cfg, []string{"edit"}, tTrans)

		w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		os.Stdout = oldStdout
		out := buf.String()

		if !strings.Contains(out, "Rule updated successfully") {
			t.Errorf("expected rule updated message, got: %q", out)
		}

		// Reload config and verify
		cfg, _ = config.LoadConfig()
		if cfg.Rules[1] != "Always log to stdout" {
			t.Errorf("expected rule 2 to be updated, got: %q", cfg.Rules[1])
		}
	}

	// 6. Delete rule via arguments (Delete index 1: index 0 internally)
	{
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		HandleRulesCommand(cfg, []string{"delete", "1"}, tTrans)

		w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		os.Stdout = oldStdout
		out := buf.String()

		if !strings.Contains(out, "Rule deleted successfully") {
			t.Errorf("expected rule deleted message, got: %q", out)
		}

		// Reload config and verify
		cfg, _ = config.LoadConfig()
		if len(cfg.Rules) != 1 || cfg.Rules[0] != "Always log to stdout" {
			t.Errorf("expected only 1 rule left ('Always log to stdout'), got: %v", cfg.Rules)
		}
	}

	// 7. Delete rule via Stdin selection
	{
		ui.StdinReader = strings.NewReader("1\n")

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		HandleRulesCommand(cfg, []string{"delete"}, tTrans)

		w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		os.Stdout = oldStdout
		out := buf.String()

		if !strings.Contains(out, "Rule deleted successfully") {
			t.Errorf("expected rule deleted message, got: %q", out)
		}

		// Reload config and verify
		cfg, _ = config.LoadConfig()
		if len(cfg.Rules) != 0 {
			t.Errorf("expected rules to be empty, got: %v", cfg.Rules)
		}
	}
}
