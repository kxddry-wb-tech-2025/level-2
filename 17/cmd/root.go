package cmd

import (
	"net"
	"os"
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
		if !checkAddress(args[0]) {
			return cmd.Usage()
		}
		opt.Address = args[0]

		return run.Run(&opt)
	},
}

func checkAddress(addr string) bool {
	_, _, err := net.SplitHostPort(addr)
	return err == nil
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
