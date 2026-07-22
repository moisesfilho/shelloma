package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"shelloma/pkg/i18n"
)

type Action int

const (
	ActionExecute Action = iota
	ActionExplain
	ActionEdit
	ActionCopy
	ActionQuit
)

var StdinReader io.Reader = os.Stdin

func PromptAction(t i18n.Translations) Action {
	return PromptActionWithReader(StdinReader, t)
}

func PromptActionWithReader(r io.Reader, t i18n.Translations) Action {
	fmt.Print(t.OptionsPrompt)

	reader := bufio.NewReader(r)
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

	reader := bufio.NewReader(StdinReader)
	newCmd, _ := reader.ReadString('\n')
	newCmd = strings.TrimSpace(newCmd)

	if newCmd == "" {
		return currentCmd
	}
	return newCmd
}

func PromptSecurityWord(expectedWord string, t i18n.Translations) bool {
	fmt.Printf(t.SecurityWordPrompt, expectedWord)
	reader := bufio.NewReader(StdinReader)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	return input == expectedWord
}
