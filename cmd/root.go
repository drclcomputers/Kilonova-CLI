package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kilocli",
	Short: "A CLI client for the competitive programming platform Kilonova ",
	Long: `Kilonova-CLI is a command-line interface (CLI) client designed for interacting 
with the Kilonova competitive programming platform. It enables users to view statements, 
search for problems, submit solutions, and retrieve submission results directly from 
the terminal.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kilocli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
