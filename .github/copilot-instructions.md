The tool reads Obsidian Markdown files and summarizes their content using an AI provider (default: OpenAI). The summaries are stored in the YAML frontmatter of each file under the property `summarize_ai`, along with a hash of the prompt used (`summarize_ai_hash`).

This project is written in Go and should follow idiomatic Go practices.

The tool supports both files and folders as input. When a folder is passed, files are processed recursively.

It accepts optional arguments to override existing summaries, customize the summarization prompt, or provide the API key (also read from the environment variable `OPENAI_API_KEY`).

Copilot should assume the YAML format is manually written and must be preserved (no corruption of structure or indentation).

The project is designed to support pluggable AI providers in the future.

When writing code, prefer clean, testable Go functions. Keep dependencies minimal.

When discussing Obsidian-specific behavior, assume the user is familiar with frontmatter and markdown file organization in Obsidian vaults.
