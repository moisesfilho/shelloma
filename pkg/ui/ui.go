package ui

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"shelloma/pkg/i18n"
)

const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"

	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[90m"

	BgDarkGray = "\033[48;5;236m"
)

type Action int

const (
	ActionExecute Action = iota
	ActionExplain
	ActionEdit
	ActionCopy
	ActionQuit
)

func PrintBanner(modelName string, lang string) {
	fmt.Printf("%s%s[Shelloma]%s %s(Model: %s | Lang: %s)%s\n", Bold, Cyan, Reset, Gray, modelName, strings.ToUpper(lang), Reset)
}

func PrintCommandCard(cmd string) {
	lines := strings.Split(cmd, "\n")
	maxLen := 0
	for _, l := range lines {
		if len(l) > maxLen {
			maxLen = len(l)
		}
	}
	if maxLen < 40 {
		maxLen = 40
	}

	border := strings.Repeat("─", maxLen+4)

	fmt.Println()
	fmt.Printf("%s┌%s┐%s\n", Cyan, border, Reset)
	for _, l := range lines {
		padding := maxLen - len(l)
		fmt.Printf("%s│%s  %s%s%s%s  │%s\n", Cyan, Reset, Bold, Yellow, l, Reset+strings.Repeat(" ", padding), Cyan)
	}
	fmt.Printf("%s└%s┘%s\n", Cyan, border, Reset)
	fmt.Println()
}

func PromptAction(t i18n.Translations) Action {
	fmt.Print(t.OptionsPrompt)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	switch input {
	case "", "y", "sim", "yes", "si", "s":
		return ActionExecute
	case "e", "ex", "explain", "explicar":
		return ActionExplain
	case "m", "mod", "edit", "modificar":
		return ActionEdit
	case "c", "copy", "copiar":
		return ActionCopy
	case "q", "n", "no", "sair", "cancel", "cancelar", "salir":
		return ActionQuit
	default:
		if strings.HasPrefix(input, "e") {
			return ActionExplain
		}
		return ActionExecute
	}
}

func EditCommand(currentCmd string, t i18n.Translations) string {
	fmt.Printf("%s%s%s %s\n", Dim, t.CurrentCommand, Reset, currentCmd)
	fmt.Printf("%s%s%s ", Bold, t.NewCommand, Reset)

	reader := bufio.NewReader(os.Stdin)
	newCmd, _ := reader.ReadString('\n')
	newCmd = strings.TrimSpace(newCmd)

	if newCmd == "" {
		return currentCmd
	}
	return newCmd
}

func ExecuteCommand(cmdStr string, t i18n.Translations) (int, string, error) {
	fmt.Printf("%s%s⚡ %s%s\n\n", Bold, Green, t.Executing, Reset)

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	cmd := exec.Command(shell, "-c", cmdStr)

	var outBuf bytes.Buffer
	multiOut := io.MultiWriter(os.Stdout, &outBuf)
	multiErr := io.MultiWriter(os.Stderr, &outBuf)

	cmd.Stdin = os.Stdin
	cmd.Stdout = multiOut
	cmd.Stderr = multiErr

	err := cmd.Run()
	fmt.Println()

	outputStr := strings.TrimSpace(outBuf.String())
	exitCode := 0

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return exitCode, outputStr, err
}

func CopyToClipboard(text string) error {
	osc52 := fmt.Sprintf("\033]52;c;%s\a", encodeBase64(text))
	fmt.Print(osc52)

	if isCommandAvailable("wl-copy") {
		cmd := exec.Command("wl-copy")
		cmd.Stdin = strings.NewReader(text)
		return cmd.Run()
	}

	if isCommandAvailable("xclip") {
		cmd := exec.Command("xclip", "-selection", "clipboard")
		cmd.Stdin = strings.NewReader(text)
		return cmd.Run()
	}

	return nil
}

func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func encodeBase64(s string) string {
	const b64 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var result []byte
	n := len(s)
	for i := 0; i < n; i += 3 {
		var b uint32
		b = uint32(s[i]) << 16
		if i+1 < n {
			b |= uint32(s[i+1]) << 8
		}
		if i+2 < n {
			b |= uint32(s[i+2])
		}
		result = append(result, b64[(b>>18)&63])
		result = append(result, b64[(b>>12)&63])
		if i+1 < n {
			result = append(result, b64[(b>>6)&63])
		} else {
			result = append(result, '=')
		}
		if i+2 < n {
			result = append(result, b64[b&63])
		} else {
			result = append(result, '=')
		}
	}
	return string(result)
}
