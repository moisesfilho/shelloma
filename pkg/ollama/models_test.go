package ollama

import (
	"testing"
)

func TestSelectBestModel(t *testing.T) {
	models := []string{"llama2:latest", "qwen2.5-coder:1.5b", "mistral:latest"}
	best := selectBestModel(models)
	if best != "qwen2.5-coder:1.5b" {
		t.Errorf("Esperava qwen2.5-coder:1.5b, obteve %s", best)
	}

	fallbackModels := []string{"custom-model:v1"}
	bestFallback := selectBestModel(fallbackModels)
	if bestFallback != "custom-model:v1" {
		t.Errorf("Esperava custom-model:v1, obteve %s", bestFallback)
	}
}
