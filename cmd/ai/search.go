// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package ai

import (
	"encoding/json"
	"fmt"
	"kncli/internal"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

var AISearchCmd = &cobra.Command{
	Use:   "find [natural language description]",
	Short: "Find problems using natural language (Gemini-powered). (online, requires GEMINI_API_KEY)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := strings.Join(args, " ")
		action := func() { aiSearchProblems(query) }
		if err := spinner.New().Title("Gemini is searching...").Action(action).Run(); err != nil {
			internal.LogError(err)
		}
	},
}

type allProblemsResponse struct {
	Status string `json:"status"`
	Data   []struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		SourceCredits string `json:"source_credits"`
		MaxScore      int    `json:"max_score"`
	} `json:"data"`
}

func aiSearchProblems(query string) {
	// Fetch all problems from the API
	url := fmt.Sprintf(internal.URL_PROBLEM, "get")
	resp, err := internal.PostJSON[allProblemsResponse](url, struct{}{})
	if err != nil {
		internal.LogError(err)
		return
	}

	if len(resp.Data) == 0 {
		fmt.Println("No problems found on Kilonova.")
		return
	}

	// Build a problem catalog for Gemini
	var catalog strings.Builder
	for i, p := range resp.Data {
		catalog.WriteString(fmt.Sprintf("%d. #%d - %s (Max Score: %d, Source: %s)\n",
			i+1, p.ID, p.Name, p.MaxScore, p.SourceCredits))
		if i > 200 {
			catalog.WriteString(fmt.Sprintf("\n... and %d more problems", len(resp.Data)-200))
			break
		}
	}

	systemPrompt := `You are a problem recommendation engine for a competitive programming platform.
Given a user's natural language query and a catalog of available problems, 
return the IDs of the 10 most relevant problems as a JSON array of integers.
Example: [42, 7, 133, 99]
Output ONLY the JSON array, nothing else.`

	userPrompt := fmt.Sprintf(`User query: "%s"

Available problems:
%s

Return the 10 most relevant problem IDs as JSON array [id1, id2, ...]:`, query, catalog.String())

	result, err := internal.GeminiGenerate(systemPrompt, userPrompt)
	if err != nil {
		internal.LogError(err)
		return
	}

	// Parse the JSON array from Gemini's response
	result = strings.TrimSpace(result)
	result = strings.TrimPrefix(result, "```json")
	result = strings.TrimPrefix(result, "```")
	result = strings.TrimSuffix(result, "```")
	result = strings.TrimSpace(result)

	var matchedIDs []int
	if err := json.Unmarshal([]byte(result), &matchedIDs); err != nil {
		// Fallback: try to extract numbers manually
		internal.LogError(fmt.Errorf("gemini returned invalid format, showing raw results:\n%s", result))
		return
	}

	// Build table rows from matched IDs
	var rows []table.Row

	for _, id := range matchedIDs {
		for _, p := range resp.Data {
			if p.ID == id {
				rows = append(rows, table.Row{
					fmt.Sprintf("%d", p.ID),
					p.Name,
					p.SourceCredits,
					fmt.Sprintf("%d", p.MaxScore),
				})
				break
			}
		}
	}

	if len(rows) == 0 {
		fmt.Println("No matching problems found. Try a different query.")
		return
	}

	columns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Name", Width: 30},
		{Title: "Source", Width: 35},
		{Title: "Max Score", Width: 10},
	}

	internal.RenderTable(columns, rows, 1)
	fmt.Printf("\n🔍 Found %d problems for: \"%s\"\n", len(rows), query)
}
