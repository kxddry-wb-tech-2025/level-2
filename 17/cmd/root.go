package cmd

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var opt Options

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
		address := strings.Split(args[0], ":")
		if len(address) != 2 {
			return cmd.Usage()
		}
		if n, err := strconv.Atoi(address[1]); err != nil {
			return err
		} else {
			opt.Port = n
		}
		opt.Host = address[0]

		return nil
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
	rootCmd.Flags().DurationVar(&opt.Timeout, "timeout", time.Second*10, "Timeout")
}

type Options struct {
	Timeout time.Duration
	Host    string
	Port    int
}
