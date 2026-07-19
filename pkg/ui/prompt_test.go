package ui

import (
	"strings"
	"testing"

	"shelloma/pkg/i18n"
)

func TestPromptActionCopyAndOtherInputs(t *testing.T) {
	trans := i18n.GetTranslations("pt")

	tests := []struct {
		input    string
		expected Action
	}{
		{"c\n", ActionCopy},
		{"copy\n", ActionCopy},
		{"copiar\n", ActionCopy},
		{"q\n", ActionQuit},
		{"sair\n", ActionQuit},
		{"y\n", ActionExecute},
		{"\n", ActionExecute},
		{"e\n", ActionExplain},
		{"m\n", ActionEdit},
	}

	for _, tt := range tests {
		reader := strings.NewReader(tt.input)
		action := PromptActionWithReader(reader, trans)
		if action != tt.expected {
			t.Errorf("Para a entrada %q, esperava ação %d, obteve %d", tt.input, tt.expected, action)
		}
	}
}
