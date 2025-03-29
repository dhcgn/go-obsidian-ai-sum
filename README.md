# go-obsidian-ai-sum

Summarize Obsidian Markdown pages using AI and store the result as a YAML file property.

## Features

- 📁 Accepts a file or folder (recursively processes subfolders)
- 🧠 Uses AI to generate summaries (default: OpenAI)
- ✍️ Stores the summary in the frontmatter under `summarize_ai`
- 🔁 Optionally overrides existing summaries
- 🧾 Adds a `summarize_ai_hash` for tracking the prompt used
- 🔧 Custom prompt support
- 🔐 API key can be provided via environment variable or CLI flag
- 🔄 AI provider is pluggable (planned for future)

## Example YAML Frontmatter

```yaml
---
title: A New Hope
year: 1977
favorite: true
summarize_ai: Summarized Text with AI
summarize_ai_prompt_hash: 3f5a9c1b7e2d1a0f
cast:
  - Mark Hamill
  - Harrison Ford
  - Carrie Fisher
---
```

## Installation

```bash
go install github.com/dhcgn/go-obsidian-ai-sum@latest
```

## Usage

```bash
go-obsidian-ai-sum \
  --path ./vault \
  --api-key sk-xxx \
  --override \
  --prompt "Summarize this file for sci-fi fans"
```

### CLI Flags

| Flag          | Description                                      |
|---------------|--------------------------------------------------|
| `--path`      | File or folder to process                        |
| `--api-key`   | (Optional) API key for the AI provider           |
| `--override`  | If set, overwrites existing summaries            |
| `--prompt`    | (Optional) Custom prompt for summarization       |

If `--api-key` is not provided, it will use the `OPENAI_API_KEY` environment variable.

## Roadmap

- [ ] Add support for multiple AI providers (e.g., Claude, Mistral)
- [ ] Add language detection and multilingual summarization
