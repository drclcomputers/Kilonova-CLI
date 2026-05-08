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

var ExplainCmd = &cobra.Command{
	Use:   "explain [Problem ID]",
	Short: "Let Gemini explain the problem statement in simple terms. (online, requires GEMINI_API_KEY)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		action := func() { explainProblem(args[0]) }
		if err := spinner.New().Title("Gemini is analyzing...").Action(action).Run(); err != nil {
			internal.LogError(err)
		}
	},
}

func explainProblem(problemID string) {
	infoText, statementText, err := internal.GetProblemContext(problemID)
	if err != nil {
		internal.LogError(err)
		return
	}

	systemPrompt := `You are a patient, enthusiastic competitive programming tutor. 
Explain the problem clearly and simply. Break it down into:
1. 📖 WHAT the problem is asking (in simple terms)
2. 🔑 KEY insights needed to solve it
3. 📊 INPUT/OUTPUT format explained with an example
4. ⚠️ EDGE CASES to watch out for

Keep it beginner-friendly but precise. Use emojis for visual clarity.`

	userPrompt := fmt.Sprintf(`Problem Info:
%s

Problem Statement:
%s

Please explain this problem:`, infoText, truncateForGemini(statementText))

	explanation, err := internal.GeminiGenerate(systemPrompt, userPrompt)
	if err != nil {
		internal.LogError(err)
		return
	}

	fmt.Println("\n" + explanation)
}
