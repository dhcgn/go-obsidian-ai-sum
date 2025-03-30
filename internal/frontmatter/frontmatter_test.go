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
			name:           "Empty file with no frontmatter",
			initialContent: ``,
			expectedContent: `---
summarize_ai: Test summary
summarize_ai_hash: TestHash
---`,
			summary: "Test summary",
			hash:    "TestHash",
		},
		{
			name:           "Empty file with no frontmatter",
			initialContent: `Come Content`,
			expectedContent: `---
summarize_ai: Test summary
summarize_ai_hash: TestHash
---
Come Content`,
			summary: "Test summary",
			hash:    "TestHash",
		},
		{
			name: "Empty file with no frontmatter - only one line break",
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
			name: "File with frontmatter and old summarize_ai keys",
			initialContent: `---
existing_key: existing_value
summarize_ai: Test summary MUST OVERRIDE
summarize_ai_hash: TestHash MUST OVERRIDE
existing_key_lower: existing_value
---
Content below frontmatter.
`,
			expectedContent: `---
existing_key: existing_value
summarize_ai: Test summary
summarize_ai_hash: TestHash
existing_key_lower: existing_value
---
Content below frontmatter.
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
			name: "File with frontmatter but no summarize_ai keys and empty line",
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
			if err := UpdateFrontmatter(tmpFile.Name(), tt.summary, nil, tt.hash); err != nil {
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

func TestUpdateFrontmatterWithTags(t *testing.T) {
	tests := []struct {
		name            string
		initialContent  string
		expectedContent string
		summary         string
		tags            []string
		hash            string
	}{
		{
			name:           "New frontmatter with single tag",
			initialContent: "Some content",
			expectedContent: `---
summarize_ai: Test summary
summarize_ai_hash: TestHash
summarize_ai_tags:
  - tag1
---
Some content`,
			summary: "Test summary",
			tags:    []string{"tag1"},
			hash:    "TestHash",
		},
		{
			name:           "New frontmatter with multiple tags",
			initialContent: "Some content",
			expectedContent: `---
summarize_ai: Test summary
summarize_ai_hash: TestHash
summarize_ai_tags:
  - tag1
  - tag2
  - tag3
---
Some content`,
			summary: "Test summary",
			tags:    []string{"tag1", "tag2", "tag3"},
			hash:    "TestHash",
		},
		{
			name: "Update existing frontmatter with tags",
			initialContent: `---
title: Example
summarize_ai: Old summary
summarize_ai_hash: OldHash
summarize_ai_tags:
  - oldtag
---
Content here
`,
			expectedContent: `---
title: Example
summarize_ai: Test summary
summarize_ai_hash: TestHash
summarize_ai_tags:
  - newtag1
  - newtag2
---
Content here
`,
			summary: "Test summary",
			tags:    []string{"newtag1", "newtag2"},
			hash:    "TestHash",
		},
		{
			name: "Update with empty tags should not include tags field",
			initialContent: `---
title: Example
summarize_ai: Old summary
summarize_ai_hash: OldHash
summarize_ai_tags:
  - oldtag
---
Content here
`,
			expectedContent: `---
title: Example
summarize_ai: Test summary
summarize_ai_hash: TestHash
---
Content here
`,
			summary: "Test summary",
			tags:    nil,
			hash:    "TestHash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "testfile-*.md")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.Write([]byte(tt.initialContent)); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpFile.Close()

			if err := UpdateFrontmatter(tmpFile.Name(), tt.summary, tt.tags, tt.hash); err != nil {
				t.Fatalf("UpdateFrontmatter failed: %v", err)
			}

			updatedContent, err := os.ReadFile(tmpFile.Name())
			if err != nil {
				t.Fatalf("Failed to read updated file: %v", err)
			}

			if string(updatedContent) != tt.expectedContent {
				t.Errorf("Content mismatch.\nExpected:\n'%s'\nGot:\n'%s'", tt.expectedContent, string(updatedContent))
			}
		})
	}
}
