// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package cmd

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"text/template"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

var printStatementCmd = &cobra.Command{
	Use:   "statement [ID] [RO or EN]",
	Short: "Print problem statement in chosen language",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		printStatement(args[0], args[1], 1)
	},
}

func init() {
	rootCmd.AddCommand(printStatementCmd)
}

func formatText(decodedText string) string {
	replacements := map[string]string{
		"$":             "",
		` \ `:           "",
		`\ldots`:        "...",
		`\leq`:          "≤",
		`\geq`:          "≥",
		`\el`:           "",
		`\in`:           "∈",
		`\le`:           "≤",
		`\qe`:           "≥",
		`\pm`:           "±",
		`\cdot`:         "•",
		`\sum_`:         "Σ ",
		`\displaystyle`: "",
		`\times`:        "x",
		`\%`:            "%",
	}

	for old, new := range replacements {
		decodedText = strings.ReplaceAll(decodedText, old, new)
	}

	replacementsRegex := []string{
		`\\text{(.*?)}`,
		`\\texttt{(.*?)}`,
		`\\bm{(.*?)}`,
		`\\textit{(.*?)}`,
		`\\rule\{[^}]+\}\{[^}]+\}`,
	}

	for _, pattern := range replacementsRegex {
		re := regexp.MustCompile(pattern)
		decodedText = re.ReplaceAllString(decodedText, "$1")
	}

	re := regexp.MustCompile(`~\[([^\]]+)\]`)
	decodedText = re.ReplaceAllString(decodedText, "$1 Download the assets to view images.")

	return decodedText
}

type ProblemInfo struct {
	Data struct {
		Name          string  `json:"name"`
		Time          float64 `json:"time_limit"`
		MemoryLimit   int     `json:"memory_limit"`
		SourceSize    int     `json:"source_size"`
		SourceCredits string  `json:"source_credits"`
	} `json:"data"`
}

type Statement struct {
	Status string `json:"status"`
	Data   struct {
		Data string `json:"data"`
	} `json:"data"`
}

func problemInfo(id string) string {
	url := fmt.Sprintf(URL_PROBLEM, id)
	body, err := MakeGetRequest(url, nil, RequestNone)
	if err != nil {
		logError(err)
		return ""
	}

	var info ProblemInfo
	if err := json.Unmarshal(body, &info); err != nil {
		logError(err)
		return ""
	}

	tmpl := `Name: {{.Name}}
ID: #{{.ID}}
Time Limit: {{.TimeLimit}}s
Memory Limit: {{.MemoryLimit}}KB
Source Size: {{.SourceSize}}KB
Credits: {{.Credits}}
`
	data := struct {
		Name        string
		ID          string
		TimeLimit   float64
		MemoryLimit int
		SourceSize  int
		Credits     string
	}{
		Name:        info.Data.Name,
		ID:          id,
		TimeLimit:   info.Data.Time,
		MemoryLimit: info.Data.MemoryLimit,
		SourceSize:  info.Data.SourceSize,
		Credits:     info.Data.SourceCredits,
	}

	t, err := template.New("problemInfo").Parse(tmpl)
	if err != nil {
		logError(err)
		return ""
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		logError(err)
		return ""
	}

	return buf.String()
}

func printStatement(id, language string, useCase int) (string, error) {
	url, err := getStatementURL(id, language)
	if err != nil {
		if useCase == 2 {
			return "", errors.New("nolang")
		}
		return "", fmt.Errorf("error fetching URL: %w", err)
	}

	body, err := MakeGetRequest(url, nil, RequestNone)
	if err != nil {
		if useCase == 2 {
			return "", errors.New("nolang")
		}
		return "", fmt.Errorf("error fetching statement: %w", err)
	}

	if strings.Contains(string(body), `"status":"error"`) {
		return "", errors.New("problem statement not available in the chosen language")
	}

	var statement Statement
	if err := json.Unmarshal(body, &statement); err != nil {
		return "", fmt.Errorf("failed to parse statement: %w", err)
	}

	text, err := decodeBase64Text(statement.Data.Data)
	if err != nil {
		if useCase == 2 {
			return "", errors.New("nolang")
		}
		return "", fmt.Errorf("failed to decode base64 text: %w", err)
	}

	decodedText := formatText(text)
	if useCase == 2 {
		fmt.Println(decodedText)
		return decodedText, nil
	}

	rendered, err := renderStatement(id, decodedText)
	if err != nil {
		return "", fmt.Errorf("failed to render statement: %w", err)
	}

	if err := runTUI(rendered); err != nil {
		return "", fmt.Errorf("failed to run TUI program: %w", err)
	}

	return decodedText, nil
}

func getStatementURL(id, language string) (string, error) {
	if language == "RO" {
		return fmt.Sprintf(URL_STATEMENT, id, STAT_FILENAME_RO), nil
	} else if language == "EN" {
		return fmt.Sprintf(URL_STATEMENT, id, STAT_FILENAME_EN), nil
	}
	return "", errors.New("invalid language choice, must be 'RO' or 'EN'")
}

func decodeBase64Text(encodedText string) (string, error) {
	decodedText, err := b64.StdEncoding.DecodeString(encodedText)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 text: %w", err)
	}
	return string(decodedText), nil
}

func renderStatement(id, decodedText string) (string, error) {
	problemInfoText := problemInfo(id)
	if problemInfoText == "" {
		return "", errors.New("failed to retrieve problem information")
	}

	renderer, err := glamour.NewTermRenderer(glamour.WithStandardStyle("dark"))
	if err != nil {
		return "", fmt.Errorf("failed to create renderer: %w", err)
	}

	rendered, err := renderer.Render(problemInfoText + "\n# STATEMENT\n\n" + decodedText)
	if err != nil {
		return "", fmt.Errorf("failed to render statement: %w", err)
	}

	return rendered, nil
}

func runTUI(rendered string) error {
	p := tea.NewProgram(NewTextModel(rendered))
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run TUI program: %w", err)
	}
	return nil
}
