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

	"github.com/charmbracelet/bubbles/table"
)

// communication
func createAnnouncement(contestID, text string) {
	url := fmt.Sprintf(internal.URL_CONTEST_CREATE_ANNOUNCEMENT, contestID)

	formData := u.Values{
		"text": {text},
	}

	body, err := internal.MakePostRequest(url, bytes.NewBufferString(formData.Encode()), internal.RequestFormAuth)
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
		fmt.Println("Failed to create an announcement!")
	} else {
		fmt.Println(data.Data)
	}
}

func updateAnnouncement(contestID, announcementID, text string) {
	url := fmt.Sprintf(internal.URL_CONTEST_UPDATE_ANNOUNCEMENT, contestID)

	formData := u.Values{
		"text": {text},
		"id":   {announcementID},
	}

	body, err := internal.MakePostRequest(url, bytes.NewBufferString(formData.Encode()), internal.RequestFormAuth)
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
		fmt.Println("Failed to update announcement!")
	} else {
		fmt.Println(data.Data)
	}
}

func deleteAnnouncement(contestID, announcementID string) {
	url := fmt.Sprintf(internal.URL_CONTEST_DELETE_ANNOUNCEMENT, contestID)

	formData := u.Values{
		"id": {announcementID},
	}

	body, err := internal.MakePostRequest(url, bytes.NewBufferString(formData.Encode()), internal.RequestFormAuth)
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
		fmt.Println("Failed to delete announcement!")
	} else {
		fmt.Println(data.Data)
	}
}

func askQuestion(contestID, text string) {
	url := fmt.Sprintf(internal.URL_CONTEST_ASK_QUESTION, contestID)

	formData := u.Values{
		"text": {text},
	}

	body, err := internal.MakePostRequest(url, bytes.NewBufferString(formData.Encode()), internal.RequestFormAuth)
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
		fmt.Println("Failed to ask question!")
	} else {
		fmt.Println(data.Data)
	}
}

func answerQuestion(contestID, questionID, text string) {
	url := fmt.Sprintf(internal.URL_CONTEST_RESPOND_QUESTION, contestID)

	formData := u.Values{
		"text":       {text},
		"questionID": {questionID},
	}

	body, err := internal.MakePostRequest(url, bytes.NewBufferString(formData.Encode()), internal.RequestFormAuth)
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
		fmt.Println("Failed to respond to question!")
	} else {
		fmt.Println(data.Data)
	}
}

func viewAnnouncementsContest(contestID string) {
	url := fmt.Sprintf(internal.URL_CONTEST_ANNOUNCEMENTS, contestID)
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
		internal.LogError(fmt.Errorf("couldn't retrieve announcements"))
		return
	} else {
		for _, announ := range data.Data {
			ok = true
			parsedTime, err := internal.ParseTime(announ.Time)
			if err != nil {
				internal.LogError(err)
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

		internal.RenderTable(Columns, Rows, 1)
	}
}

func viewAllQuestionsContest(contestID string) {
	url := fmt.Sprintf(internal.URL_CONTEST_ALL_QUESTIONS, contestID)
	body, err := internal.MakeGetRequest(url, nil, internal.RequestNone)
	if err != nil {
		internal.LogError(err)
		return
	}
	var data ContestQuestions
	if err = json.Unmarshal(body, &data); err != nil {
		internal.LogError(err)
		return
	}

	ok := false

	var Rows []table.Row

	if data.Status != internal.SUCCESS {
		internal.LogError(fmt.Errorf("couldn't retrieve questions"))
		return
	} else {
		for _, quest := range data.Data {
			ok = true
			formattedTime, err := internal.ParseTime(quest.Time)
			if err != nil {
				internal.LogError(err)
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

		internal.RenderTable(Columns, Rows, 1)
	}
}

func viewMyQuestionsContest(contestID string) {
	url := fmt.Sprintf(internal.URL_CONTEST_YOUR_QUESTIONS, contestID)
	body, err := internal.MakeGetRequest(url, nil, internal.RequestNone)
	if err != nil {
		internal.LogError(err)
		return
	}
	var data ContestQuestions
	if err = json.Unmarshal(body, &data); err != nil {
		internal.LogError(err)
		return
	}

	ok := false

	var Rows []table.Row

	if data.Status != internal.SUCCESS {
		internal.LogError(fmt.Errorf("couldn't retrieve questions"))
		return
	}
	for _, quest := range data.Data {
		ok = true
		formattedTime, err := internal.ParseTime(quest.Time)
		if err != nil {
			internal.LogError(err)
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

	internal.RenderTable(Columns, Rows, 1)

}
