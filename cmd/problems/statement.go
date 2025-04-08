// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package problems

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"text/template"

	utility "kilocli/cmd/utility"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

var PrintStatementCmd = &cobra.Command{
	Use:   "statement [ID] [RO or EN]",
	Short: "Print problem statement in chosen language",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		PrintStatement(args[0], args[1], 1)
	},
}

func init() {
}

func formatText(DecodedText string) string {

	for Old, New := range utility.Replacements {
		DecodedText = strings.ReplaceAll(DecodedText, Old, New)
	}

	for _, Pattern := range utility.ReplacementsRegex {
		Regexp := regexp.MustCompile(Pattern)
		DecodedText = Regexp.ReplaceAllString(DecodedText, "$1")
	}

	Regexp := regexp.MustCompile(`~\[([^\]]+)\]`)
	DecodedText = Regexp.ReplaceAllString(DecodedText, "$1 Download the assets to view images.")

	return DecodedText
}

type Statement struct {
	Status string `json:"status"`
	Data   struct {
		Data string `json:"data"`
	} `json:"data"`
}

func GetProblemInfo(ID string) string {
	url := fmt.Sprintf(utility.URL_PROBLEM, ID)
	ResponseBody, err := utility.MakeGetRequest(url, nil, utility.RequestNone)
	if err != nil {
		utility.LogError(err)
		return ""
	}

	var ProblemInfo utility.ProblemInfo
	if err := json.Unmarshal(ResponseBody, &ProblemInfo); err != nil {
		utility.LogError(err)
		return ""
	}

	data := struct {
		Name        string
		ID          string
		TimeLimit   float64
		MemoryLimit int
		SourceSize  int
		Credits     string
	}{
		Name:        ProblemInfo.Data.Name,
		ID:          ID,
		TimeLimit:   ProblemInfo.Data.Time,
		MemoryLimit: ProblemInfo.Data.MemoryLimit,
		SourceSize:  ProblemInfo.Data.SourceSize,
		Credits:     ProblemInfo.Data.SourceCredits,
	}

	TemplateCompleted, err := template.New("ProblemInfo").Parse(utility.TemplatePattern)
	if err != nil {
		utility.LogError(err)
		return ""
	}

	var Buffer bytes.Buffer
	if err := TemplateCompleted.Execute(&Buffer, data); err != nil {
		utility.LogError(err)
		return ""
	}

	return Buffer.String()
}

func PrintStatement(ID, language string, useCase int) (string, error) {
	url, err := getStatementURL(ID, language)
	if err != nil {
		if useCase == 2 {
			return "", errors.New(utility.NOLANG)
		}
		return "", fmt.Errorf("error fetching URL: %w", err)
	}

	ResponseBody, err := utility.MakeGetRequest(url, nil, utility.RequestNone)
	if err != nil {
		if useCase == 2 {
			return "", errors.New(utility.NOLANG)
		}
		return "", fmt.Errorf("error fetching statement: %w", err)
	}

	if strings.Contains(string(ResponseBody), `"status":utility.ERROR`) {
		return "", errors.New("problem statement not available in the chosen language")
	}

	var Statement Statement
	if err := json.Unmarshal(ResponseBody, &Statement); err != nil {
		return "", fmt.Errorf("failed to parse statement: %w", err)
	}

	text, err := utility.DecodeBase64Text(Statement.Data.Data)
	if err != nil {
		if useCase == 2 {
			return "", errors.New(utility.NOLANG)
		}
		return "", fmt.Errorf("failed to decode base64 text: %w", err)
	}

	DecodedText := formatText(text)
	if useCase == 2 {
		//fmt.Println(DecodedText)
		return DecodedText, nil
	}

	Rendered, err := renderStatement(ID, DecodedText)
	if err != nil {
		return "", fmt.Errorf("failed to render statement: %w", err)
	}

	if err := runTUI(Rendered); err != nil {
		return "", fmt.Errorf("failed to run TUI program: %w", err)
	}

	return DecodedText, nil
}

func getStatementURL(ID, Language string) (string, error) {
	if Language == "RO" {
		return fmt.Sprintf(utility.URL_STATEMENT, ID, utility.STAT_FILENAME_RO), nil
	} else if Language == "EN" {
		return fmt.Sprintf(utility.URL_STATEMENT, ID, utility.STAT_FILENAME_EN), nil
	}
	return "", errors.New("invalid language choice, must be 'RO' or 'EN'")
}

func renderStatement(ID, DecodedText string) (string, error) {
	ProblemInfoText := GetProblemInfo(ID)
	if ProblemInfoText == "" {
		return "", errors.New("failed to retrieve problem information")
	}

	Renderer, err := glamour.NewTermRenderer(glamour.WithStandardStyle("dark"))
	if err != nil {
		return "", fmt.Errorf("failed to create renderer: %w", err)
	}

	Rendered, err := Renderer.Render(ProblemInfoText + "\n# STATEMENT\n\n" + DecodedText)
	if err != nil {
		return "", fmt.Errorf("failed to render statement: %w", err)
	}

	return Rendered, nil
}

func runTUI(rendered string) error {
	p := tea.NewProgram(utility.NewTextModel(rendered))
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run TUI program: %w", err)
	}
	return nil
}
