package frontmatter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

// Frontmatter represents the YAML frontmatter structure
type Frontmatter struct {
	SummarizeAI     string `yaml:"summarize_ai,omitempty"`
	SummarizeAIHash string `yaml:"summarize_ai_hash,omitempty"`
}

// ParseFrontmatter parses the YAML frontmatter from a Markdown file
func ParseFrontmatter(filePath string) (map[string]interface{}, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	frontmatter, err := extractFrontmatter(content)
	if err != nil {
		return nil, fmt.Errorf("failed to extract frontmatter: %w", err)
	}

	var fm map[string]interface{}
	if err := yaml.Unmarshal(frontmatter, &fm); err != nil {
		return nil, fmt.Errorf("failed to unmarshal frontmatter: %w", err)
	}

	return fm, nil
}

// UpdateFrontmatter updates the YAML frontmatter with new summarize_ai and summarize_ai_hash fields
func UpdateFrontmatter(filePath string, summary string, hash string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	frontmatter, err := extractFrontmatter(content)
	if err != nil {
		return fmt.Errorf("failed to extract frontmatter: %w", err)
	}

	var fm map[string]interface{}
	if len(frontmatter) == 0 {
		// Initialize a new frontmatter if none exists
		fm = make(map[string]interface{})
	} else {
		if err := yaml.Unmarshal(frontmatter, &fm); err != nil {
			return fmt.Errorf("failed to unmarshal frontmatter: %w", err)
		}
	}

	fm["summarize_ai"] = summary
	fm["summarize_ai_hash"] = hash

	updatedFrontmatter, err := yaml.Marshal(&fm)
	if err != nil {
		return fmt.Errorf("failed to marshal updated frontmatter: %w", err)
	}

	var updatedContent []byte
	if len(frontmatter) == 0 {
		// Add new frontmatter to the top of the file
		updatedContent = append([]byte("---\n"), updatedFrontmatter...)
		updatedContent = append(updatedContent, []byte("---\n")...)
		updatedContent = append(updatedContent, content...)
	} else {
		// Replace existing frontmatter
		updatedContent = bytes.Replace(content, frontmatter, updatedFrontmatter, 1)
	}

	if err := ioutil.WriteFile(filePath, updatedContent, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write updated file: %w", err)
	}

	return nil
}

// extractFrontmatter extracts the YAML frontmatter from the content
func extractFrontmatter(content []byte) ([]byte, error) {
	start := bytes.Index(content, []byte("---"))
	if start == -1 {
		// Return an empty frontmatter if none is found
		return []byte{}, nil
	}

	end := bytes.Index(content[start+3:], []byte("---"))
	if end == -1 {
		return nil, fmt.Errorf("no frontmatter end delimiter found")
	}

	return content[start+3 : start+3+end], nil
}
