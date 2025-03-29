package summarizer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Summarizer is an interface for summarizing text
type Summarizer interface {
	Summarize(text, prompt string) (string, error)
}

// OpenAISummarizer is an implementation of Summarizer using OpenAI
type OpenAISummarizer struct {
	APIKey string
	Debug  bool
}

// Summarize generates a summary using the OpenAI API
func (s *OpenAISummarizer) Summarize(text, prompt string) (string, error) {
	url := "https://api.openai.com/v1/responses"

	escapedPrompt, err := json.Marshal(prompt + "\n\nText:\n" + text)
	if err != nil {
		return "", fmt.Errorf("failed to escape prompt to JSON: %w", err)
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
		filename := fmt.Sprintf("payload_%s.json", timestamp)
		err := os.WriteFile(filename, []byte(payload), 0644)
		if err != nil {
			fmt.Printf("Failed to write debug payload to file: %v\n", err)
		}
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d, body: %v", resp.StatusCode, string(body))
	}

	// Extract the summary from the response body (assuming JSON format)
	summary := extractSummaryFromResponse(body)
	return summary, nil
}

// extractSummaryFromResponse extracts the summary from the API response body
func extractSummaryFromResponse(body []byte) string {
	// This is a simplified example; you may need to adjust based on the actual API response format
	return string(body)
}
