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

	var node yaml.Node
	if len(frontmatter) == 0 {
		// Initialize new frontmatter
		node.Kind = yaml.MappingNode
		node.Style = yaml.TaggedStyle
	} else {
		if err := yaml.Unmarshal(frontmatter, &node); err != nil {
			return fmt.Errorf("failed to unmarshal frontmatter: %w", err)
		}
	}

	// Find or create the mapping node
	var mappingNode *yaml.Node
	if node.Kind == yaml.DocumentNode {
		mappingNode = node.Content[0]
	} else {
		mappingNode = &node
	}

	// Update or add summarize_ai and summarize_ai_hash
	addOrUpdateYAMLField(mappingNode, "summarize_ai", summary)
	addOrUpdateYAMLField(mappingNode, "summarize_ai_hash", hash)

	updatedFrontmatter, err := yaml.Marshal(&node)
	if err != nil {
		return fmt.Errorf("failed to marshal updated frontmatter: %w", err)
	}

	var updatedContent []byte
	if len(frontmatter) == 0 {
		updatedContent = append([]byte("---\n"), updatedFrontmatter...)
		updatedContent = append(updatedContent, []byte("---")...)
		if len(bytes.TrimSpace(content)) > 0 {
			updatedContent = append(updatedContent, '\n')
			updatedContent = append(updatedContent, bytes.TrimLeft(content, "\n")...)
		}
	} else {
		updatedContent = append([]byte("---\n"), updatedFrontmatter...)
		updatedContent = append(updatedContent, []byte("---")...)

		remainingContent := content[bytes.Index(content, []byte("---"))+len(frontmatter)+6:]
		if len(bytes.TrimSpace(remainingContent)) > 0 {
			updatedContent = append(updatedContent, '\n')
			updatedContent = append(updatedContent, bytes.TrimLeft(remainingContent, "\n")...)
		}
	}

	// Ensure exactly one newline at the end if needed
	if len(bytes.TrimSpace(updatedContent)) > 0 {
		updatedContent = bytes.TrimRight(updatedContent, "\n")
		updatedContent = append(updatedContent, '\n')
	}

	if err := ioutil.WriteFile(filePath, updatedContent, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write updated file: %w", err)
	}

	return nil
}

// addOrUpdateYAMLField adds or updates a field in the YAML node while preserving order
func addOrUpdateYAMLField(node *yaml.Node, key, value string) {
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			node.Content[i+1].Value = value
			return
		}
	}
	node.Content = append(node.Content,
		&yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: key,
		},
		&yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: value,
		},
	)
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
