// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package contests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kncli/internal"
	u "net/url"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/spf13/cobra"
)

var shouldDownload = false

var ContestCmd = &cobra.Command{
	Use:   "contest [command] ...",
	Short: "Manage contests",
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
	leaderboardContestCmd.Flags().BoolVarP(&shouldDownload, "download_leader", "d", false, "shouldDownload leaderboard as a CSV file.")

	modifyInfoContestCmd.AddCommand(modifyStartTimeContestCmd)
	modifyInfoContestCmd.AddCommand(modifyEndTimeContestCmd)
	modifyInfoContestCmd.AddCommand(modifyMaxSubsContestCmd)
	modifyInfoContestCmd.AddCommand(modifyVisibleContestCmd)
	modifyInfoContestCmd.AddCommand(modifyRegisterDuringContestCmd)
	modifyInfoContestCmd.AddCommand(modifyPublicLeaderboardContestCmd)
}

type ContestData struct {
	Status string `json:"status"`
	Data   []struct {
		Text     string `json:"text"`
		Time     string `json:"created_at"`
		ID       int    `json:"id"`
		Name     string `json:"name"`
		MaxScore int    `json:"max_score"`
	}
}

type ContestInfo struct {
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

type ContestQuestions struct {
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

type ContestUpdate struct {
	ContestID string
	DataForm  string
	Value     string
}

// contest
func createContest(name, contestType string) {
	formData := u.Values{
		"name": {name},
		"type": {contestType},
	}

	body, err := internal.MakePostRequest(internal.URL_CONTEST_CREATE, bytes.NewBufferString(formData.Encode()), internal.RequestFormAuth)
	if err != nil {
		internal.LogError(err)
		return
	}

	var data internal.RawKilonovaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		internal.LogError(fmt.Errorf("failed to decode response: %w", err))
		return
	}

	if data.Status != internal.SUCCESS {
		fmt.Println("Failed to create a contest!")
		return
	}
	fmt.Println("Your contest's ID: #", string(data.Data))
}

func registerContest(contestID string) {
	url := fmt.Sprintf(internal.URL_CONTEST_REGISTER, contestID)
	body, err := internal.PostJSON[internal.KilonovaResponse](url, nil)
	if err != nil {
		internal.LogError(err)
		return
	}
	fmt.Println(body.Data)
}

func startContest(contestID string) {
	url := fmt.Sprintf(internal.URL_CONTEST_START, contestID)
	body, err := internal.PostJSON[internal.KilonovaResponse](url, nil)
	if err != nil {
		internal.LogError(err)
		return
	}
	fmt.Println(body.Data)
}

func deleteContest(contestID string) {
	url := fmt.Sprintf(internal.URL_CONTEST_DELETE, contestID)
	body, err := internal.PostJSON[internal.KilonovaResponse](url, nil)
	if err != nil {
		internal.LogError(err)
		return
	}
	fmt.Println(body.Data)
}

// manage contest

func updateProblems(contestID string, problemsID []string) {
	url := fmt.Sprintf(internal.URL_CONTEST_UPDATE_PROBLEMS, contestID)

	var problemsIDInt []int
	for _, s := range problemsID {
		num, err := strconv.Atoi(s)
		if err != nil {
			internal.LogError(err)
			return
		}
		problemsIDInt = append(problemsIDInt, num)
	}

	payload := map[string]interface{}{
		"list": problemsIDInt,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		internal.LogError(err)
		return
	}

	body, err := internal.MakePostRequest(url, bytes.NewBuffer(jsonBytes), internal.RequestJSON)
	if err != nil {
		internal.LogError(err)
		return
	}

	var data internal.KilonovaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		internal.LogError(fmt.Errorf("failed to decode response: %w", err))
		return
	}

	if data.Status != internal.SUCCESS {
		fmt.Println("Failed to update problems!")
	} else {
		fmt.Println(data.Data)
	}
}

func showProblems(contestID string) {
	url := fmt.Sprintf(internal.URL_CONTEST_PROBLEMS, contestID)
	body, err := internal.MakeGetRequest(url, nil, internal.RequestNone)

	if err != nil {
		internal.LogError(err)
		return
	}
	var data ContestData
	if err = json.Unmarshal(body, &data); err != nil {
		internal.LogError(err)
		return
	}

	ok := false

	var Rows []table.Row

	if data.Status != internal.SUCCESS {
		internal.LogError(fmt.Errorf("couldn't retrieve contest problems"))
		return
	} else {
		for _, problem := range data.Data {
			ok = true
			Rows = append(Rows, table.Row{
				fmt.Sprintf("%d", problem.ID),
				problem.Name,
				fmt.Sprintf("%d", problem.MaxScore),
			})
		}
	}

	if !ok {
		fmt.Println("No problems have been added!")
	} else {
		Columns := []table.Column{
			{Title: "ID", Width: 5},
			{Title: "Name", Width: 30},
			{Title: "Max Score", Width: 10},
		}

		internal.RenderTable(Columns, Rows, 1)
	}

}

func infoContest(contestID, useCase string) ContestInfo {
	url := fmt.Sprintf(internal.URL_CONTEST, contestID)
	body, err := internal.MakeGetRequest(url, nil, internal.RequestInfo)

	if err != nil {
		internal.LogError(err)
		return ContestInfo{}
	}
	var data ContestInfo
	if err = json.Unmarshal(body, &data); err != nil {
		internal.LogError(err)
		return ContestInfo{}
	}
	if useCase != "2" {
		parsedtime1, err := internal.ParseTime(data.Data.StartTime)
		if err != nil {
			internal.LogError(err)
			return ContestInfo{}
		}
		parsedtime2, err := internal.ParseTime(data.Data.EndTime)
		if err != nil {
			internal.LogError(err)
			return ContestInfo{}
		}
		fmt.Printf("Name: %s\nStart time: %s\nEnd time: %s\nMax submissions per problem: %d\n",
			data.Data.Name, parsedtime1, parsedtime2, data.Data.MaxSubs)
		fmt.Printf("Public leaderboard: %t\nVisibility: %t\nRegistering during contest: %t\n",
			data.Data.PublicLeaderboard, data.Data.Visible, data.Data.RegisterDuringContest)
	}
	return data
}

func modifyGeneralContest(update ContestUpdate) {
	url := fmt.Sprintf(internal.URL_CONTEST_UPDATE, update.ContestID)

	formData := u.Values{
		update.DataForm: {update.Value},
	}

	body, err := internal.MakePostRequest(url, bytes.NewBufferString(formData.Encode()), internal.RequestFormAuth)
	if err != nil {
		internal.LogError(err)
		return
	}

	var data internal.KilonovaResponse
	if err = json.Unmarshal(body, &data); err != nil {
		internal.LogError(err)
		return
	}

	fmt.Println(data.Data)
}
