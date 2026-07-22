package cli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"shelloma/pkg/config"
	"shelloma/pkg/i18n"
	"shelloma/pkg/ollama"
	"shelloma/pkg/sysinfo"
	"shelloma/pkg/ui"
)

func ConnectOrRecoverOllama(cfg config.Config, t i18n.Translations) ollama.LLMProvider {
	client, err := ollama.NewClient(cfg)
	if err == nil {
		return client
	}

	fmt.Printf("%s✖ %s%s\n", ui.Red+ui.Bold, t.OllamaOfflineError, ui.Reset)
	fmt.Printf("%s%s%s %v\n", ui.Yellow+ui.Bold, t.Reason, ui.Reset, err)

	startCmd := "ollama serve"
	if runtime.GOOS == "linux" {
		if _, lookErr := exec.LookPath("systemctl"); lookErr == nil {
			startCmd = "sudo systemctl start ollama"
		}
	} else if runtime.GOOS == "darwin" {
		if _, lookErr := exec.LookPath("brew"); lookErr == nil {
			startCmd = "brew services start ollama"
		}
	}

	fmt.Printf("\n%s💡 %s%s\n", ui.Cyan+ui.Bold, t.OllamaStartSuggestion, ui.Reset)
	ui.PrintCommandCard(startCmd)

	for {
		action := HandleUserAction(nil, sysinfo.SystemContext{}, &startCmd, cfg, t)
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

func HandleUserAction(client ollama.LLMProvider, _ sysinfo.SystemContext, cmd *string, cfg config.Config, t i18n.Translations) ui.Action {
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
				fmt.Printf("%s%s %v%s\n\n", ui.Red, t.ErrorPrefix, err, ui.Reset)
			} else {
				fmt.Printf("%sℹ %s%s\n%s\n\n", ui.Cyan+ui.Bold, t.ExplanationHeader, ui.Reset, explanation)
			}
		} else {
			fmt.Printf("\n%sℹ %s%s\n%s\n\n", ui.Cyan+ui.Bold, t.ExplanationHeader, ui.Reset, t.OllamaStartSuggestion)
		}
	case ui.ActionEdit:
		*cmd = ui.EditCommand(*cmd, t)
		ui.PrintCommandCard(*cmd)
		if !cfg.DisableDangerousCheck {
			if isDanger, matched := config.CheckDangerous(*cmd, cfg.DangerousCommands); isDanger {
				ui.PrintDangerousWarning(matched, t)
			}
		}
	case ui.ActionCopy:
		if err := ui.CopyToClipboard(*cmd, t); err != nil {
			fmt.Printf("%s%s %v%s\n", ui.Red, t.CopyError, err, ui.Reset)
		} else {
			fmt.Printf("%s✔ %s%s\n", ui.Green, t.CopiedToClipboard, ui.Reset)
			return ui.ActionQuit
		}
	case ui.ActionQuit:
		fmt.Println(t.OperationCancelled)
		return ui.ActionQuit
	}
	return action
}

func ExecuteWithRecovery(client ollama.LLMProvider, sysCtx sysinfo.SystemContext, cmdStr string, cfg config.Config, t i18n.Translations) bool {
	if !cfg.DisableDangerousCheck {
		if isDanger, _ := config.CheckDangerous(cmdStr, cfg.DangerousCommands); isDanger {
			if !ui.PromptSecurityWord("CONFIRM", t) {
				fmt.Printf("%s%s%s\n", ui.Red, t.SecurityWordIncorrect, ui.Reset)
				return false
			}
		}
	}

	exitCode, output, _ := ui.ExecuteCommand(cmdStr, t)

	fmt.Printf("%s🔍 %s%s\r", ui.Gray, t.AnalyzingResult, ui.Reset)
	analysis, err := client.AnalyzeExecutionResult(cmdStr, exitCode, output, sysCtx)
	fmt.Print("                                                                      \r")

	if err == nil && analysis.Success {
		fmt.Printf("%s✔ %s%s\n", ui.Green+ui.Bold, t.Success, ui.Reset)
		if analysis.Reason != "" && analysis.Reason != t.Success && analysis.Reason != "Comando executado com sucesso" && analysis.Reason != "Completed successfully" && analysis.Reason != "Completado con éxito" {
			fmt.Printf("%s%s%s\n", ui.Gray, analysis.Reason, ui.Reset)
		}
		return true
	}

	fmt.Printf("%s✖ %s%s\n", ui.Red+ui.Bold, t.Failed, ui.Reset)
	if analysis.Reason != "" {
		fmt.Printf("%s%s%s %s\n", ui.Yellow+ui.Bold, t.Reason, ui.Reset, analysis.Reason)
	} else {
		fmt.Printf("%s%s%s %s %d\n", ui.Yellow+ui.Bold, t.Reason, ui.Reset, t.ExitCodeLabel, exitCode)
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
	if !cfg.DisableDangerousCheck {
		if isDanger, matched := config.CheckDangerous(suggestedCmd, cfg.DangerousCommands); isDanger {
			ui.PrintDangerousWarning(matched, t)
		}
	}

	for {
		action := HandleUserAction(client, sysCtx, &suggestedCmd, cfg, t)
		if action == ui.ActionExecute {
			fixSuccess := ExecuteWithRecovery(client, sysCtx, suggestedCmd, cfg, t)
			if fixSuccess {
				fmt.Printf("\n%s🔄 %s%s\n", ui.Bold+ui.Cyan, t.SuccessFixReturn, ui.Reset)
				ui.PrintCommandCard(cmdStr)
				if !cfg.DisableDangerousCheck {
					if isDanger, matched := config.CheckDangerous(cmdStr, cfg.DangerousCommands); isDanger {
						ui.PrintDangerousWarning(matched, t)
					}
				}

				for {
					prevAction := HandleUserAction(client, sysCtx, &cmdStr, cfg, t)
					if prevAction == ui.ActionExecute {
						return ExecuteWithRecovery(client, sysCtx, cmdStr, cfg, t)
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
