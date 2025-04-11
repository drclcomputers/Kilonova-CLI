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
	"kncli/internal"
	"regexp"
	"strings"

	"text/template"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

var Online = false

var PrintStatementCmd = &cobra.Command{
	Use:   "statement [ID] [RO or EN (required for online)]",
	Short: "Print problem statement in chosen language.",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Println("Starting network services for online searching ...")
			_, _ = PrintStatement(args[0], args[1], 1)
			fmt.Println("Disabling network services for online searching ...")
		} else {
			_, _ = PrintStatement(args[0], "NO_LANG_CHOSEN", 1)
		}
	},
}

func init() {
	PrintStatementCmd.Flags().BoolVarP(&Online, "online", "o", false, "Get problem statement online.")
}

func formatText(DecodedText string) string {

	for Old, New := range internal.Replacements {
		DecodedText = strings.ReplaceAll(DecodedText, Old, New)
	}

	for _, Pattern := range internal.ReplacementsRegex {
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

// Problem Details

func GetProblemInfoStructOnline(ID string) (internal.ProblemInfo, error) {
	url := fmt.Sprintf(internal.URL_PROBLEM, ID)
	ResponseBody, err := internal.MakeGetRequest(url, nil, internal.RequestNone)
	if err != nil {
		internal.LogError(err)
	}

	var ProblemInfo internal.ProblemInfo
	if err := json.Unmarshal(ResponseBody, &ProblemInfo); err != nil {
		var res internal.RawKilonovaResponse
		if err := json.Unmarshal(ResponseBody, &res); err != nil {
			internal.LogError(err)
		}
		return internal.ProblemInfo{}, fmt.Errorf("no")
	}

	return ProblemInfo, nil
}

func GetProblemInfoStructLocal(ID string) (internal.ProblemInfo, error) {
	db := internal.DBOpen()
	defer db.Close()

	query := "SELECT id, name, timelimit, memorylimit, sourcesize, credits FROM problems\nWHERE CAST(id AS TEXT) LIKE $1;"

	var data internal.Problem
	_ = db.QueryRow(query, ID).Scan(&data.Id, &data.Name, &data.Time, &data.MemoryLimit, &data.SourceSize, &data.SourceCredits)

	fmt.Println(data)

	return internal.ProblemInfo{Data: data}, nil
}

func GetProblemInfoText(ID string) string {
	var ProblemInfo internal.ProblemInfo
	if Online {
		ProblemInfo, _ = GetProblemInfoStructOnline(ID)
	} else {
		ProblemInfo, _ = GetProblemInfoStructLocal(ID)
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

	if data.Credits == "" {
		data.Credits = "-"
	}

	TemplateCompleted, err := template.New("ProblemInfo").Parse(internal.TemplatePattern)
	if err != nil {
		internal.LogError(err)
		return ""
	}

	var Buffer bytes.Buffer
	if err := TemplateCompleted.Execute(&Buffer, data); err != nil {
		internal.LogError(err)
		return ""
	}

	return Buffer.String()
}

// Problem statement

func getStatementURL(id, lang string) (string, error) {
	switch strings.ToUpper(lang) {
	case "RO":
		return fmt.Sprintf(internal.URL_STATEMENT, id, internal.STAT_FILENAME_RO), nil
	case "EN":
		return fmt.Sprintf(internal.URL_STATEMENT, id, internal.STAT_FILENAME_EN), nil
	default:
		return "", fmt.Errorf("invalid language chosen: %q. Must be 'RO' or 'EN'", lang)
	}
}

func GetStatementOnline(ID, language string) string {
	url, err := getStatementURL(ID, language)
	if err != nil {
		internal.LogError(fmt.Errorf("error fetching URL: %w", err))
	}

	ResponseBody, err := internal.MakeGetRequest(url, nil, internal.RequestNone)
	if err != nil {
		internal.LogError(fmt.Errorf("error fetching statement: %w", err))
	}

	if strings.Contains(string(ResponseBody), "notfound") {
		return internal.NOLANG
	}

	var Statement Statement
	if err := json.Unmarshal(ResponseBody, &Statement); err != nil {
		internal.LogError(fmt.Errorf("failed to parse statement: %w", err))
	}

	return Statement.Data.Data
}

func GetStatementLocal(ID string) string {
	if !internal.DBExists() {
		internal.LogError(fmt.Errorf("problem database doesn't exist! Signin or run 'database create' "))
	}
	db := internal.DBOpen()
	defer db.Close()

	query := "SELECT statement FROM problems\nWHERE CAST(id AS TEXT) LIKE $1;"

	var statement string
	_ = db.QueryRow(query, ID).Scan(&statement)

	return statement
}

func PrintStatement(ID, language string, useCase int) (string, error) { // 1 - Print, 2 - Return text
	var statement string
	if Online {
		statement = GetStatementOnline(ID, language)
	} else {
		statement = GetStatementLocal(ID)
	}

	if statement == internal.NOLANG {
		if language == "RO" {
			internal.LogError(fmt.Errorf("statement not available in Romanian. Try again in English"))
		} else if language == "EN" {
			internal.LogError(fmt.Errorf("statement not available in English. Try again in Romanian"))
		} else {
			internal.LogError(fmt.Errorf("unknoun language chosen. Must be either RO or EN"))
		}
	}

	text, err := internal.DecodeBase64Text(statement)
	if err != nil {
		if useCase == 2 {
			internal.LogError(errors.New(internal.NOLANG))
		}
		internal.LogError(fmt.Errorf("failed to decode base64 text: %w", err))
	}

	DecodedText := formatText(text)
	if useCase == 2 {
		return DecodedText, nil
	}

	Rendered, err := renderStatement(ID, DecodedText)
	if err != nil {
		return "error", fmt.Errorf("failed to render statement: %w", err)
	}

	if err := runTUI(Rendered); err != nil {
		return "error", fmt.Errorf("failed to run TUI program: %w", err)
	}

	return DecodedText, nil
}

// Others

func renderStatement(ID, DecodedText string) (string, error) {
	ProblemInfoText := GetProblemInfoText(ID)
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
	p := tea.NewProgram(internal.NewTextModel(rendered))
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run TUI program: %w", err)
	}
	return nil
}
