// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package ai

import (
	"fmt"
	"kncli/internal"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

var solveLang string

var SolveCmd = &cobra.Command{
	Use:   "solve [Problem ID]",
	Short: "Let Gemini AI generate a solution for the problem. (online, requires GEMINI_API_KEY)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		action := func() { generateSolution(args[0]) }
		if err := spinner.New().Title("Gemini is thinking...").Action(action).Run(); err != nil {
			internal.LogError(err)
		}
	},
}

func init() {
	SolveCmd.Flags().StringVarP(&solveLang, "lang", "l", "cpp", "Target language (c, cpp, python3, go, rust, java, nodejs)")
}

func generateSolution(problemID string) {
	infoText, statementText, err := internal.GetProblemContext(problemID)
	if err != nil {
		internal.LogError(err)
		return
	}

	systemPrompt := `Act as a World-Class Competitive Programming Engine.
Task: Generate a high-performance, correct solution.
Requirements:
- Language: Strictly adhere to the requested language.
- I/O: standard stdin/stdout.
- Performance: Optimal time and space complexity.
- Completeness: All necessary imports/headers.
- Constraints: Handle all edge cases and limits.
- Format: Output ONLY the raw source code. NO markdown code fences, NO prose, NO explanations. 
- Documentation: In-code comments only for complex logic.`

	userPrompt := fmt.Sprintf(`Problem Info:
%s

Problem Statement:
%s

Language: %s

Write a complete solution:`, infoText, truncateForGemini(statementText), mapLang(solveLang))

	solution, err := internal.GeminiGenerate(systemPrompt, userPrompt)
	if err != nil {
		internal.LogError(err)
		return
	}

	// Clean up code fences if Gemini adds them
	solution = strings.TrimSpace(solution)
	solution = strings.TrimPrefix(solution, "```"+mapLang(solveLang))
	solution = strings.TrimPrefix(solution, "```")
	solution = strings.TrimSuffix(solution, "```")
	solution = strings.TrimSpace(solution)

	filename := fmt.Sprintf("solution_%s.%s", problemID, getExt(solveLang))
	if err := os.WriteFile(filename, []byte(solution), 0644); err != nil {
		internal.LogError(fmt.Errorf("failed to write solution: %w", err))
		return
	}

	absPath, _ := filepath.Abs(filename)
	fmt.Printf("\n✅ Solution saved to: %s\n", absPath)
}

func mapLang(lang string) string {
	switch lang {
	case "python3", "py":
		return "python"
	case "nodejs", "js":
		return "javascript"
	case "cpp11", "cpp14", "cpp17", "cpp20":
		return "cpp"
	default:
		return lang
	}
}

func getExt(lang string) string {
	switch lang {
	case "python3", "py":
		return "py"
	case "nodejs", "js":
		return "js"
	case "go", "golang":
		return "go"
	case "rust":
		return "rs"
	case "java":
		return "java"
	case "c":
		return "c"
	default:
		return "cpp"
	}
}

func truncateForGemini(text string) string {
	// Gemini 3 Flash has 1M context, but keep prompts reasonable
	if len(text) > 32000 {
		return text[:32000] + "\n... (truncated)"
	}
	return text
}
