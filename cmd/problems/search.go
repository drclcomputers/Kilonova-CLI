// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package problems

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kncli/internal"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/eiannone/keyboard"
	"github.com/spf13/cobra"
)

var onlinesearch = false

var SearchCmd = &cobra.Command{
	Use:   "search [ID, NAME or all (all problems available)]",
	Short: "Search for problems by ID or name.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if onlinesearch {
			fmt.Println("Starting network services for online searching ...")
			searchProblemsOnline(args[0])
			fmt.Println("Disabling network services for online searching ...")
		} else {
			searchProblemsLocal(args[0])
		}
	},
}

func init() {
	SearchCmd.Flags().BoolVarP(&onlinesearch, "online", "o", false, "Online search for problems. May take longer.")
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

func fetchProblemsOnline(ProblemName string) ([]table.Row, error) {
	if ProblemName == "all" {
		ProblemName = ""
	}

	SearchData := map[string]interface{}{
		"name_fuzzy": ProblemName,
		"offset":     0,
	}

	var Rows []table.Row

	Data, err := doSearchOnline(SearchData)
	if err != nil {
		return nil, err
	}

	NumberOfPages := (Data.Data.Count + 49) / 50

	for Page := 0; Page < NumberOfPages; Page++ {
		SearchData["offset"] = Page * 50

		PageData, err := doSearchOnline(SearchData)
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

func doSearchOnline(searchData map[string]interface{}) (*SearchResponse, error) {
	payload, err := json.Marshal(searchData)
	if err != nil {
		return nil, fmt.Errorf("JSON marshal error: %w", err)
	}

	body, err := internal.MakePostRequest(internal.URL_SEARCH, bytes.NewBuffer(payload), internal.RequestJSON)
	if err != nil {
		return nil, fmt.Errorf("POST request failed: %w", err)
	}

	var res SearchResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("JSON unmarshal error: %w", err)
	}

	return &res, nil
}

func searchProblemsOnline(ProblemName string) {
	Rows, err := fetchProblemsOnline(ProblemName)
	if err != nil {
		internal.LogError(fmt.Errorf("error fetching problems: %v", err))
		return
	}

	if len(Rows) == 0 {
		fmt.Println("No problems found.")
		return
	}

	internal.GlobalRows = Rows

	Columns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Name", Width: 20},
		{Title: "Source", Width: 40},
		{Title: "Max Score", Width: 10},
	}

	internal.RenderTable(Columns, Rows, 2)

	if internal.ChosenProblem != "" {
		chooseLanguageAndShowStatement()
	}
}

func chooseLanguageAndShowStatement() {
	Online = true
	fmt.Print("\nDo you wish to see the statement in RO(r) or EN(e): ")

	if err := keyboard.Open(); err != nil {
		internal.LogError(err)
		return
	}
	defer keyboard.Close()

	for {
		key, _, err := keyboard.GetSingleKey()
		if err != nil {
			internal.LogError(fmt.Errorf("key read error: %w", err))
			return
		}

		switch key {
		case 'r', 'R':
			_, _ = PrintStatement(internal.ChosenProblem, "RO", 1)
			return
		case 'e', 'E':
			_, _ = PrintStatement(internal.ChosenProblem, "EN", 1)
			return
		case rune(keyboard.KeyEsc):
			return
		default:
			fmt.Print("Please press 'r' for RO or 'e' for EN (ESC to cancel): ")
		}
	}
}

func searchProblemsLocal(ProblemName string) {
	if !internal.DBExists() {
		internal.LogError(fmt.Errorf("problem database doesn't exist! Signin or run 'database create' "))
	}

	if internal.RefreshOrNotDB() {
		defer fmt.Println("Warning: You should refresh the database using 'database refresh' to get more problems.")
	}

	if ProblemName == "all" {
		ProblemName = ""
	}

	db := internal.DBOpen()
	defer internal.DBClose(db)

	var Rows []table.Row

	var pattern, query string

	pattern = "%" + ProblemName + "%"

	if _, err := internal.ValidateInt(ProblemName); err == nil {
		query = "SELECT id, name, credits\nFROM problems\nWHERE CAST(id AS TEXT) LIKE ?"
	} else {
		query = "SELECT id, name, credits\nFROM problems\nWHERE name LIKE ?;"
	}

	rows, err := db.Query(query, pattern)
	if err != nil {
		internal.LogError(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, credits string
		if err := rows.Scan(&id, &name, &credits); err != nil {
			internal.LogError(err)
			continue
		}

		Rows = append(Rows, table.Row{strconv.Itoa(id), name, credits, "Offline"})
	}

	if err := rows.Err(); err != nil {
		internal.LogError(err)
	}

	Columns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Name", Width: 20},
		{Title: "Source", Width: 40},
		{Title: "Max Score", Width: 10},
	}

	internal.GlobalRows = Rows

	internal.RenderTable(Columns, Rows, 2)

	if internal.ChosenProblem != "" {
		if onlinesearch {
			chooseLanguageAndShowStatement()
			return
		}
		_, _ = PrintStatement(internal.ChosenProblem, "null", 1)
	}
}
