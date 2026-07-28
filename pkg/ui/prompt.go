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

	ClearInputArea()
	DrawInputSeparator()
	DrawLegendAtBottom(legend)
	MoveToInputLine()

	fmt.Printf("%s%s%s%s", Bold, Cyan, t.OptionChoiceLabel, Reset)

	input, err := ReadFilteredInput(r)
	if err != nil {
		return ActionQuit
	}
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
	fmt.Printf("%s%s%s%s ", Bold, Cyan, t.NewCommand, Reset)

	newCmd, err := ReadFilteredInput(StdinReader)
	if err != nil || strings.TrimSpace(newCmd) == "" {
		return currentCmd
	}
	return strings.TrimSpace(newCmd)
}

func PromptSecurityWord(expectedWord string, t i18n.Translations) bool {
	prompt := fmt.Sprintf(t.SecurityWordPrompt, expectedWord)
	fmt.Printf("\n%s%s%s%s", Bold, Cyan, prompt, Reset)
	input, err := ReadFilteredInput(StdinReader)
	if err != nil {
		return false
	}
	return strings.TrimSpace(input) == expectedWord
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

func ReadFilteredInput(r io.Reader) (string, error) {
	if r != os.Stdin {
		reader := bufio.NewReader(r)
		s, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", err
		}
		return s, nil
	}
	if err := setRawMode(true); err != nil {
		reader := bufio.NewReader(r)
		return reader.ReadString('\n')
	}
	defer func() { _ = setRawMode(false) }()

	var line []rune
	reader := bufio.NewReader(os.Stdin)
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return "", err
		}
		switch r {
		case 13, 10:
			fmt.Print("\r\n")
			return string(line), nil
		case 27:
			if reader.Buffered() == 0 {
				return "", io.EOF
			}
			r2, _, err := reader.ReadRune()
			if err != nil {
				return "", io.EOF
			}
			if r2 == '[' {
				_, _, _ = reader.ReadRune()
				continue
			}
			return "", io.EOF
		case 127, 8:
			if len(line) > 0 {
				line = line[:len(line)-1]
				fmt.Print("\b \b")
			}
		default:
			if r >= 32 {
				line = append(line, r)
				fmt.Printf("%c", r)
			}
		}
	}
}

func renderLine(prompt string, line []rune, cursorPos int) {
	fmt.Print("\r\x1b[K" + prompt + string(line))
	if len(line) > cursorPos {
		fmt.Printf("\x1b[%dD", len(line)-cursorPos)
	}
}

func ReadLineWithHistory(prompt string, initialText string, history []string, t i18n.Translations) (string, error) {
	if err := setRawMode(true); err != nil {
		ClearInputArea()
		DrawInputSeparator()
		MoveToInputLine()
		fmt.Print(prompt)
		if initialText != "" {
			fmt.Print(initialText)
		}
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		DrawInputSeparator()
		MoveToContentStart()
		echo := initialText + strings.TrimSpace(input)
		if echo != "" {
			DrawContentSeparator()
			fmt.Printf("\n%s%s\n\n", prompt, echo)
		}
		return echo, nil
	}
	defer func() { _ = setRawMode(false) }()

	ClearInputArea()
	DrawInputSeparator()
	DrawLegendAtBottom(t.InputPromptLegend)
	MoveToInputLine()

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
		case 3, 4:
			fmt.Print("\r\n")
			DrawInputSeparator()
			MoveToContentStart()
			DrawContentSeparator()
			fmt.Print("\r\n")
			return "", io.EOF
		case 13, 10:
			fmt.Print("\r\n")
			DrawInputSeparator()
			MoveToContentStart()
			if len(line) > 0 {
				DrawContentSeparator()
				fmt.Printf("\r\n%s%s\r\n\r\n", prompt, string(line))
			}
			return string(line), nil
		case 127, 8:
			if cursorPos > 0 {
				line = append(line[:cursorPos-1], line[cursorPos:]...)
				cursorPos--
				renderLine(prompt, line, cursorPos)
			}
		case 27:
			if reader.Buffered() == 0 {
				fmt.Print("\r\n")
				DrawInputSeparator()
				MoveToContentStart()
				return "", io.EOF
			}
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
				case 'A':
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
				case 'B':
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
				case 'C':
					if cursorPos < len(line) {
						cursorPos++
						renderLine(prompt, line, cursorPos)
					}
				case 'D':
					if cursorPos > 0 {
						cursorPos--
						renderLine(prompt, line, cursorPos)
					}
				case '3':
					r4, _, err := reader.ReadRune()
					if err == nil && r4 == '~' {
						if cursorPos < len(line) {
							line = append(line[:cursorPos], line[cursorPos+1:]...)
							renderLine(prompt, line, cursorPos)
						}
					}
				}
				continue
			}
			fmt.Print("\r\n")
			DrawInputSeparator()
			MoveToContentStart()
			return "", io.EOF
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
	fmt.Printf("\x1b7\x1b[%d;1H\x1b[2K%s%s%s\x1b8", rows, Inverted, paddedLegend, Reset)
}

func ClearLegendAtBottom() {
	rows, _ := getTerminalSize()
	fmt.Printf("\x1b7\x1b[%d;1H\x1b[2K\x1b8", rows)
}

// Input area helpers — the bottom 4 rows are reserved as:
//   rows-3: separator line
//   rows-2: prompt / input line
//   rows-1: gap
//   rows:   footer legend

func DrawInputSeparator() {
	rows, cols := getTerminalSize()
	if rows < 5 {
		return
	}
	sep := strings.Repeat("─", cols-2)
	fmt.Printf("\x1b7\x1b[%d;1H\x1b[2K%s %s %s\x1b8", rows-3, Gray, sep, Reset)
}

func DrawContentSeparator() {
	_, cols := getTerminalSize()
	sep := strings.Repeat("─", cols-2)
	fmt.Printf("%s %s %s", Gray, sep, Reset)
}

func ClearInputArea() {
	rows, _ := getTerminalSize()
	if rows < 4 {
		return
	}
	fmt.Printf("\x1b7\x1b[%d;1H\x1b[2K\x1b8", rows-3)
	fmt.Printf("\x1b7\x1b[%d;1H\x1b[2K\x1b8", rows-2)
	fmt.Printf("\x1b7\x1b[%d;1H\x1b[2K\x1b8", rows-1)
}

func MoveToInputLine() {
	rows, _ := getTerminalSize()
	if rows < 3 {
		return
	}
	fmt.Printf("\x1b[%d;1H", rows-2)
}

func MoveToContentStart() {
	rows, _ := getTerminalSize()
	if rows > 4 {
		fmt.Printf("\x1b[%d;1H", rows-4)
	} else {
		fmt.Printf("\x1b[10;1H")
	}
}

func SetupTerminal(sysCtx sysinfo.SystemContext, model string, version string, t i18n.Translations) {
	if runtime.GOOS == "windows" {
		return
	}
	fmt.Print("\x1b[r")

	rows, cols := getTerminalSize()

	fmt.Print("\x1b[2J\x1b[H")

	fmt.Printf("%s┌%s┐%s\n", Gray, strings.Repeat("─", cols-2), Reset)

	lines := []string{
		fmt.Sprintf("      .-''-.               %sSHELLOMA - CLI ASSISTANT%s", Bold+Cyan, Reset),
		"     /   @   \\              ------------------------",
		fmt.Sprintf("    /   / \\   \\             %s: v%s", t.HeaderVersion, version),
		fmt.Sprintf("   |   /   \\   |            %s: %s", t.HeaderDirectory, sysCtx.WorkingDir),
		fmt.Sprintf("   |   \\_/    |            %s: %s (%s)", t.HeaderOS, sysCtx.OS, sysCtx.DistroName),
		fmt.Sprintf("    \\        /              %s: %s", t.HeaderShell, sysCtx.Shell),
		fmt.Sprintf("     '------'               %s: %s", t.HeaderModel, model),
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

	fmt.Printf("%s└%s┘%s\n", Gray, strings.Repeat("─", cols-2), Reset)

	// Scroll region: 1 to rows-4 (leaving separator at rows-3, input at rows-2, gap at rows-1, footer at rows)
	scrollStart := 1
	scrollEnd := rows - 4
	if scrollEnd <= scrollStart {
		scrollEnd = rows - 1
	}

	fmt.Printf("\x1b[%d;%dr", scrollStart, scrollEnd)

	DrawInputSeparator()
	ClearInputArea()
	MoveToContentStart()
}

func ResetTerminal() {
	if runtime.GOOS == "windows" {
		return
	}
	fmt.Print("\x1b[r")
	ClearInputArea()
	ClearLegendAtBottom()
}
