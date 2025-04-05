// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

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

var chosenProblem string = ""

var globalRows []table.Row

type ModelSearch struct {
	table table.Model
}

func (m ModelSearch) Init() tea.Cmd {
	return nil
}

func (m ModelSearch) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m ModelSearch) handleSelection() (tea.Model, tea.Cmd) {
	selectedIndex := m.table.Cursor()
	selectedProblem := globalRows[selectedIndex]

	chosenProblem = string(selectedProblem[0])

	return m, tea.Quit
}

func (m ModelSearch) View() string {
	return lipgloss.NewStyle().Margin(1, 2).Render(m.table.View()) + "\n(Use ↑/↓ to navigate, 'q' to quit, 'enter' to get the statement)"
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

func fetchProblems(problemName string) ([]table.Row, error) {
	if problemName == "all" {
		problemName = ""
	}

	searchData := map[string]interface{}{
		"name_fuzzy": problemName,
		"offset":     0,
	}

	var rows []table.Row

	data, err := searchAPI(searchData)
	if err != nil {
		return nil, err
	}

	pages := (data.Data.Count + 49) / 50

	for page := 0; page < pages; page++ {
		searchData["offset"] = page * 50

		pageData, err := searchAPI(searchData)
		if err != nil {
			return nil, err
		}

		for _, problem := range pageData.Data.Problems {
			if problem.Max_Score == -1 {
				problem.Max_Score = 0
			}
			if problem.SourceCredits == "" {
				problem.SourceCredits = "-"
			}
			rows = append(rows, table.Row{
				fmt.Sprintf("%d", problem.Id),
				problem.Name,
				problem.SourceCredits,
				fmt.Sprintf("%d", problem.Max_Score),
			})
		}
	}

	return rows, nil
}

func searchAPI(searchData map[string]interface{}) (*Search, error) {
	jsonData, err := json.Marshal(searchData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	body, err := MakePostRequest(URL_SEARCH, bytes.NewBuffer(jsonData), RequestJSON)
	if err != nil {
		return nil, err
	}

	var data Search
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return &data, nil
}

func searchProblems(problemName string) {
	rows, err := fetchProblems(problemName)
	if err != nil {
		logError(fmt.Errorf("error fetching problems: %v", err))
		return
	}

	if len(rows) == 0 {
		fmt.Println("No problems found.")
		return
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

	p := tea.NewProgram(ModelSearch{table: t}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logError(fmt.Errorf("error running program: %v", err))
	}

	if chosenProblem != "" {
		chooseLanguageAndShowStatement()
	}
}

func chooseLanguageAndShowStatement() {
	var choice string
	fmt.Print("\nDo you wish to see the statement in RO(r) or EN(e): ")

	if err := keyboard.Open(); err != nil {
		logError(err)
	}
	defer keyboard.Close()

	for choice == "" {
		key, _, err := keyboard.GetSingleKey()
		if err != nil {
			logError(err)
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

	printStatement(chosenProblem, choice, 1)
}
