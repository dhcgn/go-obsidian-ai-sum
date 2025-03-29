package summarizer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Summarizer is an interface for summarizing text
type Summarizer interface {
	Summarize(text, prompt string) (string, error)
}

// OpenAISummarizer is an implementation of Summarizer using OpenAI
type OpenAISummarizer struct {
	APIKey string
}

// Summarize generates a summary using the OpenAI API
func (s *OpenAISummarizer) Summarize(text, prompt string) (string, error) {
	url := "https://api.openai.com/v1/engines/davinci-codex/completions"
	payload := fmt.Sprintf(`{"prompt": "%s\n\n%s", "max_tokens": 100}`, prompt, text)
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

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Extract the summary from the response body (assuming JSON format)
	// This is a simplified example; you may need to adjust based on the actual API response format
	summary := extractSummaryFromResponse(body)
	return summary, nil
}

// extractSummaryFromResponse extracts the summary from the API response body
func extractSummaryFromResponse(body []byte) string {
	// This is a simplified example; you may need to adjust based on the actual API response format
	return string(body)
}
