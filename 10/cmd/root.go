package cmd

import (
	"fmt"
	"os"

	"gosort/internal/options"
	"gosort/internal/run"

	"github.com/spf13/cobra"
)

var (
	opt options.Options
)

var rootCmd = &cobra.Command{
	Use:   "gosort [file]",
	Short: "A simplified analogue of UNIX sort",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Disallow conflicting sort mode flags
		modeFlags := 0
		if opt.Numeric {
			modeFlags++
		}
		if opt.HumanNumeric {
			modeFlags++
		}
		if opt.Month {
			modeFlags++
		}
		if modeFlags > 1 {
			return fmt.Errorf("conflicting flags: choose only one of -n, -h, -M")
		}

		return run.Run(opt, args)
	},
}

func init() {
	// Define a dummy --help flag (no shorthand) to prevent Cobra from reserving -h.
	rootCmd.Flags().Bool("help", false, "")
	_ = rootCmd.Flags().MarkHidden("help")

	rootCmd.Flags().IntVarP(&opt.KeyColumn, "key", "k", 0, "sort by column number N (1-based, default: whole line; separator: tab)")
	rootCmd.Flags().BoolVarP(&opt.Numeric, "numeric-sort", "n", false, "compare according to numeric value")
	rootCmd.Flags().BoolVarP(&opt.Reverse, "reverse", "r", false, "reverse the result of comparisons")
	rootCmd.Flags().BoolVarP(&opt.Unique, "unique", "u", false, "output only the first of an equal run")
	rootCmd.Flags().BoolVarP(&opt.Month, "month-sort", "M", false, "compare by month name (Jan, Feb, ... Dec)")
	rootCmd.Flags().BoolVarP(&opt.IgnoreTrailingBlanks, "ignore-trailing-blanks", "b", false, "ignore trailing blanks in key comparisons")
	rootCmd.Flags().BoolVarP(&opt.Check, "check", "c", false, "check whether input is sorted; do not sort")
	rootCmd.Flags().BoolVarP(&opt.HumanNumeric, "human-numeric-sort", "h", false, "compare human readable numbers (e.g. 2K, 1M)")

	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = false
}

// Execute starts the program
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// For -c we may return a non-zero error after printing message already.
		os.Exit(1)
	}
}
