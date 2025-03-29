package cmd

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dhcgn/go-obsidian-ai-sum/internal/fswalker"
	"github.com/dhcgn/go-obsidian-ai-sum/internal/summarizer"
	"github.com/pterm/pterm"
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
	LimitChars = 50_000
)

var rootCmd = &cobra.Command{
	Use:   "go-obsidian-ai-sum",
	Short: "Summarize Obsidian Markdown pages using AI",
	Run: func(cmd *cobra.Command, args []string) {
		if apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
			if apiKey == "" {
				pterm.Error.Println("API key is required. Provide it via --api-key flag or OPENAI_API_KEY environment variable.")
				os.Exit(1)
			}
		}

		start := time.Now()
		files, err := fswalker.ReadFiles(path, override)
		pterm.Info.Printf("Reading files took: %v\n", time.Since(start))
		if err != nil {
			pterm.Error.Printf("Error reading files: %v\n", err)
			os.Exit(1)
		}

		pterm.Info.Printf("Found %d files to summarize\n", len(files))

		prompt := summarizer.LoadPrompt(prompt)
		hash := summarizer.ComputeHash(prompt)
		pterm.Info.Printf("Prompt template hash: %s\n", hash)
		summarizerInstance := summarizer.OpenAISummarizer{
			APIKey: apiKey,
			Debug:  debug,
		}

		if dryrun {
			pterm.Warning.Println("Dry run mode - no API calls will be made.")
		}

		start = time.Now()

		const workerCount = 10
		var wg sync.WaitGroup
		jobChan := make(chan string, len(files))
		errChan := make(chan error, len(files))

		// Create progress bar
		progress, _ := pterm.DefaultProgressbar.
			WithTotal(len(files)).
			WithTitle("Summarizing files").
			Start()

		var processedCount int32

		// Start workers
		for i := 0; i < workerCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for file := range jobChan {
					content, err := os.ReadFile(file)
					if err != nil {
						errChan <- fmt.Errorf("error reading file %s: %v", file, err)
						continue
					}

					if len(content) > LimitChars {
						pterm.Warning.Printf("File %s with %d exceeds %d characters, truncating...\n", file, len(content), LimitChars)
						content = content[:LimitChars]
					}

					if dryrun {
						atomic.AddInt32(&processedCount, 1)
						progress.UpdateTitle(fmt.Sprintf("(Dryrun) Processing %d/%d", atomic.LoadInt32(&processedCount), len(files)))
						<-time.After(50 * time.Millisecond)
						progress.Increment()
						continue
					}

					progress.UpdateTitle(fmt.Sprintf("Summarizing %s", file))
					summary, _, err := summarizerInstance.Summarize(string(content), file, prompt)
					if err != nil {
						errChan <- fmt.Errorf("error summarizing file %s: %v", file, err)
						continue
					}

					err = summarizer.InjectSummary(file, summary, hash)
					if err != nil {
						errChan <- fmt.Errorf("error injecting summary into file %s: %v", file, err)
					}

					atomic.AddInt32(&processedCount, 1)
					progress.UpdateTitle(fmt.Sprintf("Processed %d/%d", atomic.LoadInt32(&processedCount), len(files)))
					progress.Increment()
				}
			}()
		}

		// Send jobs to workers
		for _, file := range files {
			jobChan <- file
		}
		close(jobChan)

		// Wait for all workers to complete
		wg.Wait()
		close(errChan)

		progress.Stop()

		// Handle errors
		errorCount := 0
		for err := range errChan {
			pterm.Error.Println(err)
			errorCount++
		}

		pterm.Success.Printf("Summarization completed in %v\n", time.Since(start))
		if errorCount > 0 {
			pterm.Warning.Printf("Encountered %d errors during processing\n", errorCount)
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
	rootCmd.PersistentFlags().BoolVar(&dryrun, "dryrun", false, "Dry run mode - stops before making API calls")

	rootCmd.MarkPersistentFlagRequired("path")
}

func initConfig() {
	viper.AutomaticEnv()
}
