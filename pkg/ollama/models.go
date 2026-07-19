package ollama

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) ListModels() ([]string, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/api/tags")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama status %d", resp.StatusCode)
	}

	var tags TagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return nil, err
	}

	var names []string
	for _, m := range tags.Models {
		names = append(names, m.Name)
	}
	return names, nil
}

func selectBestModel(models []string) string {
	preferences := []string{
		"qwen2.5-coder",
		"deepseek-coder",
		"llama3.2",
		"deepseek-r1",
		"codellama",
		"mistral",
	}

	for _, pref := range preferences {
		for _, m := range models {
			if strings.Contains(strings.ToLower(m), pref) {
				return m
			}
		}
	}

	return models[0]
}
