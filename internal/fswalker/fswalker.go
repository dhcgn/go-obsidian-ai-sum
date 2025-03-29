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

// ReadFiles reads a single file or all Markdown files in a folder recursively
func ReadFiles(path string, override bool) ([]string, error) {
	var files []string

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
				if !override {
					content, err := os.ReadFile(path)
					if err != nil {
						return err
					}
					if strings.Contains(string(content), "summarize_ai:") {
						return nil
					}
				}
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to walk directory: %w", err)
		}
	} else {
		if strings.HasSuffix(info.Name(), ".md") {
			if !override {
				content, err := os.ReadFile(path)
				if err != nil {
					return nil, fmt.Errorf("failed to read file: %w", err)
				}
				if strings.Contains(string(content), "summarize_ai:") {
					return nil, nil
				}
			}
			files = append(files, path)
		}
	}

	return files, nil
}
