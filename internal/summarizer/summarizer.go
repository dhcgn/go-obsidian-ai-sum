package summarizer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// LoadPrompt loads the prompt from a flag or uses a default prompt
func LoadPrompt(flagPrompt string) string {
	if flagPrompt != "" {
		return flagPrompt
	}
	return "Summarize the following text:"
}

// ComputeHash computes the hash of the prompt (first 16 hex chars of SHA256)
func ComputeHash(prompt string) string {
	hash := sha256.Sum256([]byte(prompt))
	return hex.EncodeToString(hash[:])[:16]
}

// GenerateSummary uses the OpenAI API to generate a summary
func GenerateSummary(apiKey, text, prompt string) (string, error) {
	url := "https://api.openai.com/v1/engines/davinci-codex/completions"
	payload := fmt.Sprintf(`{"prompt": "%s\n\n%s", "max_tokens": 100}`, prompt, text)
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

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

// InjectSummary injects the summary and hash into the YAML frontmatter
func InjectSummary(filePath, summary, hash string) error {
	// This function should call the UpdateFrontmatter function from the frontmatter package
	// to update the YAML frontmatter with the new summary and hash
	return frontmatter.UpdateFrontmatter(filePath, summary, hash)
}

// extractSummaryFromResponse extracts the summary from the API response body
func extractSummaryFromResponse(body []byte) string {
	// This is a simplified example; you may need to adjust based on the actual API response format
	return string(body)
}
