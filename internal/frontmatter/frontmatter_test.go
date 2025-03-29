package frontmatter

import (
	"os"
	"testing"
)

func TestUpdateFrontmatter(t *testing.T) {
	tests := []struct {
		name            string
		initialContent  string
		expectedContent string
		summary         string
		hash            string
	}{
		{
			name:            "Empty file with no frontmatter",
			initialContent:  "",
			expectedContent: "---\nsummarize_ai: Test summary\nsummarize_ai_hash: TestHash\n---\n",
			summary:         "Test summary",
			hash:            "TestHash",
		},
		{
			name:            "File with frontmatter but no summarize_ai keys",
			initialContent:  "---\nexisting_key: existing_value\n---\nContent below frontmatter.",
			expectedContent: "---\nexisting_key: existing_value\nsummarize_ai: Test summary\nsummarize_ai_hash: TestHash\n---\nContent below frontmatter.",
			summary:         "Test summary",
			hash:            "TestHash",
		},
		{
			name:            "File with frontmatter containing unrelated keys",
			initialContent:  "---\nunrelated_key: unrelated_value\n---\nContent below frontmatter.",
			expectedContent: "---\nunrelated_key: unrelated_value\nsummarize_ai: Test summary\nsummarize_ai_hash: TestHash\n---\nContent below frontmatter.",
			summary:         "Test summary",
			hash:            "TestHash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file
			tmpFile, err := os.CreateTemp("", "testfile-*.md")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			// Write initial content to the file
			if _, err := tmpFile.Write([]byte(tt.initialContent)); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpFile.Close()

			// Call UpdateFrontmatter
			if err := UpdateFrontmatter(tmpFile.Name(), tt.summary, tt.hash); err != nil {
				t.Fatalf("UpdateFrontmatter failed: %v", err)
			}

			// Read the updated content
			updatedContent, err := os.ReadFile(tmpFile.Name())
			if err != nil {
				t.Fatalf("Failed to read updated file: %v", err)
			}

			// Compare the updated content with the expected content
			if string(updatedContent) != tt.expectedContent {
				t.Errorf("Content mismatch.\nExpected:\n%s\nGot:\n%s", tt.expectedContent, string(updatedContent))
			}
		})
	}
}
