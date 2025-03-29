package cmd

import (
	"fmt"
	"os"

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

		files, err := fswalker.ReadFiles(path, override)
		if err != nil {
			fmt.Printf("Error reading files: %v\n", err)
			os.Exit(1)
		}

		prompt := summarizer.LoadPrompt(prompt)
		hash := summarizer.ComputeHash(prompt)
		summarizerInstance := summarizer.OpenAISummarizer{
			APIKey: apiKey,
			Debug:  debug,
		}

		for _, file := range files {
			content, err := os.ReadFile(file)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", file, err)
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

	rootCmd.MarkPersistentFlagRequired("path")
}

func initConfig() {
	viper.AutomaticEnv()
}
