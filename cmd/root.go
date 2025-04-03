// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "kilocli",
	Version: "v0.1.8",
	Short:   "A CLI client for the competitive programming platform Kilonova ",
	Long: `Kilonova-CLI is a command-line interface (CLI) client designed for interacting 
with the Kilonova competitive programming platform. It enables users to view statements, 
search for problems, submit solutions, and retrieve submission results directly from 
the terminal.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
