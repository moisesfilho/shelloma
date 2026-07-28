package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHistoryFlow(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "shelloma_test_history")
	if err != nil {
		t.Fatalf("Erro ao criar dir temporario: %v", err)
	}
	defer os.RemoveAll(tempDir)

	t.Setenv("HOME", tempDir)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tempDir, ".config"))
	t.Setenv("AppData", filepath.Join(tempDir, "AppData", "Roaming"))

	// 1. Carregar historico inicial quando nao existe arquivo
	hist, err := LoadHistory()
	if err != nil {
		t.Fatalf("Erro ao carregar historico vazio: %v", err)
	}
	if len(hist) != 0 {
		t.Errorf("Esperava historico de tamanho 0, obteve %d", len(hist))
	}

	// 2. Salvar historico simples
	items := []string{"git status", "ls -la", "echo 'hello'"}
	err = SaveHistory(items)
	if err != nil {
		t.Fatalf("Erro ao salvar historico: %v", err)
	}

	hist, err = LoadHistory()
	if err != nil {
		t.Fatalf("Erro ao carregar historico: %v", err)
	}
	if len(hist) != len(items) {
		t.Fatalf("Esperava tamanho %d, obteve %d", len(items), len(hist))
	}
	for i, v := range items {
		if hist[i] != v {
			t.Errorf("Esperava item %d = %q, obteve %q", i, v, hist[i])
		}
	}

	// 3. Adicionar novos itens via AddToHistory
	AddToHistory("cat file.txt")
	hist, _ = LoadHistory()
	if len(hist) != 4 || hist[3] != "cat file.txt" {
		t.Errorf("Erro ao adicionar item ao historico: %v", hist)
	}

	// 4. Testar deduplicacao consecutiva
	AddToHistory("cat file.txt")
	hist, _ = LoadHistory()
	if len(hist) != 4 {
		t.Errorf("Esperava que item repetido consecutivo nao fosse adicionado. Tamanho: %d", len(hist))
	}

	// 5. Adicionar prompt vazio nao deve alterar o historico
	AddToHistory("")
	hist, _ = LoadHistory()
	if len(hist) != 4 {
		t.Errorf("Prompt vazio nao deveria alterar tamanho. Tamanho: %d", len(hist))
	}

	// 6. Testar limite de 1000 itens ao salvar
	largeItems := make([]string, 1005)
	for i := range largeItems {
		largeItems[i] = "item"
	}
	err = SaveHistory(largeItems)
	if err != nil {
		t.Fatalf("Erro ao salvar grande historico: %v", err)
	}
	hist, _ = LoadHistory()
	if len(hist) != 1000 {
		t.Errorf("Esperava limite de 1000 itens, obteve %d", len(hist))
	}
}
