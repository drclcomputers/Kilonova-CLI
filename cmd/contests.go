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
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

var Download bool = false

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
		layout := "2006-01-02 15:04:05"
		parsedTime, err := time.Parse(layout, args[1])
		if err != nil {
			logError(err)
		}
		NewTime := parsedTime.UTC().Format(time.RFC3339Nano)

		modifyGeneralContest(args[0], "start_time", NewTime)
	},
}

var modifyEndTimeContestCmd = &cobra.Command{
	Use:   "end [ID] [time formatted like (2006-08-09 12:30:00)]",
	Short: "Modify contest ending time.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		layout := "2006-01-02 15:04:05"
		parsedTime, err := time.Parse(layout, args[1])
		if err != nil {
			logError(err)
		}
		NewTime := parsedTime.UTC().Format(time.RFC3339Nano)

		modifyGeneralContest(args[0], "end_time", NewTime)
	},
}

var modifyMaxSubsContestCmd = &cobra.Command{
	Use:   "maxsubs [ID] [nr]",
	Short: "Modify contest max submissions per problem.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := strconv.Atoi(args[1])
		if err != nil {
			logError(err)
		}
		modifyGeneralContest(args[0], "max_subs", args[1])
	},
}

var modifyVisibleContestCmd = &cobra.Command{
	Use:   "visible [ID] [true or false]",
	Short: "Modify contest visibility.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if args[1] != "true" && args[1] != "false" {
			logError(fmt.Errorf("error: visibility must be either true or false"))
		}
		modifyGeneralContest(args[0], "visible", args[1])
	},
}

var modifyRegisterDuringContestCmd = &cobra.Command{
	Use:   "registduring [ID] [true or false]",
	Short: "Modify registering during contest.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if args[1] != "true" && args[1] != "false" {
			logError(fmt.Errorf("error: registering during contest must be either true or false"))
		}
		modifyGeneralContest(args[0], "register_during_contest", args[1])
	},
}

var modifyPublicLeaderboardContestCmd = &cobra.Command{
	Use:   "publicleader [ID] [true or false]",
	Short: "Modify leaderboard visibily to the public.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if args[1] != "true" && args[1] != "false" {
			logError(fmt.Errorf("error: public leaderboard must be either true or false"))
		}
		modifyGeneralContest(args[0], "public_leaderboard", args[1])
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

func init() {
	ContestCmd.AddCommand(createContestCmd)
	ContestCmd.AddCommand(deleteContestCmd)
	ContestCmd.AddCommand(registerContestCmd)
	ContestCmd.AddCommand(startContestCmd)
	ContestCmd.AddCommand(viewAnnouncementsContestCmd)
	ContestCmd.AddCommand(viewAllQuestionsContestCmd)
	ContestCmd.AddCommand(viewMyQuestionsContestCmd)
	ContestCmd.AddCommand(askQuestionContestCmd)
	ContestCmd.AddCommand(respondQuestionContestCmd)
	ContestCmd.AddCommand(createAnnouncementContestCmd)
	ContestCmd.AddCommand(updateAnnouncementContestCmd)
	ContestCmd.AddCommand(deleteAnnouncementContestCmd)
	ContestCmd.AddCommand(updateProblemsContestCmd)
	ContestCmd.AddCommand(showProblemsContestCmd)
	ContestCmd.AddCommand(showInfoContestCmd)
	ContestCmd.AddCommand(modifyInfoContestCmd)
	ContestCmd.AddCommand(leaderboardContestCmd)
	leaderboardContestCmd.Flags().BoolVarP(&Download, "download_leader", "d", false, "Download leaderboard as a CSV file.")

	modifyInfoContestCmd.AddCommand(modifyStartTimeContestCmd)
	modifyInfoContestCmd.AddCommand(modifyEndTimeContestCmd)
	modifyInfoContestCmd.AddCommand(modifyMaxSubsContestCmd)
	modifyInfoContestCmd.AddCommand(modifyVisibleContestCmd)
	modifyInfoContestCmd.AddCommand(modifyRegisterDuringContestCmd)
	modifyInfoContestCmd.AddCommand(modifyPublicLeaderboardContestCmd)

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

type contestinfo struct {
	Status string `json:"status"`
	Data   struct {
		StartTime             string      `json:"start_time"`
		EndTime               string      `json:"end_time"`
		MaxSubs               int         `json:"max_subs"`
		Name                  string      `json:"name"`
		Visible               bool        `json:"visible"`
		PublicLeaderboard     bool        `json:"public_leaderboard"`
		ChangeLeadboardFreeze bool        `json:"change_leaderboard_freeze"`
		IcpcSubmPenalty       json.Number `json:"icpc_submission_penalty"`
		LeadAdvFilter         bool        `json:"leaderboard_advanced_filter"`
		LeadStyle             string      `json:"leaderboard_style"`
		PerUserTime           json.Number `json:"per_user_time"`
		PublicJoin            bool        `json:"public_join"`
		QuestionCoolDown      json.Number `json:"question_contest"`
		RegisterDuringContest bool        `json:"register_during_contest"`
		SubmissionCooldown    json.Number `json:"submission_cooldown"`
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

type leaderboarddata struct {
	Status string `json:"status"`
	Data   struct {
		ProblemNames map[string]string `json:"problem_names"`
		Entries      []struct {
			User struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"user"`
			Scores map[string]int `json:"scores"`
			Total  int            `json:"total"`
		} `json:"entries"`
	} `json:"data"`
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

// communication
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

func updateAnnouncement(contest_id, announcement_id, text string) {
	url := fmt.Sprintf(URL_CONTEST_UPDATE_ANNOUNCEMENT, contest_id)

	formData := u.Values{
		"text": {text},
		"id":   {announcement_id},
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

func answerQuestion(contest_id, question_id, text string) {
	url := fmt.Sprintf(URL_CONTEST_RESPOND_QUESTION, contest_id)

	formData := u.Values{
		"text":       {text},
		"questionID": {question_id},
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
		fmt.Println("Failed to respond to question!")
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

// manage

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

func viewMyQuestionsContest(contest_id string) {
	url := fmt.Sprintf(URL_CONTEST_YOUR_QUESTIONS, contest_id)
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
			{Title: "Text", Width: 35},
			{Title: "Response", Width: 21},
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

func infoContest(contest_id, use_case string) contestinfo {
	url := fmt.Sprintf(URL_CONTEST, contest_id)
	body, err := MakeGetRequest(url, nil, RequestInfo)

	if err != nil {
		logError(err)
	}
	var data contestinfo
	if err = json.Unmarshal(body, &data); err != nil {
		logError(err)
	}
	if use_case != "2" {
		parsedtime1, err := parseSubmissionTime(data.Data.StartTime)
		if err != nil {
			logError(err)
		}
		parsedtime2, err := parseSubmissionTime(data.Data.EndTime)
		if err != nil {
			logError(err)
		}
		fmt.Printf("Name: %s\nStart time: %s\nEnd time: %s\nMax submissions per pb: %d\n",
			data.Data.Name, parsedtime1, parsedtime2, data.Data.MaxSubs)
		fmt.Printf("Public leaderboard: %t\nVisibility: %t\nRegistering during contest: %t\n",
			data.Data.PublicLeaderboard, data.Data.Visible, data.Data.RegisterDuringContest)
	}
	return data
}

func modifyGeneralContest(contest_id, datform, publicleader string) {
	url := fmt.Sprintf(URL_CONTEST_UPDATE, contest_id)

	formData := u.Values{
		datform: {publicleader},
	}

	body, err := MakePostRequest(url, bytes.NewBufferString(formData.Encode()), RequestFormAuth)
	if err != nil {
		logError(err)
	}

	var data KilonovaResponse
	if err = json.Unmarshal(body, &data); err != nil {
		logError(err)
	}

	fmt.Println(data.Data)
}

func downloadLeaderboard(contest_id string) {
	resp, err := MakeGetRequest(fmt.Sprintf(URL_CONTEST_ASSETS, contest_id), nil, RequestDownloadZip)
	if err != nil {
		logError(err)
	}

	homedir, err := os.Getwd()
	if err != nil {
		logError(fmt.Errorf("failed to get current working directory: %w", err))
		return
	}

	downFile := filepath.Join(homedir, "leaderboard_"+contest_id+".csv")
	outFile, err := os.Create(downFile)
	if err != nil {
		logError(fmt.Errorf("failed to create file %q: %w", downFile, err))
		return
	}
	defer outFile.Close()

	if err := os.WriteFile(downFile, resp, 0644); err != nil {
		logError(fmt.Errorf("failed to write to file %q: %w", downFile, err))
		return
	}

	fmt.Printf("Leaderboard to contest #%s saved to %q\n", contest_id, downFile)
}

func leaderboard(contest_id string) {
	url := fmt.Sprintf(URL_CONTEST_LEADERBOARD, contest_id)
	body, err := MakeGetRequest(url, nil, RequestNone)
	if err != nil {
		logError(err)
	}
	var data leaderboarddata
	if err = json.Unmarshal(body, &data); err != nil {
		logError(err)
	}

	var rows []table.Row

	for _, entry := range data.Data.Entries {
		var scores string
		for _, score := range entry.Scores {
			scores += fmt.Sprintf("%d  ", score)
		}
		rows = append(rows, table.Row{
			fmt.Sprintf("%d   %s ", entry.User.ID, entry.User.Name),
			scores,
			fmt.Sprintf("%d", entry.Total),
		})
	}

	var problemNamesTitle string
	for id, name := range data.Data.ProblemNames {
		problemNamesTitle += "| #" + id + " " + name + " "
	}

	problemNamesTitle += "|"

	columns := []table.Column{
		{Title: "ID | Name", Width: 25},
		{Title: problemNamesTitle, Width: 50},
		{Title: "Total", Width: 5},
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

	if Download {
		action := func() { downloadLeaderboard(contest_id) }
		if err := spinner.New().Title("Waiting for download...").Action(action).Run(); err != nil {
			logError(fmt.Errorf("error during source code download for submission #%s: %w", contest_id, err))
		}
	}
}
