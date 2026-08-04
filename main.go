package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"shelloma/pkg/cli"
	"shelloma/pkg/config"
	"shelloma/pkg/i18n"
	"shelloma/pkg/sysinfo"
	"shelloma/pkg/ui"
)

const version = "1.3.0"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("%sError loading config: %v%s\n", ui.Red, err, ui.Reset)
		cli.Exit(1, i18n.Translations{})
	}

	cli.ParseLanguageOverride(&cfg)
	t := i18n.GetTranslations(cfg.Language)

	var (
		modelFlag   string
		urlFlag     string
		langFlag    string
		yesFlag     bool
		verFlag     bool
		desktopFlag bool
	)

	cli.SetupFlags(&modelFlag, &urlFlag, &langFlag, &yesFlag, &verFlag, &desktopFlag, t, version)
	flag.Parse()
	cli.IsDesktop = desktopFlag

	if verFlag {
		fmt.Printf("Shelloma v%s\n", version)
		cli.Exit(0, t)
	}

	cli.ApplyFlagOverrides(&cfg, modelFlag, urlFlag, langFlag, yesFlag, &t)
	args := flag.Args()

	if len(args) > 0 {
		switch args[0] {
		case "config":
			cli.HandleConfigCommand(cfg, args[1:], t)
			return
		case "models", "list":
			cli.HandleModelsCommand(cfg, t)
			return
		case "logs":
			cli.HandleLogsCommand(cfg, t)
			return
		case "rules":
			cli.HandleRulesCommand(cfg, args[1:], t)
			return
		case "learn", "aprender":
			cli.HandleLearnCommand(cfg, args[1:], t)
			return
		}
	}

	ui.IsTerminalApp = len(args) == 0 || cli.IsDesktop
	sysCtx := sysinfo.GetSystemContext()

	client := cli.ConnectOrRecoverOllama(cfg, t)
	if ui.IsTerminalApp {
		ui.SetupTerminal(sysCtx, client.GetModel(), version, t)
	} else {
		ui.PrintBanner(client.GetModel(), string(i18n.NormalizeLanguage(cfg.Language)))
	}

	userQuery := cli.GetOrPromptUserQuery(args, t)
	args = []string{}

	for {
		config.AddToHistory(userQuery)
		fmt.Printf("%s⏳ %s%s\n", ui.Gray, t.ProcessingWithOllama, ui.Reset)
		cmd, err := client.GenerateCommand(sysCtx, userQuery, cfg.Temperature)
		if err != nil {
			fmt.Printf("%s%s %v%s\n", ui.Red, t.ErrorPrefix, err, ui.Reset)
			cli.Exit(1, t)
		}

		if cmd == "" {
			fmt.Printf("%s%s%s\n", ui.Yellow, t.CommandNoValid, ui.Reset)
			cli.Exit(1, t)
		}

		ui.PrintCommandCard(cmd)
		if !cfg.DisableDangerousCheck {
			if isDanger, matched := config.CheckDangerous(cmd, cfg.DangerousCommands); isDanger {
				ui.PrintDangerousWarning(matched, t)
			}
		}

		if cfg.AutoExecute {
			success, _, _ := cli.ExecuteMultiStep(client, &sysCtx, cmd, cfg, t, userQuery)
			if success {
				cli.Exit(0, t)
			}
			cli.Exit(1, t)
		}

		actionNeeded := true
		for actionNeeded {
			action := cli.HandleUserAction(client, sysCtx, &cmd, cfg, t)
			switch action {
			case ui.ActionExecute:
				success, _, _ := cli.ExecuteMultiStep(client, &sysCtx, cmd, cfg, t, userQuery)
				if !ui.IsTerminalApp {
					if success {
						exitApp(0)
					}
					exitApp(1)
				}
				history, _ := config.LoadHistory()
				ans, err := ui.ReadLineWithHistory(fmt.Sprintf("%s%s%s%s", ui.Bold, ui.Cyan, t.AnythingElsePrompt, ui.Reset), "", history, t)
				if err == io.EOF || strings.TrimSpace(ans) == "" {
					if success {
						exitApp(0)
					}
					exitApp(1)
				}
				userQuery = ans
				actionNeeded = false
			case ui.ActionQuit:
				cli.LogExecution(userQuery, cmd, "Quit", 0, "", sysCtx, cfg, client)
				exitApp(0) // No pause on explicit Quit
			case ui.ActionCopy:
				cli.LogExecution(userQuery, cmd, "Copy", 0, "", sysCtx, cfg, client)
				cli.Exit(0, t)
			case ui.ActionNewPrompt:
				userQuery = cli.GetOrPromptUserQuery(args, t)
				actionNeeded = false
			case ui.ActionRefine:
				history, _ := config.LoadHistory()
				feedback, err := ui.ReadLineWithHistory(fmt.Sprintf("%s%s%s%s", ui.Bold, ui.Cyan, t.RefinePromptLabel, ui.Reset), "", history, t)
				if err == io.EOF || strings.TrimSpace(feedback) == "" {
					exitApp(0)
				}
				fmt.Printf("%s⏳ %s%s\n", ui.Gray, t.ProcessingWithOllama, ui.Reset)
				refinedCmd, err := client.GenerateRefinedCommand(sysCtx, userQuery, cmd, feedback, cfg.Temperature)
				switch {
				case err != nil:
					fmt.Printf("\n%s%s %v%s\n", ui.Red, t.ErrorPrefix, err, ui.Reset)
				case refinedCmd == "":
					fmt.Printf("\n%s%s%s\n", ui.Yellow, t.CommandNoValid, ui.Reset)
				default:
					cmd = refinedCmd
					ui.PrintCommandCard(cmd)
					if !cfg.DisableDangerousCheck {
						if isDanger, matched := config.CheckDangerous(cmd, cfg.DangerousCommands); isDanger {
							ui.PrintDangerousWarning(matched, t)
						}
					}
				}
			case ui.ActionAdjustPrompt:
				promptStr := fmt.Sprintf("%s%s%s%s", ui.Bold, ui.Cyan, t.InitialPromptLabel, ui.Reset)
				history, _ := config.LoadHistory()
				newQuery, err := ui.ReadLineWithHistory(promptStr, userQuery, history, t)
				if err == nil && strings.TrimSpace(newQuery) != "" {
					userQuery = newQuery
					actionNeeded = false
				}
			}
		}
	}
}

func exitApp(code int) {
	ui.ClearInputArea()
	ui.ClearLegendAtBottom()
	os.Exit(code)
}
