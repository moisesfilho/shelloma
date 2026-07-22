package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed locales/*.json
var localeFS embed.FS

type Language string

const (
	EN Language = "en"
	PT Language = "pt"
	ES Language = "es"
)

type Translations struct {
	ProcessingWithOllama    string `json:"processing_with_ollama"`
	CommandNoValid          string `json:"command_no_valid"`
	Executing               string `json:"executing"`
	AnalyzingResult         string `json:"analyzing_result"`
	Success                 string `json:"success"`
	Failed                  string `json:"failed"`
	Reason                  string `json:"reason"`
	FixSuggestion           string `json:"fix_suggestion"`
	SuccessFixReturn        string `json:"success_fix_return"`
	OptionsPrompt           string `json:"options_prompt"`
	CurrentCommand          string `json:"current_command"`
	NewCommand              string `json:"new_command"`
	CopiedToClipboard       string `json:"copied_to_clipboard"`
	CopyError               string `json:"copy_error"`
	OperationCancelled      string `json:"operation_cancelled"`
	ExplainingWithOllama    string `json:"explaining_with_ollama"`
	ExplanationHeader       string `json:"explanation_header"`
	ModelsInstalled         string `json:"models_installed"`
	ConfigSaved             string `json:"config_saved"`
	UnknownKey              string `json:"unknown_key"`
	CurrentConfig           string `json:"current_config"`
	DefaultModelAuto        string `json:"default_model_auto"`
	ModelActive             string `json:"model_active"`
	SelectLanguagePrompt    string `json:"select_language_prompt"`
	LanguageName            string `json:"language_name"`
	OllamaLangInstruction   string `json:"ollama_lang_instruction"`
	OllamaReasonInstruction string `json:"ollama_reason_instruction"`
	OllamaOfflineError      string `json:"ollama_offline_error"`
	OllamaStartSuggestion   string `json:"ollama_start_suggestion"`
	ErrorPrefix             string `json:"error_prefix"`
	ErrorLoadingConfig      string `json:"error_loading_config"`
	NoInstructionProvided   string `json:"no_instruction_provided"`
	ExitCodeLabel           string `json:"exit_code_label"`
	NoClipboardUtility      string `json:"no_clipboard_utility"`
	OllamaNoModelsFound     string `json:"ollama_no_models_found"`
	DangerousCommandWarning string `json:"dangerous_command_warning"`
	SecurityWordPrompt      string `json:"security_word_prompt"`
	SecurityWordIncorrect   string `json:"security_word_incorrect"`
	LogPromptOptions        string `json:"log_prompt_options"`
	LogNoLogsFound          string `json:"log_no_logs_found"`
	LogOpeningEditor        string `json:"log_opening_editor"`
	LogInvalidOption        string `json:"log_invalid_option"`

	// Help / Usage
	HelpTitle       string `json:"help_title"`
	HelpUsage       string `json:"help_usage"`
	HelpCommands    string `json:"help_commands"`
	HelpOptions     string `json:"help_options"`
	FlagModelHelp   string `json:"flag_model_help"`
	FlagURLHelp     string `json:"flag_url_help"`
	FlagLangHelp    string `json:"flag_lang_help"`
	FlagYesHelp     string `json:"flag_yes_help"`
	FlagVersionHelp string `json:"flag_version_help"`
}

var loadedDict = make(map[Language]Translations)

func init() {
	for _, lang := range []Language{EN, PT, ES} {
		filePath := fmt.Sprintf("locales/%s.json", lang)
		data, err := localeFS.ReadFile(filePath)
		if err != nil {
			continue
		}
		var t Translations
		if err := json.Unmarshal(data, &t); err == nil {
			loadedDict[lang] = t
		}
	}
}

func NormalizeLanguage(lang string) Language {
	l := strings.ToLower(strings.TrimSpace(lang))
	switch {
	case strings.HasPrefix(l, "pt") || strings.Contains(l, "portugue"):
		return PT
	case strings.HasPrefix(l, "es") || strings.Contains(l, "spanis") || strings.Contains(l, "español") || strings.Contains(l, "espanol"):
		return ES
	default:
		return EN
	}
}

func GetTranslations(langStr string) Translations {
	lang := NormalizeLanguage(langStr)
	if t, ok := loadedDict[lang]; ok {
		return t
	}
	if t, ok := loadedDict[EN]; ok {
		return t
	}
	return Translations{}
}
