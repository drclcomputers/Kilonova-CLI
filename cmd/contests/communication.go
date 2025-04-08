// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package contests

import (
	"bytes"
	"encoding/json"
	"fmt"
	utility "kilocli/cmd/utility"
	u "net/url"

	"github.com/charmbracelet/bubbles/table"
)

// communication
func createAnnouncement(contestID, text string) {
	url := fmt.Sprintf(utility.URL_CONTEST_CREATE_ANNOUNCEMENT, contestID)

	formData := u.Values{
		"text": {text},
	}

	body, err := utility.MakePostRequest(url, bytes.NewBufferString(formData.Encode()), utility.RequestFormAuth)
	if err != nil {
		utility.LogError(err)
		return
	}

	var data utility.KilonovaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		utility.LogError(fmt.Errorf("failed to decode response: %w", err))
		return
	}

	if data.Status != utility.SUCCESS {
		fmt.Println("Failed to create an announcement!")
	} else {
		fmt.Println(string(data.Data))
	}
}

func updateAnnouncement(contestID, announcementID, text string) {
	url := fmt.Sprintf(utility.URL_CONTEST_UPDATE_ANNOUNCEMENT, contestID)

	formData := u.Values{
		"text": {text},
		"id":   {announcementID},
	}

	body, err := utility.MakePostRequest(url, bytes.NewBufferString(formData.Encode()), utility.RequestFormAuth)
	if err != nil {
		utility.LogError(err)
		return
	}

	var data utility.KilonovaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		utility.LogError(fmt.Errorf("failed to decode response: %w", err))
		return
	}

	if data.Status != utility.SUCCESS {
		fmt.Println("Failed to update announcement!")
	} else {
		fmt.Println(string(data.Data))
	}
}

func deleteAnnouncement(contestID, announcementID string) {
	url := fmt.Sprintf(utility.URL_CONTEST_DELETE_ANNOUNCEMENT, contestID)

	formData := u.Values{
		"id": {announcementID},
	}

	body, err := utility.MakePostRequest(url, bytes.NewBufferString(formData.Encode()), utility.RequestFormAuth)
	if err != nil {
		utility.LogError(err)
		return
	}

	var data utility.KilonovaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		utility.LogError(fmt.Errorf("failed to decode response: %w", err))
		return
	}

	if data.Status != utility.SUCCESS {
		fmt.Println("Failed to delete announcement!")
	} else {
		fmt.Println(string(data.Data))
	}
}

func askQuestion(contestID, text string) {
	url := fmt.Sprintf(utility.URL_CONTEST_ASK_QUESTION, contestID)

	formData := u.Values{
		"text": {text},
	}

	body, err := utility.MakePostRequest(url, bytes.NewBufferString(formData.Encode()), utility.RequestFormAuth)
	if err != nil {
		utility.LogError(err)
		return
	}

	var data utility.KilonovaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		utility.LogError(fmt.Errorf("failed to decode response: %w", err))
		return
	}

	if data.Status != utility.SUCCESS {
		fmt.Println("Failed to ask question!")
	} else {
		fmt.Println(string(data.Data))
	}
}

func answerQuestion(contestID, question_id, text string) {
	url := fmt.Sprintf(utility.URL_CONTEST_RESPOND_QUESTION, contestID)

	formData := u.Values{
		"text":       {text},
		"questionID": {question_id},
	}

	body, err := utility.MakePostRequest(url, bytes.NewBufferString(formData.Encode()), utility.RequestFormAuth)
	if err != nil {
		utility.LogError(err)
		return
	}

	var data utility.KilonovaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		utility.LogError(fmt.Errorf("failed to decode response: %w", err))
		return
	}

	if data.Status != utility.SUCCESS {
		fmt.Println("Failed to respond to question!")
	} else {
		fmt.Println(string(data.Data))
	}
}

func viewAnnouncementsContest(contestID string) {
	url := fmt.Sprintf(utility.URL_CONTEST_ANNOUNCEMENTS, contestID)
	body, err := utility.MakeGetRequest(url, nil, utility.RequestNone)
	if err != nil {
		utility.LogError(err)
		return
	}
	var data ContestData
	if err = json.Unmarshal(body, &data); err != nil {
		utility.LogError(err)
		return
	}

	ok := false

	var Rows []table.Row

	if data.Status != utility.SUCCESS {
		utility.LogError(fmt.Errorf("couldn't retrieve announcements"))
		return
	} else {
		for _, announ := range data.Data {
			ok = true
			parsedTime, err := utility.ParseTime(announ.Time)
			if err != nil {
				utility.LogError(err)
				return
			}
			Rows = append(Rows, table.Row{
				fmt.Sprintf("%d", announ.ID),
				parsedTime,
				announ.Text,
			})
		}
	}

	if !ok {
		fmt.Println("No announcements have been made!")
	} else {
		Columns := []table.Column{
			{Title: "ID", Width: 4},
			{Title: "Time", Width: 19},
			{Title: "Text", Width: 57},
		}

		utility.RenderTable(Columns, Rows, 1)
	}
}

func viewAllQuestionsContest(contestID string) {
	url := fmt.Sprintf(utility.URL_CONTEST_ALL_QUESTIONS, contestID)
	body, err := utility.MakeGetRequest(url, nil, utility.RequestNone)
	if err != nil {
		utility.LogError(err)
		return
	}
	var data ContestQuestions
	if err = json.Unmarshal(body, &data); err != nil {
		utility.LogError(err)
		return
	}

	ok := false

	var Rows []table.Row

	if data.Status != utility.SUCCESS {
		utility.LogError(fmt.Errorf("couldn't retrieve questions"))
		return
	} else {
		for _, quest := range data.Data {
			ok = true
			formattedTime, err := utility.ParseTime(quest.Time)
			if err != nil {
				utility.LogError(err)
				return
			}

			Rows = append(Rows, table.Row{
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
		Columns := []table.Column{
			{Title: "Time", Width: 19},
			{Title: "ID", Width: 5},
			{Title: "User", Width: 4},
			{Title: "Text", Width: 35},
			{Title: "Response", Width: 17},
		}

		utility.RenderTable(Columns, Rows, 1)
	}
}

func viewMyQuestionsContest(contestID string) {
	url := fmt.Sprintf(utility.URL_CONTEST_YOUR_QUESTIONS, contestID)
	body, err := utility.MakeGetRequest(url, nil, utility.RequestNone)
	if err != nil {
		utility.LogError(err)
		return
	}
	var data ContestQuestions
	if err = json.Unmarshal(body, &data); err != nil {
		utility.LogError(err)
		return
	}

	ok := false

	var Rows []table.Row

	if data.Status != utility.SUCCESS {
		utility.LogError(fmt.Errorf("couldn't retrieve questions"))
		return
	}
	for _, quest := range data.Data {
		ok = true
		formattedTime, err := utility.ParseTime(quest.Time)
		if err != nil {
			utility.LogError(err)
			return
		}

		Rows = append(Rows, table.Row{
			formattedTime,
			fmt.Sprintf("%d", quest.Id),
			quest.Text,
			quest.Response,
		})

	}

	if !ok {
		fmt.Println("No questions have been asked!")
		return
	}
	Columns := []table.Column{
		{Title: "Time", Width: 19},
		{Title: "ID", Width: 5},
		{Title: "Text", Width: 35},
		{Title: "Response", Width: 21},
	}

	utility.RenderTable(Columns, Rows, 1)

}
