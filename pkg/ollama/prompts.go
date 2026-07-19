package ollama

import (
	"strings"
)

func isValidShellCommand(cmd string) bool {
	if cmd == "" {
		return false
	}
	lowered := strings.ToLower(cmd)
	badWords := []string{"você", "pode", "tente", "criar", "verifique", "arquivo", "diretório", "caso", "este", "se ", "you ", "can ", "try "}
	for _, bw := range badWords {
		if strings.Contains(lowered, bw+" ") || strings.HasSuffix(lowered, bw) {
			return false
		}
	}
	return true
}

func cleanJSONOutput(output string) string {
	cleaned := strings.TrimSpace(output)
	if strings.HasPrefix(cleaned, "```json") {
		cleaned = strings.TrimPrefix(cleaned, "```json")
		cleaned = strings.TrimSuffix(cleaned, "```")
	} else if strings.HasPrefix(cleaned, "```") {
		cleaned = strings.TrimPrefix(cleaned, "```")
		cleaned = strings.TrimSuffix(cleaned, "```")
	}
	return strings.TrimSpace(cleaned)
}

func cleanCommandOutput(output string) string {
	cleaned := strings.TrimSpace(output)

	if strings.HasPrefix(cleaned, "```") {
		lines := strings.Split(cleaned, "\n")
		if len(lines) > 1 {
			if strings.HasPrefix(lines[0], "```") {
				lines = lines[1:]
			}
			if len(lines) > 0 && strings.HasPrefix(lines[len(lines)-1], "```") {
				lines = lines[:len(lines)-1]
			}
			cleaned = strings.Join(lines, "\n")
		}
	}

	cleaned = strings.TrimPrefix(cleaned, "`")
	cleaned = strings.TrimSuffix(cleaned, "`")
	return strings.TrimSpace(cleaned)
}
