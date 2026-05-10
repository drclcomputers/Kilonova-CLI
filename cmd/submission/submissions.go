// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package submission

import (
	"encoding/json"
	"fmt"
	"kncli/internal"

	"github.com/charmbracelet/bubbles/table"
)

var shouldDownload = false

func init() {
	SubmissionCmd.AddCommand(PrintSubmissionsCmd)
	SubmissionCmd.AddCommand(PrintSubmissionInfoCmd)

	PrintSubmissionInfoCmd.Flags().BoolVarP(&shouldDownload, "download_source", "d", false, "Download the source code of a submission.")
}

// print submissions

func getSubmissionURL(UserID, ProblemID string, OffSet int) string {
	switch {
	case UserID == "all" && ProblemID == "all":
		return fmt.Sprintf(internal.URL_SUBMISSION_LIST_NO_FILTER, OffSet)
	case UserID == "all":
		return fmt.Sprintf(internal.URL_SUBMISSION_LIST_NO_USER, OffSet, ProblemID)
	case ProblemID == "all":
		return fmt.Sprintf(internal.URL_SUBMISSION_LIST_NO_PROBLEM, OffSet, UserID)
	default:
		return fmt.Sprintf(internal.URL_SUBMISSION_LIST, OffSet, ProblemID, UserID)
	}
}

func printSubmissions(ProblemID, UserID string, FirstPage, LastPage int) {
	if UserID == "me" {
		UserID = internal.GetUserID()
	}

	if FirstPage <= 0 || LastPage <= 0 {
		internal.LogError(fmt.Errorf("invalid pages: both FirstPage and LastPage must be positive integers"))
		return
	}

	if FirstPage > LastPage {
		internal.LogError(fmt.Errorf("FirstPage cannot be greater than LastPage"))
		return
	}

	var DataSubmissions SubmissionList
	var Rows []table.Row
	startOffset := (FirstPage - 1) * 50
	endOffset := LastPage * 50

	for OffSet := startOffset; OffSet < endOffset; OffSet += 50 {
		url := getSubmissionURL(UserID, ProblemID, OffSet)

		ResponseBody, err := internal.MakeGetRequest(url, nil, internal.RequestFormAuth)
		if err != nil {
			internal.LogError(err)
			continue
		}

		if err := json.Unmarshal(ResponseBody, &DataSubmissions); err != nil {
			internal.LogError(err)
			continue
		}

		submissions := DataSubmissions.Data.Submissions
		totalCount := DataSubmissions.Data.Count

		for _, sub := range submissions {
			formattedTime, err := internal.ParseTime(sub.CreatedAt)
			if err != nil {
				internal.LogError(err)
				continue
			}

			Rows = append(Rows, table.Row{
				fmt.Sprintf("%d", sub.ProblemID),
				fmt.Sprintf("%d", sub.UserID),
				fmt.Sprintf("%d", sub.Id),
				formattedTime,
				sub.Language,
				fmt.Sprintf("%.0f", sub.Score),
			})
		}

		// Stop early if we've fetched all available submissions
		if OffSet+50 >= totalCount {
			break
		}
	}

	Columns := []table.Column{
		{Title: "Pb ID", Width: 5},
		{Title: "User ID", Width: 7},
		{Title: "Submission ID", Width: 12},
		{Title: "Time", Width: 25},
		{Title: "Language", Width: 10},
		{Title: "Score", Width: 10},
	}

	if len(Rows) == 0 {
		internal.LogError(fmt.Errorf("no submissions found"))
		return
	}

	internal.RenderTable(Columns, Rows, 1)
}

func CheckLanguages(ProblemID string, useCase int) []string {
	url := fmt.Sprintf(internal.URL_LANGS_PB, ProblemID)
	ResponseBody, err := internal.MakeGetRequest(url, nil, internal.RequestNone)
	if err != nil {
		internal.LogError(fmt.Errorf("failed to make request for problem ID %s: %w", ProblemID, err))
		return nil
	}

	var langs Languages
	if err := json.Unmarshal(ResponseBody, &langs); err != nil {
		internal.LogError(fmt.Errorf("error unmarshalling languages response: %w", err))
		return nil
	}

	switch useCase {
	case 1:
		printLanguages(langs)
		return nil
	default:
		return extractLanguageNames(langs)
	}
}

func printLanguages(langs Languages) {
	for i, lang := range langs.Data {
		fmt.Printf("%d: %s\n", i+1, lang.Name)
	}
}

func extractLanguageNames(langs Languages) []string {
	var listLangs []string
	for _, lang := range langs.Data {
		listLangs = append(listLangs, lang.Name)
	}
	return listLangs
}
