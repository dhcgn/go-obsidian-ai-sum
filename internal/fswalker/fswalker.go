package fswalker

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

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
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
				if !override {
					content, err := ioutil.ReadFile(path)
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
				content, err := ioutil.ReadFile(path)
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
