package cmd

import (
	"os"
	"strings"
	"telnet/internal/config"
	"telnet/internal/run"
	"time"

	"github.com/spf13/cobra"
)

var opt config.Options

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "telnet <host:port>",
	Short: "A telnet client",
	Long:  `A longer description for a telnet client.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}
		opt.Address = args[0]

		if strings.Index(opt.Address, ":") < 0 {
			opt.Address += ":23"
		}

		return run.Run(&opt)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().DurationVar(&opt.Timeout, "timeout", time.Second*10, "timeout")
}
