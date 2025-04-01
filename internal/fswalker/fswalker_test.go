package fswalker

import (
	"testing"
)

func TestShouldProcessFile(t *testing.T) {
	tests := []struct {
		name             string
		content          []byte
		always           bool
		includeOutOfDate bool
		expected         bool
	}{
		{
			name:     "Empty file",
			content:  []byte(""),
			always:   false,
			expected: false,
		},
		{
			name:     "Always true, process regardless of content",
			content:  []byte("Some content"),
			always:   true,
			expected: true,
		},
		{
			name:     "No frontmatter, should process",
			content:  []byte("Some content without frontmatter"),
			always:   false,
			expected: true,
		},
		{
			name: "Incomplete frontmatter",
			content: []byte(`---
summarize_ai: "summary"

Content`),
			always:   false,
			expected: true,
		},
		{
			name: "Incomplete frontmatter, missing fields",
			content: []byte(`---
summarize_ai: "summary"
---
Content`),
			always:   false,
			expected: true,
		},
		{
			name: "Complete frontmatter, not outdated",
			content: []byte(`---
summarize_ai: "summary"
summarize_ai_hash: "hash"
summarize_ai_tags: "tags"
summarize_ai_updated: "2023-01-01T12:00"
updated: "2023-01-01T12:00"
---
Content`),
			always:           false,
			includeOutOfDate: false,
			expected:         false,
		},
		{
			name: "Complete frontmatter, outdated",
			content: []byte(`---
summarize_ai: "summary"
summarize_ai_hash: "hash"
summarize_ai_tags: "tags"
summarize_ai_updated: "2023-01-01T12:00"
updated: "2023-01-02T12:00"
---
Content`),
			always:           false,
			includeOutOfDate: true,
			expected:         true,
		},
		{
			name: "Malformed frontmatter, should process",
			content: []byte(`---
summarize_ai: "summary"
summarize_ai_hash: "hash"
summarize_ai_tags: "tags"
summarize_ai_updated: "invalid-date"
updated: "2023-01-02T12:00"
---
Content`),
			always:           false,
			includeOutOfDate: true,
			expected:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldProcessFile(tt.content, tt.always, tt.includeOutOfDate)
			if result != tt.expected {
				t.Errorf("shouldProcessFile() = %v, want %v", result, tt.expected)
			}
		})
	}
}
