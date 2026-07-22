package ui

import (
	"fmt"
	"strings"

	"shelloma/pkg/i18n"
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

func PrintDangerousWarning(matchedCmd string, t i18n.Translations) {
	fmt.Printf("%s%s⚠️  %s%s\n", Bold, Red, fmt.Sprintf(t.DangerousCommandWarning, matchedCmd), Reset)
}
