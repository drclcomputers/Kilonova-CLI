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

	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

var DebugCmd = &cobra.Command{
	Use:   "debug [Submission ID]",
	Short: "Let Gemini analyze failed tests and suggest fixes. (online, requires GEMINI_API_KEY)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		action := func() { analyzeSubmission(args[0]) }
		if err := spinner.New().Title("Gemini is analyzing your submission...").Action(action).Run(); err != nil {
			internal.LogError(err)
		}
	},
}

type debugSubmissionData struct {
	Status string `json:"status"`
	Data   struct {
		ID             int     `json:"id"`
		ProblemID      int     `json:"problem_id"`
		Language       string  `json:"language"`
		Score          float64 `json:"score"`
		CompileError   bool    `json:"compile_error"`
		CompileMessage string  `json:"compile_message"`
		MaxTime        float64 `json:"max_time"`
		MaxMemory      int     `json:"max_memory"`
		Code           string  `json:"code"`
		Subtests       []struct {
			ID         int     `json:"id"`
			Done       bool    `json:"done"`
			Verdict    string  `json:"verdict"`
			Time       float64 `json:"time"`
			Memory     int     `json:"memory"`
			Percentage int     `json:"percentage"`
			Score      int     `json:"score"`
		} `json:"subtests"`
	} `json:"data"`
}

func analyzeSubmission(submissionID string) {
	url := fmt.Sprintf(internal.URL_LATEST_SUBMISSION, submissionID)
	body, err := internal.MakeGetRequest(url, nil, internal.RequestNone)
	if err != nil {
		internal.LogError(err)
		return
	}

	var sub debugSubmissionData
	if err := json.Unmarshal(body, &sub); err != nil {
		internal.LogError(fmt.Errorf("failed to parse submission: %w", err))
		return
	}

	if sub.Data.CompileError {
		fmt.Printf("\n❌ Compile Error:\n%s\n", sub.Data.CompileMessage)
		return
	}

	// Collect failed tests
	var failedTests []string
	var allTests []string
	for _, test := range sub.Data.Subtests {
		verdict := strings.TrimPrefix(test.Verdict, "translate:")
		entry := fmt.Sprintf(
			"Test #%d: %s | Time: %.3fs | Memory: %dKB | Score: %d/%d",
			test.ID, verdict, test.Time, test.Memory, test.Score, test.Percentage,
		)
		allTests = append(allTests, entry)
		if verdict != "OK" && verdict != "Accepted" {
			failedTests = append(failedTests, entry)
		}
	}

	// If no explicit failures but score < 100, include all tests
	if len(failedTests) == 0 {
		failedTests = allTests
	}

	systemPrompt := `You are an expert competitive programming debugger. 
Analyze the submission details and failed tests. Explain:
1. What went wrong (TLE, WA, RE, MLE — explain what each means)
2. Why the code likely fails for those tests
3. What approach would fix it (algorithmic fix, edge case handling, optimization)

Be specific and reference the failed test cases.`

	userPrompt := fmt.Sprintf(`Submission #%s for Problem #%d:
Score: %.0f | Language: %s | Max Time: %.3fs | Max Memory: %dKB

Failed Tests:
%s

Code (first 4000 chars):
%s

Analyze what went wrong and suggest fixes:`,
		submissionID, sub.Data.ProblemID, sub.Data.Score, sub.Data.Language,
		sub.Data.MaxTime, sub.Data.MaxMemory,
		strings.Join(failedTests, "\n"),
		truncateCode(sub.Data.Code),
	)

	analysis, err := internal.GeminiGenerate(systemPrompt, userPrompt)
	if err != nil {
		internal.LogError(err)
		return
	}

	fmt.Print("\n🔍 Debug Analysis:\n\n")
	fmt.Println(analysis)
}

func truncateCode(code string) string {
	if len(code) > 4000 {
		return code[:4000] + "\n... (truncated)"
	}
	return code
}
