package fswalker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dhcgn/go-obsidian-ai-sum/internal"
	"gopkg.in/yaml.v3"
)

var defaultIgnoreDirs = []string{
	".git",
	".obsidian",
}

func shouldIgnoreDir(name string) bool {
	for _, dir := range defaultIgnoreDirs {
		if name == dir {
			return true
		}
	}
	return false
}

// FileInfo holds information about a markdown file
type FileInfo struct {
	Path           string
	CharacterCount int
}

// shouldProcessFile checks if a file should be processed based on its content and file properties in fontmatter
// always=false and includeOutOfDate=false: process if file is incomplete (ai summary fields are imcomplete)
// always=false and includeOutOfDate=true: process if file is incomplete or the summarize_ai_updated date is older than the updated date in frontmatter
// always=true and includeOutOfDate=false: process the file regardless of its content
// always=true and includeOutOfDate=true: process the file regardless of its content
func shouldProcessFile(content []byte, always, includeOutOfDate bool) bool {
	// Empty files are skipped
	if len(content) == 0 {
		return false
	}

	// If override is true, always process
	if always {
		return true
	}

	// If the file has no frontmatter, it is incomplete and should be processed
	hasFontmatter := strings.HasPrefix(string(content), internal.FontmatterDelimiter)
	if !hasFontmatter {
		return true
	}

	// Extract content between the first and second FontmatterDelimiter
	parts := strings.SplitN(string(content), internal.FontmatterDelimiter, 3)
	if len(parts) < 3 {
		return true // If there's no second delimiter, consider the file incomplete
	}
	frontmatterContent := parts[1]

	// Check for existing summary fields, if any of them are missing the file is incomplete and will be processed
	hasAllProps := strings.Contains(frontmatterContent, internal.FontmatterSummarize_ai+":") &&
		strings.Contains(frontmatterContent, internal.FontmatterSummarize_ai_hash+":") &&
		strings.Contains(frontmatterContent, internal.FontmatterSummarize_ai_tags+":") &&
		strings.Contains(frontmatterContent, internal.FontmatterSummarize_ai_updated+":")
	if !hasAllProps {
		return true
	}

	// At this point, we have a file with an existing summary
	// Only continue if we're checking for outdated files
	if !includeOutOfDate {
		return false
	}

	var frontmatter map[string]interface{}
	if err := yaml.Unmarshal([]byte(frontmatterContent), &frontmatter); err != nil {
		return false
	}

	// For includeOutOfDate mode, we need both date fields
	updatedStr, hasUpdated := frontmatter[internal.FontmatterUpdated].(string)
	summarizeAIUpdatedStr, hasSummarizeUpdated := frontmatter[internal.FontmatterSummarize_ai_updated].(string)

	// If either date field is missing, we consider the file incomplete and process it
	if !hasUpdated || !hasSummarizeUpdated {
		return false
	}

	updated, err1 := time.Parse("2006-01-02T15:04", updatedStr)
	summarizeAIUpdated, err2 := time.Parse("2006-01-02T15:04", summarizeAIUpdatedStr)

	if err1 != nil || err2 != nil {
		return false
	}

	return updated.After(summarizeAIUpdated)
}

// ReadFiles reads a single file or all Markdown files in a folder recursively
func ReadFiles(path string, override bool, onlyOutOfDate bool) ([]FileInfo, error) {
	var files []FileInfo

	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat path: %w", err)
	}

	if info.IsDir() {
		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && shouldIgnoreDir(info.Name()) {
				return filepath.SkipDir
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				if shouldProcessFile(content, override, onlyOutOfDate) {
					files = append(files, FileInfo{
						Path:           path,
						CharacterCount: len(content),
					})
				}
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to walk directory: %w", err)
		}
	} else {
		if strings.HasSuffix(info.Name(), ".md") {
			content, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to read file: %w", err)
			}

			if shouldProcessFile(content, override, onlyOutOfDate) {
				files = append(files, FileInfo{
					Path:           path,
					CharacterCount: len(content),
				})
			}
		}
	}

	return files, nil
}
