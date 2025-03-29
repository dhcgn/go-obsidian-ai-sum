package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/dhcgn/go-obsidian-ai-sum/internal/fswalker"
	"github.com/dhcgn/go-obsidian-ai-sum/internal/summarizer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	path     string
	apiKey   string
	prompt   string
	override bool
	debug    bool
	dryrun   bool
)

const (
	LimitChars = 10_000
)

var rootCmd = &cobra.Command{
	Use:   "go-obsidian-ai-sum",
	Short: "Summarize Obsidian Markdown pages using AI",
	Run: func(cmd *cobra.Command, args []string) {
		if apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
			if apiKey == "" {
				fmt.Println("API key is required. Provide it via --api-key flag or OPENAI_API_KEY environment variable.")
				os.Exit(1)
			}
		}

		start := time.Now()
		files, err := fswalker.ReadFiles(path, override)
		fmt.Println("Reading files took:", time.Since(start))
		if err != nil {
			fmt.Printf("Error reading files: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Found", len(files), "files to summarize")

		prompt := summarizer.LoadPrompt(prompt)
		hash := summarizer.ComputeHash(prompt)
		fmt.Println("Prompt template hash:", hash)
		summarizerInstance := summarizer.OpenAISummarizer{
			APIKey: apiKey,
			Debug:  debug,
		}

		if dryrun {
			fmt.Println("Dry run mode - no API calls will be made.")
		}

		start = time.Now()
		for _, file := range files {
			content, err := os.ReadFile(file)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", file, err)
				continue
			}

			// only the first 10.000 characters are sent to the API
			if len(content) > LimitChars {
				content = content[:LimitChars]
			}

			if dryrun {
				continue
			}

			summary, _, err := summarizerInstance.Summarize(string(content), file, prompt)
			if err != nil {
				fmt.Printf("Error summarizing file %s: %v\n", file, err)
				continue
			}

			err = summarizer.InjectSummary(file, summary, hash)
			if err != nil {
				fmt.Printf("Error injecting summary into file %s: %v\n", file, err)
			}
		}
		fmt.Println("Summarization took:", time.Since(start))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&path, "path", "", "Path to file or folder")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key for the AI provider")
	rootCmd.PersistentFlags().StringVar(&prompt, "prompt", "", "Custom prompt for summarization")
	rootCmd.PersistentFlags().BoolVar(&override, "override", false, "Override existing summaries")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode to log payloads")
	rootCmd.PersistentFlags().BoolVar(&dryrun, "dryrun", false, "Dry run mode - stops before making API calls")

	rootCmd.MarkPersistentFlagRequired("path")
}

func initConfig() {
	viper.AutomaticEnv()
}
