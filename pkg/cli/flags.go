package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"shelloma/pkg/config"
	"shelloma/pkg/i18n"
	"shelloma/pkg/ui"
)

func ParseLanguageOverride(cfg *config.Config) {
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

func SetupFlags(modelFlag, urlFlag, langFlag *string, yesFlag, verFlag *bool, t i18n.Translations, version string) {
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

func ApplyFlagOverrides(cfg *config.Config, modelFlag, urlFlag, langFlag string, yesFlag bool, t *i18n.Translations) {
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

func GetOrPromptUserQuery(args []string, ver string, t i18n.Translations) string {
	query := strings.Join(args, " ")
	if strings.TrimSpace(query) == "" {
		fmt.Printf("%s%s[Shelloma v%s]%s Prompt: ", ui.Bold, ui.Cyan, ver, ui.Reset)
		var err error
		query, err = readLineFromStdin()
		if err != nil || strings.TrimSpace(query) == "" {
			fmt.Printf("\n%s\n", t.NoInstructionProvided)
			os.Exit(0)
		}
	}
	return query
}

func readLineFromStdin() (string, error) {
	var line string
	_, err := fmt.Scanln(&line)
	return line, err
}
