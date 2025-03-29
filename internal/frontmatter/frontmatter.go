package frontmatter

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

// Frontmatter represents the YAML frontmatter structure
type Frontmatter struct {
	Title          string   `yaml:"title"`
	Year           int      `yaml:"year"`
	Favorite       bool     `yaml:"favorite"`
	SummarizeAI    string   `yaml:"summarize_ai,omitempty"`
	SummarizeAIHash string  `yaml:"summarize_ai_hash,omitempty"`
	Cast           []string `yaml:"cast"`
}

// ParseFrontmatter parses the YAML frontmatter from a Markdown file
func ParseFrontmatter(filePath string) (*Frontmatter, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	frontmatter, err := extractFrontmatter(content)
	if err != nil {
		return nil, fmt.Errorf("failed to extract frontmatter: %w", err)
	}

	var fm Frontmatter
	if err := yaml.Unmarshal(frontmatter, &fm); err != nil {
		return nil, fmt.Errorf("failed to unmarshal frontmatter: %w", err)
	}

	return &fm, nil
}

// UpdateFrontmatter updates the YAML frontmatter with new summarize_ai and summarize_ai_hash fields
func UpdateFrontmatter(filePath string, summary string, hash string) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	frontmatter, err := extractFrontmatter(content)
	if err != nil {
		return fmt.Errorf("failed to extract frontmatter: %w", err)
	}

	var fm Frontmatter
	if err := yaml.Unmarshal(frontmatter, &fm); err != nil {
		return fmt.Errorf("failed to unmarshal frontmatter: %w", err)
	}

	fm.SummarizeAI = summary
	fm.SummarizeAIHash = hash

	updatedFrontmatter, err := yaml.Marshal(&fm)
	if err != nil {
		return fmt.Errorf("failed to marshal updated frontmatter: %w", err)
	}

	updatedContent := bytes.Replace(content, frontmatter, updatedFrontmatter, 1)
	if err := ioutil.WriteFile(filePath, updatedContent, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write updated file: %w", err)
	}

	return nil
}

// extractFrontmatter extracts the YAML frontmatter from the content
func extractFrontmatter(content []byte) ([]byte, error) {
	start := bytes.Index(content, []byte("---"))
	if start == -1 {
		return nil, fmt.Errorf("no frontmatter start delimiter found")
	}

	end := bytes.Index(content[start+3:], []byte("---"))
	if end == -1 {
		return nil, fmt.Errorf("no frontmatter end delimiter found")
	}

	return content[start+3 : start+3+end], nil
}
