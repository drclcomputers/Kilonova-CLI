// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eiannone/keyboard"
	"github.com/spf13/cobra"
)

type Problem struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	SourceCredits string `json:"source_credits"`
	MaxScore      int    `json:"max_score"`
}

var searchCmd = &cobra.Command{
	Use:   "search [ID, NAME or all (all problems available)]",
	Short: "Search for problems by ID or name.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		searchProblems(args[0])
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

// TABLE

var chosenPb string = ""

var globalRows []table.Row

type modelsearch struct {
	table table.Model
}

func (m modelsearch) Init() tea.Cmd {
	return nil
}

func (m modelsearch) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "esc":
			return m, tea.Quit
		case "enter":
			return m.handleSelection()
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m modelsearch) handleSelection() (tea.Model, tea.Cmd) {
	selectedIndex := m.table.Cursor()
	selectedProblem := globalRows[selectedIndex]

	chosenPb = string(selectedProblem[0])

	return m, tea.Quit
}

func (m modelsearch) View() string {
	return lipgloss.NewStyle().Margin(1, 2).Render(m.table.View()) + "\n(Use ↑/↓ to navigate, 'q' to quit, 'enter' to get the statement)"
}

// search problems
type search struct {
	Data struct {
		Count    int `json:"count"`
		Problems []struct {
			Id             int    `json:"id"`
			Name           string `json:"name"`
			Source_Credits string `json:"source_credits"`
			Max_Score      int    `json:"max_score"`
		}
	} `json:"data"`
}

func searchProblems(problem_name string) {
	if problem_name == "all" {
		problem_name = ""
	}

	searchData := map[string]interface{}{
		"name_fuzzy": problem_name,
		"offset":     0,
	}

	jsonData, err := json.Marshal(searchData)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	body, err := makeRequest("POST", URL_SEARCH, bytes.NewBuffer(jsonData), "2")
	if err != nil {
		logErr(err)
	}

	var data search
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		os.Exit(1)
	}

	if data.Data.Count == 0 {
		fmt.Println("No problems found.")
		return
	}

	pages := data.Data.Count / 50
	if data.Data.Count%50 != 0 {
		pages++
	}

	var rows []table.Row

	for offset := 0; offset < data.Data.Count; offset += 50 {
		searchData["offset"] = offset

		jsonData, err := json.Marshal(searchData)
		if err != nil {
			logErr(err)
		}

		body, err = makeRequest("POST", URL_SEARCH, bytes.NewBuffer(jsonData), "2")
		if err != nil {
			logErr(err)
		}

		err = json.Unmarshal(body, &data)
		if err != nil {
			logErr(err)
		}

		for _, problem := range data.Data.Problems {
			if problem.Max_Score == -1 {
				problem.Max_Score = 0
			}
			if problem.Source_Credits == "" {
				problem.Source_Credits = "-"
			}
			rows = append(rows, table.Row{
				fmt.Sprintf("%d", problem.Id),
				problem.Name,
				problem.Source_Credits,
				fmt.Sprintf("%d", problem.Max_Score),
			})
		}
	}

	globalRows = rows

	columns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Name", Width: 20},
		{Title: "Source", Width: 40},
		{Title: "Max Score", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	t.SetStyles(table.DefaultStyles())

	p := tea.NewProgram(modelsearch{table: t}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}

	if chosenPb != "" {
		var choice string = ""
		fmt.Print("\nDo you wish to see the statement in RO(r) or EN(e): ")
		//fmt.Scanf("%s", &choice)\

		if err := keyboard.Open(); err != nil {
			logErr(err)
		}
		defer keyboard.Close()

		for choice == "" {
			key, _, err := keyboard.GetSingleKey()
			if err != nil {
				log.Fatal(err)
			}

			switch {
			case key == rune(keyboard.KeyEsc):
				choice = "ESC"
			case key == rune('r') || key == rune('R'):
				choice = "RO"
			case key == rune('e') || key == rune('E'):
				choice = "EN"
			default:
				choice = "ESC"
			}

		}

		if choice == "ESC" {
			return
		}

		printStatement(chosenPb, choice, 1)
	}
}
