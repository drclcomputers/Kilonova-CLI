// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package contests

import (
	"encoding/json"
	"fmt"
	utility "kilocli/cmd/utility"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/huh/spinner"
)

func downloadLeaderboard(contestID string) {
	resp, err := utility.MakeGetRequest(fmt.Sprintf(utility.URL_CONTEST_ASSETS, contestID), nil, utility.RequestDownloadZip)
	if err != nil {
		utility.LogError(err)
		return
	}

	homedir, err := os.Getwd()
	if err != nil {
		utility.LogError(fmt.Errorf("failed to get current working directory: %w", err))
		return
	}

	downFile := filepath.Join(homedir, "leaderboard_"+contestID+".csv")
	outFile, err := os.Create(downFile)
	if err != nil {
		utility.LogError(fmt.Errorf("failed to create file %q: %w", downFile, err))
		return
	}
	defer outFile.Close()

	if err := os.WriteFile(downFile, resp, 0644); err != nil {
		utility.LogError(fmt.Errorf("failed to write to file %q: %w", downFile, err))
		return
	}

	fmt.Printf("Leaderboard to contest #%s saved to %q\n", contestID, downFile)
}

func leaderboard(contestID string) {
	url := fmt.Sprintf(utility.URL_CONTEST_LEADERBOARD, contestID)
	body, err := utility.MakeGetRequest(url, nil, utility.RequestNone)
	if err != nil {
		utility.LogError(err)
		return
	}
	var data LeaderboardData
	if err = json.Unmarshal(body, &data); err != nil {
		utility.LogError(err)
		return
	}

	var Rows []table.Row

	for _, entry := range data.Data.Entries {
		var scores string
		for _, score := range entry.Scores {
			scores += fmt.Sprintf("%d  ", score)
		}
		Rows = append(Rows, table.Row{
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

	Columns := []table.Column{
		{Title: "ID | Name", Width: 25},
		{Title: problemNamesTitle, Width: 50},
		{Title: "Total", Width: 5},
	}

	utility.RenderTable(Columns, Rows, 1)

	if shouldDownload {
		action := func() { downloadLeaderboard(contestID) }
		if err := spinner.New().Title("Waiting for shouldDownload...").Action(action).Run(); err != nil {
			utility.LogError(fmt.Errorf("error during source code shouldDownload for submission #%s: %w", contestID, err))
			return
		}
	}
}
