package cmd

import (
	"context"
	"fmt"
	"os"
	"time"
	"wget/internal/config"
	"wget/internal/lib/scraper"

	"github.com/spf13/cobra"
)

var cfg config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wget <url>",
	Short: "wget is used to download urls and resources recursively",
	Long: `Use wget to download urls and resources recursively.
You only need to plug one link in.
All the other links will get fetched.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 || len(args) > 1 {
			return cmd.Usage()
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg.StartURL = args[0]
		s := scraper.New(&cfg)

		fmt.Println("downloading...")
		fmt.Println("max depth", cfg.MaxDepth)
		fmt.Println("timeout", cfg.Timeout)
		fmt.Println("workers", cfg.Workers)

		if err := s.Start(ctx); err != nil {
			println("failed scraping", err.Error())
			return err
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wget.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().IntVarP(&cfg.MaxDepth, "depth", "d", 2, "max recursion depth")
	rootCmd.Flags().StringVarP(&cfg.OutputDir, "output-dir", "o", "", "output directory")
	rootCmd.Flags().IntVarP(&cfg.Workers, "workers", "w", 5, "number of workers")
	rootCmd.Flags().DurationVarP(&cfg.Timeout, "timeout", "t", 10*time.Second, "timeout")
	rootCmd.Flags().StringVarP(&cfg.UserAgent, "user-agent", "a", "Wget/1.0", "user agent")
	rootCmd.Flags().BoolVarP(&cfg.IgnoreRobots, "ignore-robots", "r", false, "ignore robots.txt")
}
