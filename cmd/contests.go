// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	u "net/url"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var ContestCmd = &cobra.Command{
	Use:   "contest [command] ...",
	Short: "Manage contests",
}

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
	Use:   "questions [ID]",
	Short: "View all contest questions.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		viewAllQuestionsContest(args[0])
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
	Use:   "update [pb_1] [pb_2] ... [pb_n]",
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

func init() {
	ContestCmd.AddCommand(createContestCmd)
	ContestCmd.AddCommand(deleteContestCmd)
	ContestCmd.AddCommand(registerContestCmd)
	ContestCmd.AddCommand(startContestCmd)
	ContestCmd.AddCommand(viewAnnouncementsContestCmd)
	ContestCmd.AddCommand(viewAllQuestionsContestCmd)
	ContestCmd.AddCommand(askQuestionContestCmd)
	ContestCmd.AddCommand(createAnnouncementContestCmd)
	ContestCmd.AddCommand(updateAnnouncementContestCmd)
	ContestCmd.AddCommand(deleteAnnouncementContestCmd)
	ContestCmd.AddCommand(updateProblemsContestCmd)
	ContestCmd.AddCommand(showProblemsContestCmd)

	rootCmd.AddCommand(ContestCmd)
}

type contestdata struct {
	Status string `json:"status"`
	Data   []struct {
		Text     string `json:"text"`
		Time     string `json:"created_at"`
		ID       int    `json:"id"`
		Name     string `json:"name"`
		MaxScore int    `json:"max_score"`
	}
}

type contestquestions struct {
	Status string `json:"status"`
	Data   []struct {
		Text          string `json:"text"`
		Time          string `json:"asked_at"`
		RespondedTIme string `json:"responded_at"`
		Response      string `json:"response"`
		Id            int    `json:"id"`
		AuthorID      int    `json:"author_id"`
	}
}

// contest
func createContest(name, contesttype string) {
	formData := u.Values{
		"name": {name},
		"type": {contesttype},
	}

	body, err := MakePostRequest(URL_CONTEST_CREATE, bytes.NewBufferString(formData.Encode()), RequestFormAuth)
	if err != nil {
		logError(err)
	}

	var data RawKilonovaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		logError(fmt.Errorf("failed to decode response: %w", err))
	}

	if data.Status != "success" {
		fmt.Println("Failed to create a contest!")
	} else {
		fmt.Println("Your contest's ID: #", string(data.Data))
	}
}

func registerContest(contest_id string) {
	url := fmt.Sprintf(URL_CONTEST_REGISTER, contest_id)
	body, err := PostJSON[KilonovaResponse](url, nil)
	if err != nil {
		logError(err)
	}
	fmt.Println(body.Data)
}

func startContest(contest_id string) {
	url := fmt.Sprintf(URL_CONTEST_START, contest_id)
	body, err := PostJSON[KilonovaResponse](url, nil)
	if err != nil {
		logError(err)
	}
	fmt.Println(body.Data)
}

func deleteContest(contest_id string) {
	url := fmt.Sprintf(URL_CONTEST_DELETE, contest_id)
	body, err := PostJSON[KilonovaResponse](url, nil)
	if err != nil {
		logError(err)
	}
	fmt.Println(body.Data)
}

// manage contest

func createAnnouncement(contest_id, text string) {
	url := fmt.Sprintf(URL_CONTEST_CREATE_ANNOUNCEMENT, contest_id)

	formData := u.Values{
		"text": {text},
	}

	body, err := MakePostRequest(url, bytes.NewBufferString(formData.Encode()), RequestFormAuth)
	if err != nil {
		logError(err)
	}

	var data KilonovaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		logError(fmt.Errorf("failed to decode response: %w", err))
	}

	if data.Status != "success" {
		fmt.Println("Failed to create an announcement!")
	} else {
		fmt.Println(string(data.Data))
	}
}

func updateAnnouncement(contest_id, announ_id, text string) {
	url := fmt.Sprintf(URL_CONTEST_UPDATE_ANNOUNCEMENT, contest_id)

	formData := u.Values{
		"text": {text},
		"id":   {announ_id},
	}

	body, err := MakePostRequest(url, bytes.NewBufferString(formData.Encode()), RequestFormAuth)
	if err != nil {
		logError(err)
	}

	var data KilonovaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		logError(fmt.Errorf("failed to decode response: %w", err))
	}

	if data.Status != "success" {
		fmt.Println("Failed to update announcement!")
	} else {
		fmt.Println(string(data.Data))
	}
}

func deleteAnnouncement(contest_id, id string) {
	url := fmt.Sprintf(URL_CONTEST_DELETE_ANNOUNCEMENT, contest_id)

	formData := u.Values{
		"id": {id},
	}

	body, err := MakePostRequest(url, bytes.NewBufferString(formData.Encode()), RequestFormAuth)
	if err != nil {
		logError(err)
	}

	var data KilonovaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		logError(fmt.Errorf("failed to decode response: %w", err))
	}

	if data.Status != "success" {
		fmt.Println("Failed to delete announcement!")
	} else {
		fmt.Println(string(data.Data))
	}
}

func askQuestion(contest_id, text string) {
	url := fmt.Sprintf(URL_CONTEST_ASK_QUESTION, contest_id)

	formData := u.Values{
		"text": {text},
	}

	body, err := MakePostRequest(url, bytes.NewBufferString(formData.Encode()), RequestFormAuth)
	if err != nil {
		logError(err)
	}

	var data KilonovaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		logError(fmt.Errorf("failed to decode response: %w", err))
	}

	if data.Status != "success" {
		fmt.Println("Failed to ask question!")
	} else {
		fmt.Println(string(data.Data))
	}
}

func updateProblems(contest_id string, problems_id []string) {
	url := fmt.Sprintf(URL_CONTEST_UPDATE_PROBLEMS, contest_id)

	var problems_id_int []int
	for _, s := range problems_id {
		num, err := strconv.Atoi(s)
		if err != nil {
			logError(err)
		}
		problems_id_int = append(problems_id_int, num)
	}

	payload := map[string]interface{}{
		"list": problems_id_int,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		logError(err)
	}

	body, err := MakePostRequest(url, bytes.NewBuffer(jsonBytes), RequestJSON)
	if err != nil {
		logError(err)
	}

	var data KilonovaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		logError(fmt.Errorf("failed to decode response: %w", err))
	}

	if data.Status != "success" {
		fmt.Println("Failed to update problems!")
	} else {
		fmt.Println(string(data.Data))
	}
}

func viewAnnouncementsContest(contest_id string) {
	url := fmt.Sprintf(URL_CONTEST_ANNOUNCEMENTS, contest_id)
	body, err := MakeGetRequest(url, nil, RequestNone)
	if err != nil {
		logError(err)
	}
	var data contestdata
	if err = json.Unmarshal(body, &data); err != nil {
		logError(err)
	}

	ok := false

	var rows []table.Row

	if data.Status != "success" {
		logError(fmt.Errorf("couldn't retrieve announcements"))
	} else {
		for _, announ := range data.Data {
			ok = true
			parsedTime, err := parseSubmissionTime(announ.Time)
			if err != nil {
				logError(err)
			}
			rows = append(rows, table.Row{
				fmt.Sprintf("%d", announ.ID),
				parsedTime,
				announ.Text,
			})
		}
	}

	if !ok {
		fmt.Println("No announcements have been made!")
	} else {
		columns := []table.Column{
			{Title: "ID", Width: 4},
			{Title: "Time", Width: 19},
			{Title: "Text", Width: 57},
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(20),
		)

		t.SetStyles(table.DefaultStyles())

		p := tea.NewProgram(&Model{table: t}, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			logError(fmt.Errorf("error running program: %w", err))
		}
	}
}

func viewAllQuestionsContest(contest_id string) {
	url := fmt.Sprintf(URL_CONTEST_ALL_QUESTIONS, contest_id)
	body, err := MakeGetRequest(url, nil, RequestNone)
	if err != nil {
		logError(err)
	}
	var data contestquestions
	if err = json.Unmarshal(body, &data); err != nil {
		logError(err)
	}

	ok := false

	var rows []table.Row

	if data.Status != "success" {
		logError(fmt.Errorf("couldn't retrieve questions"))
	} else {
		for _, quest := range data.Data {
			ok = true
			formattedTime, err := parseSubmissionTime(quest.Time)
			if err != nil {
				logError(err)
			}

			rows = append(rows, table.Row{
				formattedTime,
				fmt.Sprintf("%d", quest.Id),
				fmt.Sprintf("%d", quest.AuthorID),
				quest.Text,
				quest.Response,
			})

		}
	}

	if !ok {
		fmt.Println("No questions have been asked!")
	} else {
		columns := []table.Column{
			{Title: "Time", Width: 19},
			{Title: "ID", Width: 5},
			{Title: "User", Width: 4},
			{Title: "Text", Width: 35},
			{Title: "Response", Width: 17},
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(20),
		)

		t.SetStyles(table.DefaultStyles())

		p := tea.NewProgram(&Model{table: t}, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			logError(fmt.Errorf("error running program: %w", err))
		}
	}
}

func showProblems(contest_id string) {
	url := fmt.Sprintf(URL_CONTEST_PROBLEMS, contest_id)
	body, err := MakeGetRequest(url, nil, RequestNone)

	if err != nil {
		logError(err)
	}
	var data contestdata
	if err = json.Unmarshal(body, &data); err != nil {
		logError(err)
	}

	ok := false

	var rows []table.Row

	if data.Status != "success" {
		logError(fmt.Errorf("couldn't retrieve contest problems"))
	} else {
		for _, pb := range data.Data {
			ok = true
			rows = append(rows, table.Row{
				fmt.Sprintf("%d", pb.ID),
				pb.Name,
				fmt.Sprintf("%d", pb.MaxScore),
			})
		}
	}

	if !ok {
		fmt.Println("No problems have been added!")
	} else {
		columns := []table.Column{
			{Title: "ID", Width: 5},
			{Title: "Name", Width: 30},
			{Title: "Max Score", Width: 10},
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(20),
		)

		t.SetStyles(table.DefaultStyles())

		p := tea.NewProgram(&Model{table: t}, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			logError(fmt.Errorf("error running program: %w", err))
		}
	}

}
