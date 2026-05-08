// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

// Package cmd implements the command-line interface for the Kilonova CLI application.
// It defines the root command and all subcommands for interacting with the Kilonova
// competitive programming platform.
package cmd

import (
	ai "kncli/cmd/ai"
	contest "kncli/cmd/contests"
	db "kncli/cmd/database"
	problem "kncli/cmd/problems"
	"kncli/cmd/project"
	"kncli/cmd/submission"
	"kncli/cmd/user"
	"kncli/internal"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands.
// It defines the CLI's name, version, and description, and serves as the
// container for all subcommands related to contests, problems, projects,
// submissions, users, and database operations.
var RootCmd = &cobra.Command{
	Use:     "kncli",
	Version: internal.Version,
	Short:   "A CLI client for the competitive programming platform Kilonova",
	Long: `Kilonova-CLI is a command-line interface (CLI) client designed for interacting 
with the Kilonova competitive programming platform. It enables users to view statements, 
search for problems, submit solutions, and retrieve submission results directly from 
the terminal.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This function is called by main.main() and only needs to happen once to the rootCmd.
// It handles command execution and error reporting, exiting with status 1 on failure.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		RootCmd.Println(err)
		os.Exit(1)
	}
}

// init initializes the root command by adding all subcommand groups.
// This function is automatically called by the Go runtime when the package is
// initialized and sets up the command structure for the CLI application.
func init() {
	// Add contest-related commands
	RootCmd.AddCommand(contest.ContestCmd)

	// Add problem-related commands
	RootCmd.AddCommand(problem.GetAssetsCmd)
	RootCmd.AddCommand(problem.SearchCmd)
	RootCmd.AddCommand(problem.PrintStatementCmd)

	// Add project-related commands
	RootCmd.AddCommand(project.InitProjectCmd)
	RootCmd.AddCommand(project.GetRandPbCmd)

	// Add submission-related commands
	RootCmd.AddCommand(submission.CheckLangsCmd)
	RootCmd.AddCommand(submission.UploadCodeCmd)
	RootCmd.AddCommand(submission.SubmissionCmd)

	// Add user-related commands
	RootCmd.AddCommand(user.SettingsCmd)
	RootCmd.AddCommand(user.SigninCmd)
	RootCmd.AddCommand(user.LogoutCmd)
	RootCmd.AddCommand(user.UserGetDetailsCmd)
	RootCmd.AddCommand(user.UserSolvedProblemsCmd)

	// Add database-related commands
	RootCmd.AddCommand(db.DatabaseCmd)

	// Add AI/Gemini-powered commands
	RootCmd.AddCommand(ai.SolveCmd)
	RootCmd.AddCommand(ai.ExplainCmd)
	RootCmd.AddCommand(ai.HintCmd)
	RootCmd.AddCommand(ai.DebugCmd)
	RootCmd.AddCommand(ai.AISearchCmd)
}
