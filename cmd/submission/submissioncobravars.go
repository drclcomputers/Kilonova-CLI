// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package submission

import (
	"fmt"
	utility "kilocli/cmd/utility"
	"strconv"

	"github.com/spf13/cobra"
)

var UploadCodeCmd = &cobra.Command{
	Use:   "submit [ID] [LANGUAGE] [FILENAME] [Contest ID (optional)]",
	Short: "Submit solution to problem.",
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 3 {
			uploadCode(args[0], args[1], args[2], "NO")
		} else {
			uploadCode(args[0], args[1], args[2], args[3])
		}
	},
}

var CheckLangsCmd = &cobra.Command{
	Use:   "langs [ID]",
	Short: "View available languages for solutions.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		CheckLanguages(args[0], 1)
	},
}

var SubmissionCmd = &cobra.Command{
	Use:   "submission [command] ...",
	Short: "View details about submissions.",
}

var PrintSubmissionsCmd = &cobra.Command{
	Use:   "list [Problem ID or all (all problems)] [User ID, me (personal submissions), all (all users)] [1st page] [last page]",
	Short: "View sent submissions to a problem.",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		FirstPage, err := strconv.Atoi(args[2])
		if err != nil {
			utility.LogError(fmt.Errorf("invalid first page number: %v", err))
			return
		}

		LastPage, err := strconv.Atoi(args[3])
		if err != nil {
			utility.LogError(fmt.Errorf("invalid last page number: %v", err))
			return
		}

		printSubmissions(args[0], args[1], FirstPage, LastPage)
	},
}

var PrintSubmissionInfoCmd = &cobra.Command{
	Use:   "info [Submission ID]",
	Short: "View a detailed description of a sent submission.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		printDetailsSubmission(args[0])
	},
}
