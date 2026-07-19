package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"shelloma/pkg/config"
	"shelloma/pkg/i18n"
	"shelloma/pkg/ollama"
	"shelloma/pkg/sysinfo"
	"shelloma/pkg/ui"
)

const version = "1.0.2"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("%sError loading config: %v%s\n", ui.Red, err, ui.Reset)
		os.Exit(1)
	}

	parseLanguageOverride(&cfg)
	t := i18n.GetTranslations(cfg.Language)

	var (
		modelFlag string
		urlFlag   string
		langFlag  string
		yesFlag   bool
		verFlag   bool
	)

	setupFlags(&modelFlag, &urlFlag, &langFlag, &yesFlag, &verFlag, t)
	flag.Parse()

	if verFlag {
		fmt.Printf("Shelloma v%s\n", version)
		os.Exit(0)
	}

	applyFlagOverrides(&cfg, modelFlag, urlFlag, langFlag, yesFlag, &t)
	args := flag.Args()

	if len(args) > 0 {
		switch args[0] {
		case "config":
			handleConfigCommand(cfg, args[1:], t)
			return
		case "models", "list":
			handleModelsCommand(cfg, t)
			return
		}
	}

	userQuery := getOrPromptUserQuery(args, version)
	sysCtx := sysinfo.GetSystemContext()

	client := connectOrRecoverOllama(cfg, t)
	ui.PrintBanner(client.GetModel(), string(i18n.NormalizeLanguage(cfg.Language)))

	fmt.Printf("%s⏳ %s%s\r", ui.Gray, t.ProcessingWithOllama, ui.Reset)
	cmd, err := client.GenerateCommand(sysCtx, userQuery, cfg.Temperature)
	if err != nil {
		fmt.Printf("\n%sError: %v%s\n", ui.Red, err, ui.Reset)
		os.Exit(1)
	}
	fmt.Print("                                                                \r")

	if cmd == "" {
		fmt.Printf("%s%s%s\n", ui.Yellow, t.CommandNoValid, ui.Reset)
		os.Exit(1)
	}

	ui.PrintCommandCard(cmd)

	if cfg.AutoExecute {
		if executeWithRecovery(client, sysCtx, cmd, true, t) {
			os.Exit(0)
		}
		os.Exit(1)
	}

	for {
		action := handleUserAction(client, sysCtx, &cmd, t)
		if action == ui.ActionExecute {
			if executeWithRecovery(client, sysCtx, cmd, false, t) {
				os.Exit(0)
			}
			os.Exit(1)
		} else if action == ui.ActionQuit {
			os.Exit(0)
		}
	}
}

func parseLanguageOverride(cfg *config.Config) {
	for i, arg := range os.Args {
		if (arg == "-l" || arg == "--lang" || strings.HasPrefix(arg, "-l=") || strings.HasPrefix(arg, "--lang=")) && i+1 < len(os.Args) {
			val := os.Args[i+1]
			if strings.Contains(arg, "=") {
				parts := strings.Split(arg, "=")
				val = parts[1]
			}
			cfg.Language = string(i18n.NormalizeLanguage(val))
		}
	}
}

func setupFlags(modelFlag, urlFlag, langFlag *string, yesFlag, verFlag *bool, t i18n.Translations) {
	flag.StringVar(modelFlag, "m", "", t.FlagModelHelp)
	flag.StringVar(modelFlag, "model", "", t.FlagModelHelp)
	flag.StringVar(urlFlag, "url", "", t.FlagURLHelp)
	flag.StringVar(langFlag, "l", "", t.FlagLangHelp)
	flag.StringVar(langFlag, "lang", "", t.FlagLangHelp)
	flag.BoolVar(yesFlag, "y", false, t.FlagYesHelp)
	flag.BoolVar(yesFlag, "yes", false, t.FlagYesHelp)
	flag.BoolVar(verFlag, "v", false, t.FlagVersionHelp)
	flag.BoolVar(verFlag, "version", false, t.FlagVersionHelp)

	flag.Usage = func() {
		fmt.Printf("%sShelloma v%s%s - %s\n\n", ui.Bold+ui.Cyan, version, ui.Reset, t.HelpTitle)
		fmt.Println(t.HelpUsage)
		fmt.Println("  shelloma \"instruction\"")
		fmt.Println("  shelloma [options] \"instruction\"")
		fmt.Println()
		fmt.Println(t.HelpCommands)
		fmt.Println("  shelloma models          List installed Ollama models")
		fmt.Println("  shelloma config          Show current configuration")
		fmt.Println("  shelloma config set model <model_name>")
		fmt.Println("  shelloma config set lang <en|pt|es>")
		fmt.Println()
		fmt.Println(t.HelpOptions)
		flag.PrintDefaults()
	}
}

func applyFlagOverrides(cfg *config.Config, modelFlag, urlFlag, langFlag string, yesFlag bool, t *i18n.Translations) {
	if modelFlag != "" {
		cfg.Model = modelFlag
	}
	if urlFlag != "" {
		cfg.OllamaURL = urlFlag
	}
	if langFlag != "" {
		cfg.Language = string(i18n.NormalizeLanguage(langFlag))
		*t = i18n.GetTranslations(cfg.Language)
	}
	if yesFlag {
		cfg.AutoExecute = true
	}
}

func getOrPromptUserQuery(args []string, ver string) string {
	query := strings.Join(args, " ")
	if strings.TrimSpace(query) == "" {
		fmt.Printf("%s%s[Shelloma v%s]%s Prompt: ", ui.Bold, ui.Cyan, ver, ui.Reset)
		var err error
		query, err = readLineFromStdin()
		if err != nil || strings.TrimSpace(query) == "" {
			fmt.Println("\nNo instruction provided. Exiting.")
			os.Exit(0)
		}
	}
	return query
}

func connectOrRecoverOllama(cfg config.Config, t i18n.Translations) ollama.LLMProvider {
	client, err := ollama.NewClient(cfg)
	if err == nil {
		return client
	}

	fmt.Printf("%s✖ %s%s\n", ui.Red+ui.Bold, t.OllamaOfflineError, ui.Reset)
	fmt.Printf("%s%s%s %v\n", ui.Yellow+ui.Bold, t.Reason, ui.Reset, err)

	startCmd := "sudo systemctl start ollama"
	if _, lookErr := exec.LookPath("systemctl"); lookErr != nil {
		startCmd = "ollama serve"
	}

	fmt.Printf("\n%s💡 %s%s\n", ui.Cyan+ui.Bold, t.OllamaStartSuggestion, ui.Reset)
	ui.PrintCommandCard(startCmd)

	for {
		action := handleUserAction(nil, sysinfo.SystemContext{}, &startCmd, t)
		if action == ui.ActionExecute {
			exitCode, _, _ := ui.ExecuteCommand(startCmd, t)
			if exitCode == 0 {
				fmt.Printf("%s✔ %s%s\n", ui.Green+ui.Bold, t.Success, ui.Reset)
				fmt.Printf("%s⏳ %s%s\r", ui.Gray, t.ProcessingWithOllama, ui.Reset)

				var retryErr error
				for attempt := 1; attempt <= 6; attempt++ {
					time.Sleep(1 * time.Second)
					client, retryErr = ollama.NewClient(cfg)
					if retryErr == nil {
						fmt.Print("                                                                \r")
						return client
					}
				}
				fmt.Print("                                                                \r")
				fmt.Printf("%s%v%s\n", ui.Red, retryErr, ui.Reset)
			}
			os.Exit(1)
		} else if action == ui.ActionQuit {
			os.Exit(1)
		}
	}
}

func handleUserAction(client ollama.LLMProvider, sysCtx sysinfo.SystemContext, cmd *string, t i18n.Translations) ui.Action {
	action := ui.PromptAction(t)
	switch action {
	case ui.ActionExecute:
		return ui.ActionExecute
	case ui.ActionExplain:
		if client != nil {
			fmt.Printf("\n%s⏳ %s%s\r", ui.Gray, t.ExplainingWithOllama, ui.Reset)
			explanation, err := client.ExplainCommand(*cmd)
			fmt.Print("                                                \r")
			if err != nil {
				fmt.Printf("%sError: %v%s\n\n", ui.Red, err, ui.Reset)
			} else {
				fmt.Printf("%sℹ %s%s\n%s\n\n", ui.Cyan+ui.Bold, t.ExplanationHeader, ui.Reset, explanation)
			}
		} else {
			fmt.Printf("\n%sℹ %s%s\nInicia o serviço do Ollama no sistema Linux.\n\n", ui.Cyan+ui.Bold, t.ExplanationHeader, ui.Reset)
		}
	case ui.ActionEdit:
		*cmd = ui.EditCommand(*cmd, t)
		ui.PrintCommandCard(*cmd)
	case ui.ActionCopy:
		if err := ui.CopyToClipboard(*cmd); err != nil {
			fmt.Printf("%s%s %v%s\n", ui.Red, t.CopyError, err, ui.Reset)
		} else {
			fmt.Printf("%s✔ %s%s\n\n", ui.Green, t.CopiedToClipboard, ui.Reset)
		}
	case ui.ActionQuit:
		fmt.Println(t.OperationCancelled)
		return ui.ActionQuit
	}
	return action
}

func executeWithRecovery(client ollama.LLMProvider, sysCtx sysinfo.SystemContext, cmdStr string, autoExec bool, t i18n.Translations) bool {
	exitCode, output, _ := ui.ExecuteCommand(cmdStr, t)

	fmt.Printf("%s🔍 %s%s\r", ui.Gray, t.AnalyzingResult, ui.Reset)
	analysis, err := client.AnalyzeExecutionResult(cmdStr, exitCode, output, sysCtx)
	fmt.Print("                                                                      \r")

	if err == nil && analysis.Success {
		fmt.Printf("%s✔ %s%s\n", ui.Green+ui.Bold, t.Success, ui.Reset)
		if analysis.Reason != "" && analysis.Reason != "Comando executado com sucesso" && analysis.Reason != "Completed successfully" {
			fmt.Printf("%s%s%s\n", ui.Gray, analysis.Reason, ui.Reset)
		}
		return true
	}

	fmt.Printf("%s✖ %s%s\n", ui.Red+ui.Bold, t.Failed, ui.Reset)
	if analysis.Reason != "" {
		fmt.Printf("%s%s%s %s\n", ui.Yellow+ui.Bold, t.Reason, ui.Reset, analysis.Reason)
	} else {
		fmt.Printf("%s%s%s Exit Code: %d\n", ui.Yellow+ui.Bold, t.Reason, ui.Reset, exitCode)
	}

	suggestedCmd := analysis.SuggestedCommand
	if suggestedCmd == "" || suggestedCmd == cmdStr {
		fmt.Printf("%s⏳ %s%s\r", ui.Gray, t.ProcessingWithOllama, ui.Reset)
		var fixErr error
		suggestedCmd, fixErr = client.GenerateFixCommand(sysCtx, cmdStr, output)
		fmt.Print("                                                       \r")
		if fixErr != nil || suggestedCmd == "" {
			suggestedCmd = "ls -la"
		}
	}

	fmt.Printf("\n%s💡 %s%s\n", ui.Cyan+ui.Bold, t.FixSuggestion, ui.Reset)
	ui.PrintCommandCard(suggestedCmd)

	for {
		action := handleUserAction(client, sysCtx, &suggestedCmd, t)
		if action == ui.ActionExecute {
			fixSuccess := executeWithRecovery(client, sysCtx, suggestedCmd, autoExec, t)
			if fixSuccess {
				fmt.Printf("\n%s🔄 %s%s\n", ui.Bold+ui.Cyan, t.SuccessFixReturn, ui.Reset)
				ui.PrintCommandCard(cmdStr)

				for {
					prevAction := handleUserAction(client, sysCtx, &cmdStr, t)
					if prevAction == ui.ActionExecute {
						return executeWithRecovery(client, sysCtx, cmdStr, autoExec, t)
					} else if prevAction == ui.ActionQuit {
						return false
					}
				}
			}
			return false
		} else if action == ui.ActionQuit {
			return false
		}
	}
}

func handleModelsCommand(cfg config.Config, t i18n.Translations) {
	client := connectOrRecoverOllama(cfg, t)

	models, err := client.ListModels()
	if err != nil {
		fmt.Printf("%sError: %v%s\n", ui.Red, err, ui.Reset)
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

func handleConfigCommand(cfg config.Config, args []string, t i18n.Translations) {
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
		default:
			fmt.Printf("%s%s: %s%s\n", ui.Red, t.UnknownKey, key, ui.Reset)
			os.Exit(1)
		}

		if err := config.SaveConfig(cfg); err != nil {
			fmt.Printf("%sError: %v%s\n", ui.Red, err, ui.Reset)
			os.Exit(1)
		}

		fmt.Printf("%s✔ %s (%s = %s)%s\n", ui.Green, t.ConfigSaved, key, val, ui.Reset)
		return
	}

	path, _ := config.GetConfigPath()
	fmt.Printf("%s%s (%s):%s\n", ui.Bold+ui.Cyan, t.CurrentConfig, path, ui.Reset)
	fmt.Printf("  ollama_url:   %s\n", cfg.OllamaURL)
	fmt.Printf("  model:        %s %s\n", cfg.Model, t.DefaultModelAuto)
	fmt.Printf("  language:     %s (%s)\n", cfg.Language, i18n.GetTranslations(cfg.Language).LanguageName)
	fmt.Printf("  temperature:  %.2f\n", cfg.Temperature)
	fmt.Printf("  auto_execute: %t\n", cfg.AutoExecute)
}

func readLineFromStdin() (string, error) {
	var line string
	_, err := fmt.Scanln(&line)
	return line, err
}
