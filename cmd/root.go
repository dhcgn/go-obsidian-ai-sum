package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	path     string
	apiKey   string
	prompt   string
	override bool
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

		// Add logic to handle the input path and pass it to the internal packages
		// Add logic to read the API key from the environment variable if not provided via flag
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

	rootCmd.MarkPersistentFlagRequired("path")
}

func initConfig() {
	viper.AutomaticEnv()
}
