package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"shelloma/pkg/config"
	"shelloma/pkg/i18n"
	"shelloma/pkg/logger"
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
			return ui.ActionCopy
		}
	case ui.ActionQuit:
		fmt.Println(t.OperationCancelled)
		return ui.ActionQuit
	}
	return action
}

func ExecuteWithRecovery(client ollama.LLMProvider, sysCtx sysinfo.SystemContext, cmdStr string, cfg config.Config, t i18n.Translations) (bool, int, string) {
	if !cfg.DisableDangerousCheck {
		if isDanger, _ := config.CheckDangerous(cmdStr, cfg.DangerousCommands); isDanger {
			if !ui.PromptSecurityWord("CONFIRM", t) {
				fmt.Printf("%s%s%s\n", ui.Red, t.SecurityWordIncorrect, ui.Reset)
				return false, 0, ""
			}
		}
	}

	exitCode, output, _ := ui.ExecuteCommand(cmdStr, t)

	fmt.Printf("%s🔍 %s%s\r", ui.Gray, t.AnalyzingResult, ui.Reset)
	analysis, errResult := client.AnalyzeExecutionResult(cmdStr, exitCode, output, sysCtx)
	fmt.Print("                                                                      \r")

	if errResult == nil && analysis.Success {
		fmt.Printf("%s✔ %s%s\n", ui.Green+ui.Bold, t.Success, ui.Reset)
		if analysis.Reason != "" && analysis.Reason != t.Success && analysis.Reason != "Comando executado com sucesso" && analysis.Reason != "Completed successfully" && analysis.Reason != "Completado con éxito" {
			fmt.Printf("%s%s%s\n", ui.Gray, analysis.Reason, ui.Reset)
		}
		return true, exitCode, output
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
			fixSuccess, fixExitCode, fixOutput := ExecuteWithRecovery(client, sysCtx, suggestedCmd, cfg, t)
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
						return false, fixExitCode, fixOutput
					}
				}
			}
			return false, fixExitCode, fixOutput
		} else if action == ui.ActionQuit {
			return false, exitCode, output
		}
	}
}

func LogExecution(userPrompt, cmdStr, action string, exitCode int, output string, sysCtx sysinfo.SystemContext, cfg config.Config, client ollama.LLMProvider) {
	isDanger, matched := config.CheckDangerous(cmdStr, cfg.DangerousCommands)
	dangerousAlertShown := !cfg.DisableDangerousCheck && isDanger

	var modelName string
	if client != nil {
		modelName = client.GetModel()
	}

	entry := logger.LogEntry{
		UserQuery:              userPrompt,
		SuggestedCommand:       cmdStr,
		UserAction:             action,
		ExitCode:               exitCode,
		CommandOutput:          output,
		WorkingDir:             sysCtx.WorkingDir,
		User:                   sysCtx.User,
		OS:                     sysCtx.OS,
		OllamaURL:              cfg.OllamaURL,
		Model:                  modelName,
		Temperature:            cfg.Temperature,
		AutoExecute:            cfg.AutoExecute,
		DangerousAlertShown:    dangerousAlertShown,
		MatchedDangerousWord:   matched,
		DangerousCheckBypassed: cfg.DisableDangerousCheck,
	}

	_ = logger.WriteLogEntry(entry)
}

func SplitCommandSteps(cmd string) []string {
	var steps []string
	for _, line := range strings.Split(cmd, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			steps = append(steps, trimmed)
		}
	}
	return steps
}

func handleCdCommand(cmdStr string, sysCtx *sysinfo.SystemContext) error {
	parts := strings.Fields(cmdStr)
	target := ""
	if len(parts) > 1 {
		target = strings.TrimSpace(strings.TrimPrefix(cmdStr, "cd"))
		target = strings.Trim(target, "\"'")
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		target = home
	}

	// Expand ~
	if strings.HasPrefix(target, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			target = filepath.Join(home, target[1:])
		}
	}

	target = os.ExpandEnv(target)

	err := os.Chdir(target)
	if err == nil {
		if wd, getErr := os.Getwd(); getErr == nil {
			sysCtx.WorkingDir = wd
		}
	}
	return err
}

func ExecuteMultiStep(client ollama.LLMProvider, sysCtx *sysinfo.SystemContext, cmdStr string, cfg config.Config, t i18n.Translations, userQuery string) (bool, int, string) {
	steps := SplitCommandSteps(cmdStr)
	if len(steps) <= 1 {
		success, ec, out := ExecuteWithRecovery(client, *sysCtx, cmdStr, cfg, t)
		LogExecution(userQuery, cmdStr, "Execute", ec, out, *sysCtx, cfg, client)
		return success, ec, out
	}

	fmt.Printf("\n%s⛓️  %s (%d %s)%s\n", ui.Bold+ui.Cyan, t.MultiStepHeader, len(steps), t.StepsLabel, ui.Reset)

	var lastExitCode = 0
	var lastOutput = ""

	for i, step := range steps {
		stepNum := i + 1
		fmt.Printf("\n%s👉 [%d/%d] %s%s\n", ui.Bold+ui.Yellow, stepNum, len(steps), step, ui.Reset)
		if !cfg.AutoExecute {
			fmt.Printf(t.ConfirmStepPrompt, stepNum, len(steps))
			reader := bufio.NewReader(ui.StdinReader)
			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(strings.ToLower(choice))
			if choice != "" && choice != "y" && choice != "yes" && choice != "sim" && choice != "si" && choice != "s" {
				fmt.Println(t.OperationCancelled)
				return false, lastExitCode, lastOutput
			}
		}

		if strings.HasPrefix(strings.TrimSpace(step), "cd ") || strings.TrimSpace(step) == "cd" {
			err := handleCdCommand(strings.TrimSpace(step), sysCtx)
			if err != nil {
				fmt.Printf("%s%s %v%s\n", ui.Red, t.ErrorPrefix, err, ui.Reset)
				LogExecution(userQuery, step, "Execute", 1, err.Error(), *sysCtx, cfg, client)
				return false, 1, err.Error()
			}
			fmt.Printf("%s✔ %s: %s%s\n", ui.Green, t.Success, sysCtx.WorkingDir, ui.Reset)
			LogExecution(userQuery, step, "Execute", 0, "Changed directory to "+sysCtx.WorkingDir, *sysCtx, cfg, client)
			lastExitCode = 0
			lastOutput = "Changed directory to " + sysCtx.WorkingDir
			continue
		}

		success, ec, out := ExecuteWithRecovery(client, *sysCtx, step, cfg, t)
		LogExecution(userQuery, step, "Execute", ec, out, *sysCtx, cfg, client)

		lastExitCode = ec
		lastOutput = out

		if !success {
			return false, ec, out
		}
	}

	return true, lastExitCode, lastOutput
}
