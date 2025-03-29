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
			name: "Empty file with no frontmatter",
			initialContent: `
`,
			expectedContent: `---
summarize_ai: Test summary
summarize_ai_hash: TestHash
---
`,
			summary: "Test summary",
			hash:    "TestHash",
		},
		{
			name: "File with frontmatter but no summarize_ai keys",
			initialContent: `---
existing_key: existing_value
---
Content below frontmatter.
`,
			expectedContent: `---
existing_key: existing_value
summarize_ai: Test summary
summarize_ai_hash: TestHash
---
Content below frontmatter.
`,
			summary: "Test summary",
			hash:    "TestHash",
		},
		{
			name: "File with frontmatter containing unrelated keys",
			initialContent: `---
unrelated_key: unrelated_value
---
Content below frontmatter.
`,
			expectedContent: `---
unrelated_key: unrelated_value
summarize_ai: Test summary
summarize_ai_hash: TestHash
---
Content below frontmatter.
`,
			summary: "Test summary",
			hash:    "TestHash",
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
				// Find the first differing character
				minLen := min(len(updatedContent), len(tt.expectedContent))
				var diffIndex int
				for diffIndex = 0; diffIndex < minLen-1; diffIndex++ {
					if updatedContent[diffIndex] != tt.expectedContent[diffIndex] {
						break
					}
				}
				t.Errorf("Content mismatch.\nExpected:\n%s\nGot:\n%s\nFirst difference at index %d: expected '%c' (0x%x), got '%c' (0x%x)",
					tt.expectedContent, string(updatedContent), diffIndex, tt.expectedContent[diffIndex], tt.expectedContent[diffIndex],
					updatedContent[diffIndex], updatedContent[diffIndex])
			}
		})
	}
}
