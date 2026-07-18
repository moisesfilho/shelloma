package i18n

import (
	"testing"
)

func TestNormalizeLanguage(t *testing.T) {
	tests := []struct {
		input    string
		expected Language
	}{
		{"en", EN},
		{"EN", EN},
		{"english", EN},
		{"pt", PT},
		{"pt-BR", PT},
		{"portugues", PT},
		{"es", ES},
		{"es-ES", ES},
		{"español", ES},
		{"spanish", ES},
		{"unknown", EN}, // Fallback default: English
	}

	for _, tt := range tests {
		got := NormalizeLanguage(tt.input)
		if got != tt.expected {
			t.Errorf("NormalizeLanguage(%q) = %q; esperava %q", tt.input, got, tt.expected)
		}
	}
}

func TestGetTranslations(t *testing.T) {
	enTrans := GetTranslations("en")
	if enTrans.LanguageName != "English" {
		t.Errorf("Esperava English, obteve %s", enTrans.LanguageName)
	}

	ptTrans := GetTranslations("pt-BR")
	if ptTrans.LanguageName != "Português Brasileiro" {
		t.Errorf("Esperava Português Brasileiro, obteve %s", ptTrans.LanguageName)
	}

	esTrans := GetTranslations("es")
	if esTrans.LanguageName != "Español" {
		t.Errorf("Esperava Español, obteve %s", esTrans.LanguageName)
	}
}
