package summarizer

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"

	"github.com/dhcgn/go-obsidian-ai-sum/internal/frontmatter"
)

var (
	// DefaultPrompt is the default prompt used for summarization
	//go:embed embed/prompt.md
	defaultPrompt string
)

// LoadPrompt loads the prompt from a flag or uses a default prompt
func LoadPrompt(flagPrompt string) string {
	if flagPrompt != "" {
		return flagPrompt
	}
	return defaultPrompt
}

// ComputeHash computes the hash of the prompt (first 16 hex chars of SHA256)
func ComputeHash(prompt string) string {
	hash := sha256.Sum256([]byte(prompt))
	return hex.EncodeToString(hash[:])[:16]
}

// InjectSummary injects the summary and hash into the YAML frontmatter
func InjectSummary(filePath, summary string, tags []string, hash string) error {
	// This function should call the UpdateFrontmatter function from the frontmatter package
	// to update the YAML frontmatter with the new summary and hash
	return frontmatter.UpdateFrontmatter(filePath, summary, tags, hash)
}
