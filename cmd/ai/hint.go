// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package ai

import (
	"fmt"
	"kncli/internal"

	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

var hintLevel int

var HintCmd = &cobra.Command{
	Use:   "hint [Problem ID]",
	Short: "Get progressive hints without seeing the full solution. (online, requires GEMINI_API_KEY)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		action := func() { getHint(args[0], hintLevel) }
		if err := spinner.New().Title("Gemini is preparing a hint...").Action(action).Run(); err != nil {
			internal.LogError(err)
		}
	},
}

func init() {
	HintCmd.Flags().IntVarP(&hintLevel, "level", "l", 1, "Hint level (1=gentle nudge, 2=approach hint, 3=almost solution)")
}

func getHint(problemID string, level int) {
	if level < 1 {
		level = 1
	}
	if level > 3 {
		level = 3
	}

	infoText, statementText, err := internal.GetProblemContext(problemID)
	if err != nil {
		internal.LogError(err)
		return
	}

	hintInstructions := map[int]string{
		1: `Give a VERY GENTLE hint. Do NOT reveal the solution. 
Suggest what topic/technique to think about (e.g., "Think about using a prefix sum array" or "Consider sorting first").
Do NOT give any code or algorithm steps.`,
		2: `Give a MEDIUM hint. Describe the approach at a high level.
Mention the algorithm/data structure to use, but do NOT write any code.
Give a rough outline of the steps without implementation details.`,
		3: `Give a STRONG hint, almost the solution. 
Describe the algorithm in detail, mention key implementation details and edge cases.
But still do NOT write complete code — just describe what needs to happen.`,
	}

	systemPrompt := fmt.Sprintf(`You are a competitive programming coach. 
%s
Use emojis to make hints friendly. Label the hint level clearly.
REMEMBER: Never reveal the full solution or write complete code.`, hintInstructions[level])

	userPrompt := fmt.Sprintf(`Problem Info:
%s

Problem Statement:
%s

Give a Level %d hint for this problem:`, infoText, truncateForGemini(statementText), level)

	hint, err := internal.GeminiGenerate(systemPrompt, userPrompt)
	if err != nil {
		internal.LogError(err)
		return
	}

	fmt.Printf("\n💡 Hint (Level %d/%s):\n\n%s\n",
		level, map[int]string{1: "3 — Gentle", 2: "3 — Medium", 3: "3 — Strong"}[level], hint)

	if level < 3 {
		fmt.Printf("\n💪 Stuck? Try: kncli hint %s --level %d\n", problemID, level+1)
	}
}
