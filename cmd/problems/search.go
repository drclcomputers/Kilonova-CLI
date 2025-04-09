// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package problems

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	utility "kncli/cmd/utility"

	"github.com/charmbracelet/bubbles/table"
	"github.com/eiannone/keyboard"
	"github.com/spf13/cobra"
)

type Problem struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	SourceCredits string `json:"source_credits"`
	MaxScore      int    `json:"max_score"`
}

var SearchCmd = &cobra.Command{
	Use:   "search [ID, NAME or all (all problems available)]",
	Short: "Search for problems by ID or name.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		searchProblems(args[0])
	},
}

func init() {
}

type SearchResponse struct {
	Status string `json:"status"`
	Data   struct {
		Count    int `json:"count"`
		Problems []struct {
			Id            int    `json:"id"`
			Name          string `json:"name"`
			SourceCredits string `json:"source_credits"`
			MaxScore      int    `json:"max_score"`
		}
	} `json:"data"`
}

func fetchProblems(ProblemName string) ([]table.Row, error) {
	if ProblemName == "all" {
		ProblemName = ""
	}

	SearchData := map[string]interface{}{
		"name_fuzzy": ProblemName,
		"offset":     0,
	}

	var Rows []table.Row

	Data, err := doSearch(SearchData)
	if err != nil {
		return nil, err
	}

	NumberOfPages := (Data.Data.Count + 49) / 50

	for Page := 0; Page < NumberOfPages; Page++ {
		SearchData["offset"] = Page * 50

		PageData, err := doSearch(SearchData)
		if err != nil {
			return nil, err
		}

		for _, Problem := range PageData.Data.Problems {
			if Problem.MaxScore == -1 {
				Problem.MaxScore = 0
			}
			if Problem.SourceCredits == "" {
				Problem.SourceCredits = "-"
			}
			Rows = append(Rows, table.Row{
				strconv.Itoa(Problem.Id),
				Problem.Name,
				Problem.SourceCredits,
				strconv.Itoa(Problem.MaxScore),
			})
		}
	}

	return Rows, nil
}

func doSearch(searchData map[string]interface{}) (*SearchResponse, error) {
	payload, err := json.Marshal(searchData)
	if err != nil {
		return nil, fmt.Errorf("JSON marshal error: %w", err)
	}

	body, err := utility.MakePostRequest(utility.URL_SEARCH, bytes.NewBuffer(payload), utility.RequestJSON)
	if err != nil {
		return nil, fmt.Errorf("POST request failed: %w", err)
	}

	var res SearchResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("JSON unmarshal error: %w", err)
	}

	return &res, nil
}

func searchProblems(ProblemName string) {
	Rows, err := fetchProblems(ProblemName)
	if err != nil {
		utility.LogError(fmt.Errorf("error fetching problems: %v", err))
		return
	}

	if len(Rows) == 0 {
		fmt.Println("No problems found.")
		return
	}

	utility.GlobalRows = Rows

	Columns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Name", Width: 20},
		{Title: "Source", Width: 40},
		{Title: "Max Score", Width: 10},
	}

	utility.RenderTable(Columns, Rows, 2)

	if utility.ChosenProblem != "" {
		chooseLanguageAndShowStatement()
	}
}

func chooseLanguageAndShowStatement() {
	fmt.Print("\nDo you wish to see the statement in RO(r) or EN(e): ")

	if err := keyboard.Open(); err != nil {
		utility.LogError(err)
		return
	}
	defer keyboard.Close()

	for {
		key, _, err := keyboard.GetSingleKey()
		if err != nil {
			utility.LogError(fmt.Errorf("key read error: %w", err))
			return
		}

		switch key {
		case 'r', 'R':
			PrintStatement(utility.ChosenProblem, "RO", 1)
			return
		case 'e', 'E':
			PrintStatement(utility.ChosenProblem, "EN", 1)
			return
		case rune(keyboard.KeyEsc):
			return
		default:
			fmt.Print("Please press 'r' for RO or 'e' for EN (ESC to cancel): ")
		}
	}
}
