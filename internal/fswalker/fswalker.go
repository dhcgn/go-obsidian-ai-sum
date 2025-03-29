package fswalker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// ReadFiles reads a single file or all Markdown files in a folder recursively
func ReadFiles(path string, override bool) ([]FileInfo, error) {
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

				if len(content) == 0 {
					return nil
				}

				if !override && strings.Contains(string(content), "summarize_ai:") {
					return nil
				}

				files = append(files, FileInfo{
					Path:           path,
					CharacterCount: len(content),
				})
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

			if !override && strings.Contains(string(content), "summarize_ai:") {
				return nil, nil
			}

			files = append(files, FileInfo{
				Path:           path,
				CharacterCount: len(content),
			})
		}
	}

	return files, nil
}
