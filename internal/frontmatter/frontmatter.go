package frontmatter

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// UpdateFrontmatter updates (or creates) only the summarize_ai, summarize_ai_hash,
// summarize_ai_tags, and summarize_ai_updated keys in the frontmatter, leaving all other text content untouched.
func UpdateFrontmatter(filePath, summary string, tags []string, hash string) error {
	// Read the original file content.
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	content := string(contentBytes)

	// Get the current timestamp in the desired format.
	currentTimestamp := time.Now().Format("2006-01-02T15:04")

	// Check if file starts with a frontmatter block.
	if strings.HasPrefix(content, "---") {
		// Split content into lines.
		lines := strings.Split(content, "\n")
		// Find the closing delimiter line (the second occurrence of a line that equals '---').
		closingIndex := -1
		for i := 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "---" {
				closingIndex = i
				break
			}
		}
		if closingIndex == -1 {
			return fmt.Errorf("no closing frontmatter delimiter found")
		}

		// Process the existing frontmatter (lines[1:closingIndex]).
		frontLines := lines[1:closingIndex]
		var newFront []string
		updatedAI, updatedHash, updatedTags, updatedTimestamp := false, false, false, false

		for i := 0; i < len(frontLines); i++ {
			line := frontLines[i]
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "summarize_ai:") && !strings.HasPrefix(trimmed, "summarize_ai_hash:") && !strings.HasPrefix(trimmed, "summarize_ai_updated:") {
				newFront = append(newFront, fmt.Sprintf("summarize_ai: \"%s\"", summary))
				updatedAI = true
			} else if strings.HasPrefix(trimmed, "summarize_ai_hash:") {
				newFront = append(newFront, "summarize_ai_hash: "+hash)
				updatedHash = true
			} else if strings.HasPrefix(trimmed, "summarize_ai_tags:") {
				// Skip this line and the following indented lines.
				if len(tags) > 0 {
					newFront = append(newFront, "summarize_ai_tags:")
					for _, tag := range tags {
						newFront = append(newFront, "  - "+tag)
					}
				}
				updatedTags = true
				// Skip all subsequent indented lines.
				j := i + 1
				for j < len(frontLines) && frontLines[j] != "" && (frontLines[j][0] == ' ' || frontLines[j][0] == '\t') {
					j++
				}
				i = j - 1 // Adjust loop index.
			} else if strings.HasPrefix(trimmed, "summarize_ai_updated:") {
				newFront = append(newFront, "summarize_ai_updated: "+currentTimestamp)
				updatedTimestamp = true
			} else {
				newFront = append(newFront, line)
			}
		}

		// If any key is missing, add it.
		if !updatedAI {
			newFront = append(newFront, fmt.Sprintf("summarize_ai: \"%s\"", summary))
		}
		if !updatedHash {
			newFront = append(newFront, "summarize_ai_hash: "+hash)
		}
		if !updatedTags && len(tags) > 0 {
			newFront = append(newFront, "summarize_ai_tags:")
			for _, tag := range tags {
				newFront = append(newFront, "  - "+tag)
			}
		}
		if !updatedTimestamp {
			newFront = append(newFront, "summarize_ai_updated: "+currentTimestamp)
		}

		// Reassemble the frontmatter block exactly.
		newFrontmatter := "---\n" + strings.Join(newFront, "\n") + "\n---"

		// Reassemble the final file content:
		// Keep the remainder (after the frontmatter) exactly as-is.
		remainder := ""
		if closingIndex+1 < len(lines) {
			remainder = strings.Join(lines[closingIndex+1:], "\n")
		}
		finalContent := ""
		if remainder != "" {
			finalContent = newFrontmatter + "\n" + remainder
		} else {
			finalContent = newFrontmatter
		}
		// Do not add any extra newline at the end.
		return os.WriteFile(filePath, []byte(finalContent), os.ModePerm)
	}

	// No frontmatter exists: create a new frontmatter block and prepend it.
	var front []string
	front = append(front, fmt.Sprintf("summarize_ai: \"%s\"", summary))
	front = append(front, "summarize_ai_hash: "+hash)
	if len(tags) > 0 {
		front = append(front, "summarize_ai_tags:")
		for _, tag := range tags {
			front = append(front, "  - "+tag)
		}
	}
	front = append(front, "summarize_ai_updated: "+currentTimestamp)
	newFrontmatter := "---\n" + strings.Join(front, "\n") + "\n---"
	finalContent := ""
	if content != "" {
		finalContent = newFrontmatter + "\n" + content
	} else {
		finalContent = newFrontmatter
	}
	return os.WriteFile(filePath, []byte(finalContent), os.ModePerm)
}
