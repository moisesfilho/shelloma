package main

import (
	"flag"
	"fmt"
	"os"

	"shelloma/pkg/cli"
	"shelloma/pkg/config"
	"shelloma/pkg/i18n"
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

	cli.ParseLanguageOverride(&cfg)
	t := i18n.GetTranslations(cfg.Language)

	var (
		modelFlag string
		urlFlag   string
		langFlag  string
		yesFlag   bool
		verFlag   bool
	)

	cli.SetupFlags(&modelFlag, &urlFlag, &langFlag, &yesFlag, &verFlag, t, version)
	flag.Parse()

	if verFlag {
		fmt.Printf("Shelloma v%s\n", version)
		os.Exit(0)
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
		}
	}

	userQuery := cli.GetOrPromptUserQuery(args, version, t)
	sysCtx := sysinfo.GetSystemContext()

	client := cli.ConnectOrRecoverOllama(cfg, t)
	ui.PrintBanner(client.GetModel(), string(i18n.NormalizeLanguage(cfg.Language)))

	fmt.Printf("%s⏳ %s%s\r", ui.Gray, t.ProcessingWithOllama, ui.Reset)
	cmd, err := client.GenerateCommand(sysCtx, userQuery, cfg.Temperature)
	if err != nil {
		fmt.Printf("\n%s%s %v%s\n", ui.Red, t.ErrorPrefix, err, ui.Reset)
		os.Exit(1)
	}
	fmt.Print("                                                                \r")

	if cmd == "" {
		fmt.Printf("%s%s%s\n", ui.Yellow, t.CommandNoValid, ui.Reset)
		os.Exit(1)
	}

	ui.PrintCommandCard(cmd)

	if cfg.AutoExecute {
		if cli.ExecuteWithRecovery(client, sysCtx, cmd, true, t) {
			os.Exit(0)
		}
		os.Exit(1)
	}

	for {
		action := cli.HandleUserAction(client, sysCtx, &cmd, t)
		if action == ui.ActionExecute {
			if cli.ExecuteWithRecovery(client, sysCtx, cmd, false, t) {
				os.Exit(0)
			}
			os.Exit(1)
		} else if action == ui.ActionQuit {
			os.Exit(0)
		}
	}
}
