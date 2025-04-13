// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package submission

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"kncli/internal"
	"strconv"
	"strings"
)

func init() {
	SubmissionCmd.AddCommand(ShowTestsSubmissionCmd)
}

func ShowTestsSubmission(ID string) {
	url := fmt.Sprintf(internal.URL_LATEST_SUBMISSION, ID)
	body, err := internal.MakeGetRequest(url, nil, internal.RequestNone)
	if err != nil {
		internal.LogError(err)
	}

	var data SubmissionDetails
	if err = json.Unmarshal(body, &data); err != nil {
		internal.LogError(err)
	}

	var Rows []table.Row

	for _, test := range data.Data.Subtests {
		Rows = append(Rows, table.Row{
			strconv.Itoa(data.Data.ProblemID),
			strconv.Itoa(test.ID), strings.TrimPrefix(test.Verdict, "translate:"),
			strconv.FormatFloat(test.Time, 'f', -1, 64),
			strconv.Itoa(test.Memory), strconv.Itoa(test.Percentage),
			strconv.Itoa(test.Score), strconv.Itoa(test.Percentage / 100 * test.Score),
		})
	}

	Columns := []table.Column{
		{Title: "Pb ID", Width: 5},
		{Title: "Test ID", Width: 11},
		{Title: "Verdict", Width: 22},
		{Title: "Time", Width: 8},
		{Title: "Memory", Width: 8},
		{Title: "%", Width: 3},
		{Title: "Max", Width: 3},
		{Title: "Obtained", Width: 8},
	}

	internal.RenderTable(Columns, Rows, 1)
}
