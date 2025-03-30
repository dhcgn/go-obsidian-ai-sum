[![Go](https://github.com/dhcgn/go-obsidian-ai-sum/actions/workflows/go.yml/badge.svg)](https://github.com/dhcgn/go-obsidian-ai-sum/actions/workflows/go.yml)
# go-obsidian-ai-sum

**go-obsidian-ai-sum** is a CLI tool that leverages AI to generate concise summaries for your Obsidian Markdown files and injects them directly into the YAML frontmatter. It can process individual files or entire directories recursively, making it perfect for automating vault maintenance.

## Motivation

This tool was created to optimize note management using AI capabilities. By automatically generating summaries for each note, it enables powerful ways to organize and understand your knowledge base:

- **Smart Organization:** Use AI-generated summaries to automatically categorize and rearrange notes based on their content
- **Quick Overview:** Get instant context about any note without opening it
- **Enhanced Search:** Leverage summaries to find related content more effectively
- **Future Extensibility:** Potential for an Obsidian plugin to provide these features directly in the editor

The summaries are stored in the frontmatter, making them easily accessible for other tools and automations.

## Why Go?

Go was chosen as the implementation language for several key reasons:

- **Concurrency:** Efficient handling of parallel API calls through goroutines, essential for processing multiple files
- **Cross-Platform:** Easy compilation to different operating systems without dependencies
- **Fast Development:** Quick prototyping and implementation of CLI tools
- **Performance:** Native execution speed for file operations and text processing
- **Simple Deployment:** Single binary distribution without runtime dependencies

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

## ⚠️⚠️⚠️ Warning ⚠️⚠️⚠️

> **WARNING:** This tool will modify your Markdown files directly. Ensure you have backups or work with copies. Use at your own risk—there is no warranty for any changes made.

> **PRIVACY NOTICE:** This proof-of-concept currently uses OpenAI's API endpoints. When using this tool, your note content will be transmitted to OpenAI's servers for processing. If you have sensitive or private information, please be aware of this data transfer. Future versions will support local LLMs for enhanced privacy.

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

| Flag                   | Description                                                                |
| ---------------------- | -------------------------------------------------------------------------- |
| `--path`               | Path to a Markdown file or folder                                          |
| `--api-key`            | API key for the AI provider (or use `OPENAI_API_KEY` environment variable) |
| `--override`           | Overwrite existing summaries                                               |
| `--prompt`             | Custom prompt for summarization                                            |
| `--dryrun`             | Run in simulation mode (no API calls)                                      |
| `--random-file-access` | Process files in a random order (optional)                                 |
| `--top`                | Process only this many files (0 for all)                                   |

## Roadmap

- Support additional AI providers (e.g., Claude, Mistral)
- Multilingual summarization and language detection
- Enhanced customization for frontmatter updates

