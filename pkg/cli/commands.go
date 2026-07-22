package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"shelloma/pkg/config"
	"shelloma/pkg/i18n"
	"shelloma/pkg/logger"
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

func HandleLogsCommand(_ config.Config, t i18n.Translations) {
	logPath, err := logger.GetLogFilePath()
	if err != nil {
		fmt.Printf("%s%s %v%s\n", ui.Red, t.ErrorPrefix, err, ui.Reset)
		os.Exit(1)
	}

	info, err := os.Stat(logPath)
	if os.IsNotExist(err) || info.Size() == 0 {
		fmt.Printf("%s%s%s\n", ui.Yellow, t.LogNoLogsFound, ui.Reset)
		return
	}

	fmt.Print(t.LogPromptOptions)
	choice := readLine()

	switch choice {
	case "1":
		data, err := os.ReadFile(logPath)
		if err != nil {
			fmt.Printf("%s%s %v%s\n", ui.Red, t.ErrorPrefix, err, ui.Reset)
			os.Exit(1)
		}
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			var entry logger.LogEntry
			if err := json.Unmarshal([]byte(line), &entry); err != nil {
				fmt.Println(line)
				continue
			}
			timeStr := entry.Timestamp
			if len(timeStr) > 19 {
				timeStr = timeStr[:19]
			}
			dangerStr := ""
			if entry.DangerousAlertShown {
				dangerStr = fmt.Sprintf(" %s[⚠️ DANGER: %s]%s", ui.Red, entry.MatchedDangerousWord, ui.Reset)
			}
			actionColor := ui.Reset
			switch entry.UserAction {
			case "Execute":
				actionColor = ui.Green
			case "Explain":
				actionColor = ui.Cyan
			case "Edit":
				actionColor = ui.Yellow
			case "Copy":
				actionColor = ui.Magenta
			case "Quit", "Cancel":
				actionColor = ui.Gray
			}
			fmt.Printf("%s[%s]%s User: %s | Prompt: %q\n", ui.Gray, timeStr, ui.Reset, entry.User, entry.UserQuery)
			fmt.Printf("  └─> Command: %s%s%s | Action: %s%s%s | Exit: %d%s\n", ui.Yellow, entry.SuggestedCommand, ui.Reset, actionColor, entry.UserAction, ui.Reset, entry.ExitCode, dangerStr)
		}
	case "2":
		fmt.Printf("%s%s%s\n", ui.Cyan, t.LogOpeningEditor, ui.Reset)
		openLogFileInEditor(logPath)
	case "q", "Q", "quit", "sair":
		return
	default:
		fmt.Printf("%s%s%s\n", ui.Red, t.LogInvalidOption, ui.Reset)
		os.Exit(1)
	}
}

func openLogFileInEditor(logPath string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd.exe", "/C", "start", "", logPath)
	case "darwin":
		cmd = exec.Command("open", logPath)
	default:
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = os.Getenv("VISUAL")
		}
		if editor != "" {
			cmd = exec.Command(editor, logPath)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		} else {
			if _, lookErr := exec.LookPath("xdg-open"); lookErr == nil {
				cmd = exec.Command("xdg-open", logPath)
			} else {
				if _, lookErr := exec.LookPath("sensible-editor"); lookErr == nil {
					cmd = exec.Command("sensible-editor", logPath)
				} else {
					cmd = exec.Command("nano", logPath)
				}
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
			}
		}
	}

	if cmd != nil {
		_ = cmd.Start()
	}
}

func HandleRulesCommand(cfg config.Config, args []string, t i18n.Translations) {
	if len(args) == 0 {
		listRules(cfg, t)
		return
	}

	subCmd := strings.ToLower(args[0])
	switch subCmd {
	case "add":
		addRule(&cfg, args[1:], t)
	case "list":
		listRules(cfg, t)
	case "edit":
		editRule(&cfg, args[1:], t)
	case "delete", "remove":
		deleteRule(&cfg, args[1:], t)
	default:
		fmt.Printf("%sInvalid rules subcommand: %s. Use: add, list, edit, delete.%s\n", ui.Red, subCmd, ui.Reset)
	}
}

func listRules(cfg config.Config, t i18n.Translations) {
	if len(cfg.Rules) == 0 {
		fmt.Println(t.NoRulesFound)
		return
	}
	fmt.Printf("%s%s%s\n", ui.Bold+ui.Cyan, t.RulesHeader, ui.Reset)
	for i, r := range cfg.Rules {
		fmt.Printf(" %d. %s\n", i+1, r)
	}
}

func addRule(cfg *config.Config, args []string, t i18n.Translations) {
	ruleText := ""
	if len(args) > 0 {
		ruleText = strings.Join(args, " ")
	} else {
		fmt.Print(t.EnterRulePrompt + ": ")
		ruleText = readLine()
	}

	if ruleText == "" {
		return
	}

	cfg.Rules = append(cfg.Rules, ruleText)
	if err := config.SaveConfig(*cfg); err != nil {
		fmt.Printf("%s%s %v%s\n", ui.Red, t.ErrorPrefix, err, ui.Reset)
		os.Exit(1)
	}
	fmt.Printf("%s✔ %s%s\n", ui.Green, t.RuleAdded, ui.Reset)
}

func editRule(cfg *config.Config, args []string, t i18n.Translations) {
	if len(cfg.Rules) == 0 {
		fmt.Println(t.NoRulesFound)
		return
	}

	var index = -1
	var ruleText string

	if len(args) > 0 {
		idx, err := strconv.Atoi(args[0])
		if err == nil && idx > 0 && idx <= len(cfg.Rules) {
			index = idx - 1
			if len(args) > 1 {
				ruleText = strings.Join(args[1:], " ")
			}
		}
	}

	if index < 0 || index >= len(cfg.Rules) {
		listRules(*cfg, t)
		fmt.Print(t.EnterIndexPrompt)
		input := readLine()
		idx, err := strconv.Atoi(input)
		if err != nil || idx <= 0 || idx > len(cfg.Rules) {
			fmt.Printf("%s%s%s\n", ui.Red, t.InvalidIndex, ui.Reset)
			return
		}
		index = idx - 1
	}

	if ruleText == "" {
		fmt.Printf("%s [%d]: ", t.EnterRulePrompt, index+1)
		ruleText = readLine()
	}

	if ruleText == "" {
		return
	}

	cfg.Rules[index] = ruleText
	if err := config.SaveConfig(*cfg); err != nil {
		fmt.Printf("%s%s %v%s\n", ui.Red, t.ErrorPrefix, err, ui.Reset)
		os.Exit(1)
	}
	fmt.Printf("%s✔ %s%s\n", ui.Green, t.RuleEdited, ui.Reset)
}

func deleteRule(cfg *config.Config, args []string, t i18n.Translations) {
	if len(cfg.Rules) == 0 {
		fmt.Println(t.NoRulesFound)
		return
	}

	var index = -1

	if len(args) > 0 {
		idx, err := strconv.Atoi(args[0])
		if err == nil && idx > 0 && idx <= len(cfg.Rules) {
			index = idx - 1
		}
	}

	if index < 0 || index >= len(cfg.Rules) {
		listRules(*cfg, t)
		fmt.Print(t.EnterIndexPrompt)
		input := readLine()
		idx, err := strconv.Atoi(input)
		if err != nil || idx <= 0 || idx > len(cfg.Rules) {
			fmt.Printf("%s%s%s\n", ui.Red, t.InvalidIndex, ui.Reset)
			return
		}
		index = idx - 1
	}

	cfg.Rules = append(cfg.Rules[:index], cfg.Rules[index+1:]...)
	if err := config.SaveConfig(*cfg); err != nil {
		fmt.Printf("%s%s %v%s\n", ui.Red, t.ErrorPrefix, err, ui.Reset)
		os.Exit(1)
	}
	fmt.Printf("%s✔ %s%s\n", ui.Green, t.RuleDeleted, ui.Reset)
}

func readLine() string {
	var buf []byte
	b := make([]byte, 1)
	for {
		n, err := ui.StdinReader.Read(b)
		if n > 0 {
			if b[0] == '\n' {
				break
			}
			buf = append(buf, b[0])
		}
		if err != nil {
			break
		}
	}
	return strings.TrimSpace(string(buf))
}
