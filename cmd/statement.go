// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package cmd

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

var printStatementCmd = &cobra.Command{
	Use:   "statement [ID] [RO or EN]",
	Short: "Print problem statement in chosen language",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		printStatement(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(printStatementCmd)
}

func formatText(decodedtext string) string {
	replacements := map[string]string{
		"$":              "",
		" \\ ":           "",
		"\\ldots":        "...",
		"\\leq":          "≤",
		"\\geq":          "≥",
		"\\el":           "",
		"\\in":           "∈",
		"\\le":           "≤",
		"\\qe":           "≥",
		"\\pm":           "±",
		"\\cdot":         "•",
		"\\sum_":         "Σ ",
		"\\displaystyle": "",
		"\\times":        "x",
		"\\%":            "%",
	}

	for old, new := range replacements {
		decodedtext = strings.ReplaceAll(decodedtext, old, new)
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
		decodedtext = re.ReplaceAllString(decodedtext, "$1")
	}

	re := regexp.MustCompile(`~\[([^\]]+)\]`)
	decodedtext = re.ReplaceAllString(decodedtext, "$1 Download the assets to view images.")

	return decodedtext
}

// print statement
type info struct {
	Data struct {
		Name          string  `json:"name"`
		Time          float64 `json:"time_limit"`
		MemoryLimit   int     `json:"memory_limit"`
		SourceSize    int     `json:"source_size"`
		SourceCredits string  `json:"source_credits"`
	} `json:"data"`
}

type statement struct {
	Status string `json:"status"`
	Data   struct {
		Data string `json:"data"`
	} `json:"data"`
}

func problemInfo(id string) string {
	//info
	url := fmt.Sprintf(URL_PROBLEM, id)
	body, err := makeRequest("GET", url, nil, "0")
	if err != nil {
		logErr(err)
		os.Exit(1)
	}

	var info info
	if err := json.Unmarshal(body, &info); err != nil {
		logErr(err)
		os.Exit(1)
	}
	return fmt.Sprintf("Name: %s\n\nID: #%s\n\nTime Limit: %.2fs\n\nMemory Limit: %dKB\n\nSource Size: %dKB\n\nCredits: %s\n\n",
		info.Data.Name, id, info.Data.Time, info.Data.MemoryLimit,
		info.Data.SourceSize, info.Data.SourceCredits)

}

func printStatement(id, language string) {
	var url string

	//statement

	renderer, err := glamour.NewTermRenderer(glamour.WithStandardStyle("dark"))
	if err != nil {
		logErr(err)
		return
	}

	if language == "RO" {
		url = fmt.Sprintf(URL_STATEMENT, id, STAT_FILENAME_RO)
	} else {
		url = fmt.Sprintf(URL_STATEMENT, id, STAT_FILENAME_EN)
	}
	body, err := makeRequest("GET", url, nil, "0")
	if err != nil {
		logErr(err)
		return
	}

	if strings.Contains(string(body), `"status":"error"`) {
		log.Fatal("Error: Problem statement is not available in the chosen language!")
	}

	var data statement
	if err := json.Unmarshal(body, &data); err != nil {
		logErr(err)
		return
	}
	text, err := b64.StdEncoding.DecodeString(data.Data.Data)
	if err != nil {
		logErr(err)
		return
	}
	decodedtext := formatText(string(text))

	rendered, err := renderer.Render(problemInfo(id) + "\n# STATEMENT\n\n" + decodedtext)
	if err != nil {
		logErr(err)
		return
	}

	p := tea.NewProgram(newTextModel(rendered))
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
