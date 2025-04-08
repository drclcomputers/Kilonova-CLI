// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package contests

import (
	utility "kilocli/cmd/utility"
	"strconv"

	"github.com/spf13/cobra"
)

var createContestCmd = &cobra.Command{
	Use:   "create [name] [type]",
	Short: "Create a contest.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		createContest(args[0], args[1])
	},
}

var registerContestCmd = &cobra.Command{
	Use:   "register [ID]",
	Short: "Register in a contest.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		registerContest(args[0])
	},
}

var startContestCmd = &cobra.Command{
	Use:   "start [ID]",
	Short: "Start a contest.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		startContest(args[0])
	},
}

var deleteContestCmd = &cobra.Command{
	Use:   "delete [ID]",
	Short: "Delete a contest.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deleteContest(args[0])
	},
}

var viewAnnouncementsContestCmd = &cobra.Command{
	Use:   "announcements [ID]",
	Short: "View contest announcements.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		viewAnnouncementsContest(args[0])
	},
}

var viewAllQuestionsContestCmd = &cobra.Command{
	Use:   "allquestions [ID]",
	Short: "View all contest questions.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		viewAllQuestionsContest(args[0])
	},
}

var viewMyQuestionsContestCmd = &cobra.Command{
	Use:   "myquestions [ID]",
	Short: "View yout contest questions.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		viewMyQuestionsContest(args[0])
	},
}

var askQuestionContestCmd = &cobra.Command{
	Use:   "ask [ID] [text]",
	Short: "Ask a question in a contest.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		askQuestion(args[0], args[1])
	},
}

var respondQuestionContestCmd = &cobra.Command{
	Use:   "respond [ID] [question ID] [text]",
	Short: "Respond to a question.",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		answerQuestion(args[0], args[1], args[2])
	},
}

var createAnnouncementContestCmd = &cobra.Command{
	Use:   "createannoun [ID] [text]",
	Short: "Create an announcement.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		createAnnouncement(args[0], args[1])
	},
}

var updateAnnouncementContestCmd = &cobra.Command{
	Use:   "updateannoun [ID] [Announ. ID] [text]",
	Short: "Update an announcement.",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		updateAnnouncement(args[0], args[1], args[2])
	},
}

var deleteAnnouncementContestCmd = &cobra.Command{
	Use:   "delannoun [Contest ID] [Announ. ID]",
	Short: "Delete an announcement.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		deleteAnnouncement(args[0], args[1])
	},
}

var updateProblemsContestCmd = &cobra.Command{
	Use:   "update [problem_1] [problem_2] ... [problem_n]",
	Short: "Update the problems in your contest.",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		contestID := args[0]
		problemIDs := args[1:]

		updateProblems(contestID, problemIDs)
	},
}

var showProblemsContestCmd = &cobra.Command{
	Use:   "problems [ID]",
	Short: "Show problems in the contest.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		showProblems(args[0])
	},
}

var showInfoContestCmd = &cobra.Command{
	Use:   "info [ID]",
	Short: "Show a brief description of the contest.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		infoContest(args[0], "1")
	},
}

var modifyInfoContestCmd = &cobra.Command{
	Use:   "settings [command]",
	Short: "Adjust contest settings.",
}

var modifyStartTimeContestCmd = &cobra.Command{
	Use:   "start [ID] [time formatted like (2006-08-09 12:30:00)]",
	Short: "Modify contest starting time.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		NewTime, err := utility.ParseTime(args[1])
		if err != nil {
			utility.LogError(err)
			return
		}
		modifyGeneralContest(ContestUpdate{args[0], "start_time", NewTime})
	},
}

var modifyEndTimeContestCmd = &cobra.Command{
	Use:   "end [ID] [time formatted like (2006-08-09 12:30:00)]",
	Short: "Modify contest ending time.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		NewTime, err := utility.ParseTime(args[1])
		if err != nil {
			utility.LogError(err)
			return
		}

		modifyGeneralContest(ContestUpdate{args[0], "end_time", NewTime})
	},
}

var modifyMaxSubsContestCmd = &cobra.Command{
	Use:   "maxsubs [ID] [nr]",
	Short: "Modify contest max submissions per problem.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := strconv.Atoi(args[1])
		if err != nil {
			utility.LogError(err)
			return
		}
		modifyGeneralContest(ContestUpdate{args[0], "max_subs", args[1]})
	},
}

var modifyVisibleContestCmd = &cobra.Command{
	Use:   "visible [ID] [true or false]",
	Short: "Modify contest visibility.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		value, err := utility.ValidateBoolean(args[1])
		if err != nil {
			utility.LogError(err)
			return
		}
		modifyGeneralContest(ContestUpdate{args[0], "visible", strconv.FormatBool(value)})
	},
}

var modifyRegisterDuringContestCmd = &cobra.Command{
	Use:   "registduring [ID] [true or false]",
	Short: "Modify registering during contest.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		value, err := utility.ValidateBoolean(args[1])
		if err != nil {
			utility.LogError(err)
			return
		}
		modifyGeneralContest(ContestUpdate{args[0], "register_during_contest", strconv.FormatBool(value)})
	},
}

var modifyPublicLeaderboardContestCmd = &cobra.Command{
	Use:   "publicleader [ID] [true or false]",
	Short: "Modify leaderboard visibily to the public.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		value, err := utility.ValidateBoolean(args[1])
		if err != nil {
			utility.LogError(err)
			return
		}
		modifyGeneralContest(ContestUpdate{args[0], "public_leaderboard", strconv.FormatBool(value)})
	},
}

var leaderboardContestCmd = &cobra.Command{
	Use:   "leaderboard [ID]",
	Short: "Show contest leaderboard.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		leaderboard(args[0])
	},
}
