package summarizer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/oliveagle/jsonpath"
)

// Summarizer is an interface for summarizing text
type Summarizer interface {
	Summarize(text, prompt string) (string, []string, error)
}

// OpenAISummarizer is an implementation of Summarizer using OpenAI
type OpenAISummarizer struct {
	APIKey string
	Debug  bool
}

const (
	PlaceholerText = `{{Text}}`
	PlaceholerPath = `{{Obsidian_Vault_Path}}`
)

// Summarize generates a summary using the OpenAI API
func (s *OpenAISummarizer) Summarize(text, filepath, prompt string, warn func(string)) (string, []string, error) {
	url := "https://api.openai.com/v1/responses"

	if strings.Contains(prompt, PlaceholerText) {
		prompt = strings.ReplaceAll(prompt, PlaceholerText, text)
	} else {
		return "", nil, fmt.Errorf("prompt must contain " + PlaceholerText + " placeholder")
	}

	if strings.Contains(prompt, PlaceholerPath) {
		prompt = strings.ReplaceAll(prompt, PlaceholerPath, filepath)
	} else {
		warn("Warning: prompt does not contain " + PlaceholerPath + " placeholder")
	}

	escapedPrompt, err := json.Marshal(prompt)
	if err != nil {
		return "", nil, fmt.Errorf("failed to escape prompt to JSON: %w", err)
	}

	payload := fmt.Sprintf(`{
		"model": "gpt-4o-mini",
		"input": [
			{
				"role": "user",
				"content": [
					{
						"type": "input_text",
						"text": %s
					}
				]
			}
		],
		"text": {
			"format": {
				"type": "json_schema",
				"name": "text_summary",
				"strict": true,
				"schema": {
					"type": "object",
					"properties": {
						"summary": {
							"type": "string",
							"description": "A summary of the text."
						},
						"tags": {
							"type": "array",
							"description": "An array of tags associated with the text.",
							"items": {
								"type": "string"
							}
						}
					},
					"required": ["summary", "tags"],
					"additionalProperties": false
				}
			}
		},
		"reasoning": {},
		"tools": [],
		"temperature": 1,
		"max_output_tokens": 10000,
		"top_p": 1,
		"store": false
	}`, string(escapedPrompt))

	if s.Debug {
		timestamp := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("debug_%s_payload.json", timestamp)
		err := os.WriteFile(filename, []byte(payload), 0644)
		if err != nil {
			fmt.Printf("Failed to write debug payload to file: %v\n", err)
		}
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return "", nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if s.Debug {
		timestamp := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("debug_%s_body.json", timestamp)
		err := os.WriteFile(filename, body, 0644)
		if err != nil {
			fmt.Printf("Failed to write debug body to file: %v\n", err)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("unexpected status code: %d, body: %v", resp.StatusCode, string(body))
	}

	// Use JSONPath to extract the `text` field
	var jsonData map[string]interface{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		return "", nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	textPath := "$.output[0].content[0].text"
	extractedText, err := jsonpath.JsonPathLookup(jsonData, textPath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to extract text using JSONPath: %w", err)
	}

	// Parse the extracted text as JSON to get `summary` and `tags`
	var parsedOutput struct {
		Summary string   `json:"summary"`
		Tags    []string `json:"tags"`
	}
	if err := json.Unmarshal([]byte(extractedText.(string)), &parsedOutput); err != nil {
		return "", nil, fmt.Errorf("failed to parse extracted text: %w", err)
	}

	return parsedOutput.Summary, parsedOutput.Tags, nil
}
