package cmd

import (
	"grep/internal/config"
	"grep/internal/run"
	"grep/internal/stream"
	"os"

	"github.com/spf13/cobra"
)

var opt config.Flags

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "grep [FILENAME (optional)]",
	Short: "Grep -- a utility to search for regular expressions.",
	Long:  `Search text for regular expressions / plain text strings. Don't input a file name if you want to read STDIN.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 || len(args) > 2 {
			return cmd.Usage()
		}
		return run.Run(args, opt, stream.NewProcessor(&opt))
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
	var c int

	rootCmd.Flags().IntVarP(&opt.After, "after-context", "A", 0, "show N lines after each found expression")
	rootCmd.Flags().IntVarP(&opt.Before, "before-context", "B", 0, "show N lines before each found expression")

	rootCmd.Flags().IntVarP(&c, "context", "C", 0, "show N lines before and after each found expression")
	opt.After = max(opt.After, c)
	opt.Before = max(opt.Before, c)

	rootCmd.Flags().BoolVarP(&opt.OnlyCount, "count", "c", false, "show only matching count")
	rootCmd.Flags().BoolVarP(&opt.IgnoreCase, "ignore-case", "i", false, "ignore case matching")
	rootCmd.Flags().BoolVarP(&opt.Invert, "invert", "v", false, "invert matching")
	rootCmd.Flags().BoolVarP(&opt.FixedString, "fixed-string", "F", false, "fix string instead of regexp")
	rootCmd.Flags().BoolVarP(&opt.PrintNumbers, "print-numbers", "n", false, "print line numbers")
}
