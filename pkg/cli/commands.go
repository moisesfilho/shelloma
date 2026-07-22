package cli

import (
	"fmt"
	"os"
	"strings"

	"shelloma/pkg/config"
	"shelloma/pkg/i18n"
	"shelloma/pkg/ui"
)

func HandleModelsCommand(cfg config.Config, t i18n.Translations) {
	client := ConnectOrRecoverOllama(cfg, t)

	models, err := client.ListModels()
	if err != nil {
		fmt.Printf("%s%s %v%s\n", ui.Red, t.ErrorPrefix, err, ui.Reset)
		os.Exit(1)
	}

	fmt.Printf("%s%s (%s):%s\n", ui.Bold+ui.Cyan, t.ModelsInstalled, cfg.OllamaURL, ui.Reset)
	for _, m := range models {
		active := ""
		if m == client.GetModel() {
			active = fmt.Sprintf(" %s%s%s", ui.Green, t.ModelActive, ui.Reset)
		}
		fmt.Printf(" - %s%s\n", m, active)
	}
}

func HandleConfigCommand(cfg config.Config, args []string, t i18n.Translations) {
	if len(args) >= 3 && args[0] == "set" {
		key := strings.ToLower(args[1])
		val := args[2]

		switch key {
		case "model":
			cfg.Model = val
		case "ollama_url", "url":
			cfg.OllamaURL = val
		case "lang", "language":
			cfg.Language = string(i18n.NormalizeLanguage(val))
			t = i18n.GetTranslations(cfg.Language)
		case "dangerous_commands", "dangerous":
			commands := strings.Split(val, ",")
			for i, cmd := range commands {
				commands[i] = strings.TrimSpace(cmd)
			}
			cfg.DangerousCommands = commands
		case "disable_dangerous_check":
			cfg.DisableDangerousCheck = val == "true" || val == "1" || val == "yes"
		default:
			fmt.Printf("%s%s: %s%s\n", ui.Red, t.UnknownKey, key, ui.Reset)
			os.Exit(1)
		}

		if err := config.SaveConfig(cfg); err != nil {
			fmt.Printf("%s%s %v%s\n", ui.Red, t.ErrorPrefix, err, ui.Reset)
			os.Exit(1)
		}

		fmt.Printf("%s✔ %s (%s = %s)%s\n", ui.Green, t.ConfigSaved, key, val, ui.Reset)
		return
	}

	path, _ := config.GetConfigPath()
	fmt.Printf("%s%s (%s):%s\n", ui.Bold+ui.Cyan, t.CurrentConfig, path, ui.Reset)
	fmt.Printf("  ollama_url:              %s\n", cfg.OllamaURL)
	fmt.Printf("  model:                   %s %s\n", cfg.Model, t.DefaultModelAuto)
	fmt.Printf("  language:                %s (%s)\n", cfg.Language, i18n.GetTranslations(cfg.Language).LanguageName)
	fmt.Printf("  temperature:             %.2f\n", cfg.Temperature)
	fmt.Printf("  auto_execute:            %t\n", cfg.AutoExecute)
	fmt.Printf("  dangerous_commands:      %s\n", strings.Join(cfg.DangerousCommands, ", "))
	fmt.Printf("  disable_dangerous_check: %t\n", cfg.DisableDangerousCheck)
}
