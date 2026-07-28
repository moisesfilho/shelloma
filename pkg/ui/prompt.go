package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"shelloma/pkg/i18n"
	"shelloma/pkg/sysinfo"
)

type Action int

const (
	ActionExecute Action = iota
	ActionExplain
	ActionEdit
	ActionCopy
	ActionQuit
	ActionNewPrompt
	ActionRefine
	ActionAdjustPrompt
)

var StdinReader io.Reader = os.Stdin

func PromptAction(t i18n.Translations) Action {
	return PromptActionWithReader(StdinReader, t)
}

func PromptActionWithReader(r io.Reader, t i18n.Translations) Action {
	legend := t.OptionsPrompt
	legend = strings.Replace(legend, "Opções: ", "", 1)
	legend = strings.Replace(legend, "Options: ", "", 1)
	legend = strings.Replace(legend, "Opciones: ", "", 1)
	legend = strings.TrimSuffix(legend, ": ")
	legend = strings.TrimSpace(legend)

	DrawLegendAtBottom(legend)

	fmt.Print(t.OptionChoiceLabel)

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
	case "p", "prompt", "new", "novo", "nuevo":
		return ActionNewPrompt
	case "r", "refine", "refinar", "complement", "complementar":
		return ActionRefine
	case "a", "adjust", "ajustar", "alterar":
		return ActionAdjustPrompt
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

func setRawMode(raw bool) error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("raw mode not supported on windows")
	}
	var cmd *exec.Cmd
	if raw {
		cmd = exec.Command("stty", "raw", "-echo")
	} else {
		cmd = exec.Command("stty", "-raw", "echo")
	}
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func renderLine(prompt string, line []rune, cursorPos int) {
	fmt.Print("\r\x1b[K" + prompt + string(line))
	if len(line) > cursorPos {
		fmt.Printf("\x1b[%dD", len(line)-cursorPos)
	}
}

func ReadLineWithHistory(prompt string, initialText string, history []string, t i18n.Translations) (string, error) {
	if err := setRawMode(true); err != nil {
		fmt.Print(prompt)
		if initialText != "" {
			fmt.Print(initialText)
		}
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return initialText + strings.TrimSpace(input), nil
	}
	defer func() { _ = setRawMode(false) }()

	DrawLegendAtBottom(t.InputPromptLegend)

	line := []rune(initialText)
	cursorPos := len(line)
	historyIndex := len(history)
	tempInput := []rune(initialText)

	renderLine(prompt, line, cursorPos)

	reader := bufio.NewReader(os.Stdin)

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return "", err
		}

		switch r {
		case 3, 4: // Ctrl+C, Ctrl+D
			fmt.Print("\r\n")
			return "", io.EOF
		case 13, 10: // Enter
			fmt.Print("\r\n")
			return string(line), nil
		case 127, 8: // Backspace
			if cursorPos > 0 {
				line = append(line[:cursorPos-1], line[cursorPos:]...)
				cursorPos--
				renderLine(prompt, line, cursorPos)
			}
		case 27: // Escape sequence
			r2, _, err := reader.ReadRune()
			if err != nil {
				continue
			}
			if r2 == '[' {
				r3, _, err := reader.ReadRune()
				if err != nil {
					continue
				}
				switch r3 {
				case 'A': // Up
					if len(history) > 0 && historyIndex > 0 {
						if historyIndex == len(history) {
							tempInput = make([]rune, len(line))
							copy(tempInput, line)
						}
						historyIndex--
						line = []rune(history[historyIndex])
						cursorPos = len(line)
						renderLine(prompt, line, cursorPos)
					}
				case 'B': // Down
					if historyIndex < len(history) {
						historyIndex++
						if historyIndex == len(history) {
							line = make([]rune, len(tempInput))
							copy(line, tempInput)
						} else {
							line = []rune(history[historyIndex])
						}
						cursorPos = len(line)
						renderLine(prompt, line, cursorPos)
					}
				case 'C': // Right
					if cursorPos < len(line) {
						cursorPos++
						renderLine(prompt, line, cursorPos)
					}
				case 'D': // Left
					if cursorPos > 0 {
						cursorPos--
						renderLine(prompt, line, cursorPos)
					}
				case '3': // Delete
					r4, _, err := reader.ReadRune()
					if err == nil && r4 == '~' {
						if cursorPos < len(line) {
							line = append(line[:cursorPos], line[cursorPos+1:]...)
							renderLine(prompt, line, cursorPos)
						}
					}
				}
			}
		default:
			if r >= 32 {
				line = append(line[:cursorPos], append([]rune{r}, line[cursorPos:]...)...)
				cursorPos++
				renderLine(prompt, line, cursorPos)
			}
		}
	}
}

func stripANSI(s string) string {
	var sb strings.Builder
	inESC := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' {
			inESC = true
			continue
		}
		if inESC {
			if s[i] == 'm' || s[i] == 'K' || s[i] == 'H' || s[i] == 'J' {
				inESC = false
			}
			continue
		}
		sb.WriteByte(s[i])
	}
	return sb.String()
}

func getTerminalSize() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 24, 80 // Default fallback
	}
	var rows, cols int
	_, err = fmt.Sscanf(strings.TrimSpace(string(out)), "%d %d", &rows, &cols)
	if err != nil || rows <= 0 || cols <= 0 {
		return 24, 80 // Default fallback
	}
	return rows, cols
}

func DrawLegendAtBottom(legend string) {
	rows, cols := getTerminalSize()
	visLen := len([]rune(stripANSI(legend)))
	padWidth := cols - visLen
	if padWidth < 0 {
		padWidth = 0
	}
	paddedLegend := legend + strings.Repeat(" ", padWidth)
	// \x1b7 (salvar cursor DEC), \x1b8 (restaurar cursor DEC)
	fmt.Printf("\x1b7\x1b[%d;1H\x1b[2K%s%s%s\x1b8", rows, Inverted, paddedLegend, Reset)
}

func ClearLegendAtBottom() {
	rows, _ := getTerminalSize()
	fmt.Printf("\x1b7\x1b[%d;1H\x1b[2K\x1b8", rows)
}

func SetupTerminal(sysCtx sysinfo.SystemContext, model string, version string, t i18n.Translations) {
	if runtime.GOOS == "windows" {
		return
	}
	// 1. Resetar qualquer região de rolagem anterior
	fmt.Print("\x1b[r")

	rows, cols := getTerminalSize()

	// 2. Limpar tela e mover cursor para o topo (1;1)
	fmt.Print("\x1b[2J\x1b[H")

	// 3. Imprimir borda superior do cabeçalho preenchendo toda a largura
	fmt.Printf("%s┌%s┐%s\n", Gray, strings.Repeat("─", cols-2), Reset)

	// 4. Imprimir cada linha de conteúdo envelopada nas bordas laterais
	lines := []string{
		fmt.Sprintf("      /\\                 %sSHELLOMA - CLI ASSISTANT%s", Bold+Cyan, Reset),
		"     /  \\                ------------------------",
		fmt.Sprintf("    / /\\ \\               %s: v%s", t.HeaderVersion, version),
		fmt.Sprintf("   ( (  ) )              %s: %s", t.HeaderDirectory, sysCtx.WorkingDir),
		fmt.Sprintf("    \\ \\/ /               %s: %s (%s)", t.HeaderOS, sysCtx.OS, sysCtx.DistroName),
		fmt.Sprintf("     \\__/                %s: %s", t.HeaderShell, sysCtx.Shell),
		fmt.Sprintf("                         %s: %s", t.HeaderModel, model),
	}

	for _, line := range lines {
		visLen := len([]rune(stripANSI(line)))
		padWidth := (cols - 4) - visLen
		if padWidth < 0 {
			padWidth = 0
		}
		padded := line + strings.Repeat(" ", padWidth)
		fmt.Printf("%s│ %s%s │%s\n", Gray, Reset+padded, Gray, Reset)
	}

	// 5. Imprimir borda inferior do cabeçalho preenchendo toda a largura
	fmt.Printf("%s└%s┘%s\n", Gray, strings.Repeat("─", cols-2), Reset)

	// Definir margens de rolagem (de 1 a penúltima linha rows-1), deixando o rodapé (linha rows) fixo
	scrollStart := 1
	scrollEnd := rows - 1
	if scrollEnd <= scrollStart {
		scrollEnd = rows
	}

	fmt.Printf("\x1b[%d;%dr", scrollStart, scrollEnd)
	fmt.Printf("\x1b[10;1H")
}

func ResetTerminal() {
	if runtime.GOOS == "windows" {
		return
	}
	// 1. Resetar margens de rolagem
	fmt.Print("\x1b[r")
	// 2. Limpar a legenda de rodapé
	ClearLegendAtBottom()
}
