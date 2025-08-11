package cmd

import (
	"os"
	"shell/internal/shell"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "shell",
	Short: "A Basic unix shell written in Go",
	Long: `A shell written in Go.
It suppots:
- cd,
- pwd,
- echo,
- kill,
- ps,

- exec,
- pipelines,
- EOF / Ctrl+C,

- conditionals (&&, ||)
- environment variables with $VAR,
- redirects < or >`,
	Run: func(cmd *cobra.Command, args []string) {
		sh := shell.New()
		sh.Run()
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.shell.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
