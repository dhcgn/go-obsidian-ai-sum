# go-obsidian-ai-sum

**go-obsidian-ai-sum** is a CLI tool that leverages AI to generate concise summaries for your Obsidian Markdown files and injects them directly into the YAML frontmatter. It can process individual files or entire directories recursively, making it perfect for automating vault maintenance.

## Features

- **Recursive Processing:** Scans a single file or entire folders (with subfolders) for Markdown files.
- **AI-Generated Summaries:** Uses OpenAI (default) to produce precise and concise summaries.
- **Frontmatter Injection:** Automatically adds/updates `summarize_ai`, `summarize_ai_hash`, and `summarize_ai_tags` fields.
- **Custom Prompt Support:** Allows the use of a custom prompt to tailor the summarization.
- **Override Existing Summaries:** Optionally overwrite previously generated summaries.
- **Cost Estimation:** Provides a rough estimate of API costs based on content length.
- **Dry Run Mode:** Simulate the summarization process without making any API calls.
- **Random File Order:** Option to process files in a random order.
- **Pluggable AI Provider:** Designed for future support of alternative AI providers.

## What It Does

- **Reads** Markdown files (or entire folders).
- **Sends** file content along with a prompt to an AI service.
- **Receives** a summary and associated tags.
- **Injects** the results into the file's YAML frontmatter—automatically updating or creating the necessary keys.

## Warning

> **WARNING:** This tool will modify your Markdown files directly. Ensure you have backups or work with copies. Use at your own risk—there is no warranty for any changes made.

## Installation

Install the latest version via:

```bash
go install github.com/dhcgn/go-obsidian-ai-sum@latest
```

## Usage

```bash
go-obsidian-ai-sum --path ./vault --api-key YOUR_API_KEY [--override] [--prompt "Custom prompt"]
```

### CLI Flags

| Flag                  | Description                                                    |
|-----------------------|----------------------------------------------------------------|
| `--path`              | Path to a Markdown file or folder                              |
| `--api-key`           | API key for the AI provider (or use `OPENAI_API_KEY` environment variable) |
| `--override`          | Overwrite existing summaries                                   |
| `--prompt`            | Custom prompt for summarization                                |
| `--dryrun`            | Run in simulation mode (no API calls)                          |
| `--random-file-access`| Process files in a random order (optional)                     |

## Roadmap

- Support additional AI providers (e.g., Claude, Mistral)
- Multilingual summarization and language detection
- Enhanced customization for frontmatter updates

