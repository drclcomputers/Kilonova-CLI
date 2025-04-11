// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package cmd

import (
	contest "kncli/cmd/contests"
	db "kncli/cmd/database"
	problem "kncli/cmd/problems"
	"kncli/cmd/project"
	"kncli/cmd/submission"
	"kncli/cmd/user"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "kncli",
	Version: "v0.2.7",
	Short:   "A CLI client for the competitive programming platform Kilonova ",
	Long: `Kilonova-CLI is a command-line interface (CLI) client designed for interacting 
with the Kilonova competitive programming platform. It enables users to view statements, 
search for problems, submit solutions, and retrieve submission results directly from 
the terminal.`,
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		RootCmd.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(contest.ContestCmd)

	RootCmd.AddCommand(problem.GetAssetsCmd)
	RootCmd.AddCommand(problem.SearchCmd)
	RootCmd.AddCommand(problem.PrintStatementCmd)

	RootCmd.AddCommand(project.InitProjectCmd)
	RootCmd.AddCommand(project.GetRandPbCmd)

	RootCmd.AddCommand(submission.CheckLangsCmd)
	RootCmd.AddCommand(submission.UploadCodeCmd)
	RootCmd.AddCommand(submission.SubmissionCmd)

	RootCmd.AddCommand(user.SettingsCmd)
	RootCmd.AddCommand(user.SigninCmd)
	RootCmd.AddCommand(user.LogoutCmd)
	RootCmd.AddCommand(user.UserGetDetailsCmd)
	RootCmd.AddCommand(user.UserSolvedProblemsCmd)

	RootCmd.AddCommand(db.DatabaseCmd)

}
