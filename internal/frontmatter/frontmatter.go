package frontmatter

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	delimiter = []byte("---")
)

// UpdateFrontmatter updates the YAML frontmatter with new summarize_ai and summarize_ai_hash fields.
func UpdateFrontmatter(filePath string, summary string, hash string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	lines := bytes.Split(content, []byte("\n"))

	// If no frontmatter exists, prepend one.
	if !bytes.Equal(bytes.TrimSpace(lines[0]), delimiter) {
		newFM := fmt.Sprintf("---\nsummarize_ai: %s\nsummarize_ai_hash: %s\n---\n", summary, hash)
		newContent := newFM + string(content)
		return os.WriteFile(filePath, []byte(newContent), os.ModePerm)
	}

	// Locate the closing delimiter.
	closingIndex := -1
	for i := 1; i < len(lines); i++ {
		if bytes.Equal(bytes.TrimSpace(lines[i]), delimiter) {
			closingIndex = i
			break
		}
	}
	if closingIndex == -1 {
		return fmt.Errorf("no closing frontmatter delimiter found")
	}

	// Extract frontmatter content (lines between the delimiters).
	fmLines := lines[1:closingIndex]
	fmContent := bytes.Join(fmLines, []byte("\n"))

	// Use yaml.Node to preserve key order.
	var node yaml.Node
	if len(bytes.TrimSpace(fmContent)) == 0 {
		node.Kind = yaml.MappingNode
	} else {
		if err := yaml.Unmarshal(fmContent, &node); err != nil {
			return fmt.Errorf("failed to unmarshal frontmatter: %w", err)
		}
	}

	var mappingNode *yaml.Node
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		mappingNode = node.Content[0]
	} else {
		mappingNode = &node
	}

	// Update or add the fields.
	addOrUpdateYAMLField(mappingNode, "summarize_ai", summary)
	addOrUpdateYAMLField(mappingNode, "summarize_ai_hash", hash)

	updatedFM, err := yaml.Marshal(&node)
	if err != nil {
		return fmt.Errorf("failed to marshal updated frontmatter: %w", err)
	}

	// Reassemble the file.
	newLines := [][]byte{}
	newLines = append(newLines, delimiter)
	fmUpdatedLines := bytes.Split(bytes.TrimSuffix(updatedFM, []byte("\n")), []byte("\n"))
	newLines = append(newLines, fmUpdatedLines...)
	newLines = append(newLines, delimiter)
	restLines := lines[closingIndex+1:]
	newContent := bytes.Join(newLines, []byte("\n"))
	if len(restLines) > 0 {
		newContent = append(newContent, '\n')
		newContent = append(newContent, bytes.Join(restLines, []byte("\n"))...)
	}
	newContent = bytes.TrimRight(newContent, "\n")
	newContent = append(newContent, '\n')

	return os.WriteFile(filePath, newContent, os.ModePerm)
}

// addOrUpdateYAMLField adds or updates a field in the YAML node while preserving order.
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
