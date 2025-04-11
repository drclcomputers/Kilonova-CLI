// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package user

import (
	"encoding/json"
	"fmt"
	"kncli/internal"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/bubbles/table"
)

// get info about user
func getUserBio(UserName string) string {
	res, err := http.Get(fmt.Sprintf("https://kilonova.ro/profile/%s", UserName))
	if err != nil {
		internal.LogError(err)
		return internal.ERROR
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		internal.LogError(err)
		return internal.ERROR
	}

	bio := doc.Find("div.segment-panel.reset-list.statement-content.enhance-tables p").First().Text()

	return bio
}

func userGetDetails(UserID, useCase string) bool {
	var url string
	if UserID == "me" {
		url = internal.URL_SELF
	} else {
		url = fmt.Sprintf(internal.URL_USER, UserID)
	}

	// Fetch user data
	ResponseBody, err := internal.MakeGetRequest(url, nil, internal.RequestFormGuest)
	if err != nil {
		internal.LogError(fmt.Errorf("error fetching user details: %w", err))
		return false
	}

	var dataUser UserDetailResponse
	if err := json.Unmarshal(ResponseBody, &dataUser); err != nil {
		internal.LogError(fmt.Errorf("error unmarshalling user data: %w", err))
		return false
	}

	if dataUser.Data.DisplayName == "" {
		dataUser.Data.DisplayName = "-"
	}

	switch useCase {
	case "isadmin":
		return dataUser.Data.Admin
	default:
		printUserDetails(dataUser)

		return false
	}
}

func userGetSolvedProblems(UserID string) {
	var url string
	if UserID == "me" {
		url = internal.URL_SELF_PROBLEMS
	} else {
		url = fmt.Sprintf(internal.URL_USER_PROBLEMS, UserID)
	}

	ResponseBody, err := internal.MakeGetRequest(url, nil, internal.RequestFormGuest)
	if err != nil {
		internal.LogError(fmt.Errorf("error fetching solved problems: %w", err))
		return
	}

	var dataUser UserSolvedProblems
	if err := json.Unmarshal(ResponseBody, &dataUser); err != nil {
		internal.LogError(fmt.Errorf("error unmarshalling solved problems: %w", err))
		return
	}

	Rows := prepareTableRows(dataUser)

	Columns := []table.Column{
		{Title: "Problem ID", Width: 7},
		{Title: "Name", Width: 20},
		{Title: "Source", Width: 40},
		{Title: "Score", Width: 7},
	}

	internal.RenderTable(Columns, Rows, 1)
}
