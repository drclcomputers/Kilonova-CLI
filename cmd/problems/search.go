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

	utility "kilocli/cmd/utility"

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

type Search struct {
	Data struct {
		Count    int `json:"count"`
		Problems []struct {
			Id            int    `json:"id"`
			Name          string `json:"name"`
			SourceCredits string `json:"source_credits"`
			Max_Score     int    `json:"max_score"`
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

	Data, err := search(SearchData)
	if err != nil {
		return nil, err
	}

	NumberOfPages := (Data.Data.Count + 49) / 50

	for Page := 0; Page < NumberOfPages; Page++ {
		SearchData["offset"] = Page * 50

		PageData, err := search(SearchData)
		if err != nil {
			return nil, err
		}

		for _, Problem := range PageData.Data.Problems {
			if Problem.Max_Score == -1 {
				Problem.Max_Score = 0
			}
			if Problem.SourceCredits == "" {
				Problem.SourceCredits = "-"
			}
			Rows = append(Rows, table.Row{
				strconv.Itoa(Problem.Id),
				Problem.Name,
				Problem.SourceCredits,
				strconv.Itoa(Problem.Max_Score),
			})
		}
	}

	return Rows, nil
}

func search(SearchData map[string]interface{}) (*Search, error) {
	JSONData, err := json.Marshal(SearchData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	ResponseBody, err := utility.MakePostRequest(utility.URL_SEARCH, bytes.NewBuffer(JSONData), utility.RequestJSON)
	if err != nil {
		return nil, err
	}

	var Data Search
	if err := json.Unmarshal(ResponseBody, &Data); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return &Data, nil
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
	var LanguageChoice string
	fmt.Print("\nDo you wish to see the statement in RO(r) or EN(e): ")

	if err := keyboard.Open(); err != nil {
		utility.LogError(err)
		return
	}
	defer keyboard.Close()

	for LanguageChoice == "" {
		Key, _, err := keyboard.GetSingleKey()
		if err != nil {
			utility.LogError(err)
			return
		}

		switch {
		case Key == rune(keyboard.KeyEsc):
			LanguageChoice = "ESC"
		case Key == rune('r') || Key == rune('R'):
			LanguageChoice = "RO"
		case Key == rune('e') || Key == rune('E'):
			LanguageChoice = "EN"
		default:
			LanguageChoice = "ESC"
		}
	}

	if LanguageChoice == "ESC" {
		return
	}

	PrintStatement(utility.ChosenProblem, LanguageChoice, 1)
}
