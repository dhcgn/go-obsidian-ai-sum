package fswalker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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

func shouldProcessFile(content []byte, override, onlyOutOfDate bool) bool {
	// Empty files are skipped
	if len(content) == 0 {
		return false
	}

	// If override is true, always process
	if override {
		return true
	}

	// Check for existing summary
	hasSummary := strings.Contains(string(content), "summarize_ai:")
	if !hasSummary {
		return true
	}

	// At this point, we have a file with an existing summary
	// Only continue if we're checking for outdated files
	if !onlyOutOfDate {
		return false
	}

	var frontmatter map[string]interface{}
	if err := yaml.Unmarshal(content, &frontmatter); err != nil {
		return false
	}

	// For onlyOutOfDate mode, we need both date fields
	updatedStr, hasUpdated := frontmatter["updated"].(string)
	summarizeAIUpdatedStr, hasSummarizeUpdated := frontmatter["summarize_ai_updated"].(string)

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
