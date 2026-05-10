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
			_, _ = PrintStatement(args[0], args[1], true, 1)
			fmt.Println("Disabling network services for online searching ...")
		} else if Online {
			_, _ = PrintStatement(args[0], "NO_LANG_CHOSEN", true, 1)
		} else {
			_, _ = PrintStatement(args[0], "", false, 1)
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
		return internal.ProblemInfo{}, fmt.Errorf("failed to fetch problem %s: %w", ID, err)
	}

	var ProblemInfo internal.ProblemInfo
	if err := json.Unmarshal(ResponseBody, &ProblemInfo); err != nil {
		return internal.ProblemInfo{}, fmt.Errorf("failed to parse problem info: %w", err)
	}

	return ProblemInfo, nil
}

func GetProblemInfoStructLocal(ID string) (internal.ProblemInfo, error) {
	db := internal.DBOpen()
	if db == nil {
		return internal.ProblemInfo{}, fmt.Errorf("problem database is not available")
	}
	defer internal.DBClose()

	query := "SELECT id, name, timelimit, memorylimit, sourcesize, credits FROM problems\nWHERE CAST(id AS TEXT) LIKE ?;"

	var data internal.Problem
	err := db.QueryRow(query, ID).Scan(&data.Id, &data.Name, &data.Time, &data.MemoryLimit, &data.SourceSize, &data.SourceCredits)
	if err != nil {
		return internal.ProblemInfo{}, fmt.Errorf("problem %s not found in database: %w", ID, err)
	}

	return internal.ProblemInfo{Data: data}, nil
}

func GetProblemInfoText(ID string, online bool) string {
	var ProblemInfo internal.ProblemInfo
	var err error
	if online {
		ProblemInfo, err = GetProblemInfoStructOnline(ID)
	} else {
		ProblemInfo, err = GetProblemInfoStructLocal(ID)
	}
	if err != nil {
		internal.LogError(fmt.Errorf("failed to get problem info: %w", err))
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

func GetStatementOnline(ID, language string, useCase int) string {
	return getStatementOnline(ID, language, useCase, false)
}

func GetStatementOnlineQuiet(ID, language string, useCase int) string {
	return getStatementOnline(ID, language, useCase, true)
}

func getStatementOnline(ID, language string, useCase int, suppressMissing bool) string {
	url, err := getStatementURL(ID, language)
	if err != nil {
		return internal.NOLANG
	}

	var ResponseBody []byte
	if useCase == 1 {
		ResponseBody, err = internal.MakeGetRequest(url, nil, internal.RequestNone)
	} else if useCase == 2 {
		ResponseBody, err = internal.MakeGetRequest(url, nil, internal.RequestDatabase)
	}
	if err != nil {
		if !suppressMissing || !isMissingStatementErr(err) {
			internal.LogError(fmt.Errorf("error fetching statement for problem %s: %w", ID, err))
		}
		return internal.NOLANG
	}

	if bytes.Contains(ResponseBody, []byte("notfound")) {
		return internal.NOLANG
	}

	var Statement Statement
	if err := json.Unmarshal(ResponseBody, &Statement); err != nil {
		internal.LogError(fmt.Errorf("failed to parse statement: %w", err))
	}

	return Statement.Data.Data
}

func isMissingStatementErr(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "attachment does not exist") || strings.Contains(msg, "notfound")
}

func GetStatementLocal(ID string) string {
	db := internal.DBOpen()
	if db == nil {
		return ""
	}
	defer internal.DBClose()

	query := "SELECT statement FROM problems\nWHERE CAST(id AS TEXT) LIKE ?;"

	var statement string
	if err := db.QueryRow(query, ID).Scan(&statement); err != nil {
		return ""
	}
	return statement
}

func PrintStatement(ID, language string, online bool, useCase int) (string, error) { // 1 - Print, 2 - Return text
	var statement string
	if online {
		statement = GetStatementOnline(ID, language, 1)
	} else {
		if !internal.DBExists() {
			internal.LogError(fmt.Errorf("problem database doesn't exist! Signin or run 'database create' "))
			return "", fmt.Errorf("problem database doesn't exist")
		}

		if internal.RefreshOrNotDB() {
			defer fmt.Println("Warning: You should refresh the database using 'database refresh' to get more problems.")
		}
		if !internal.ProblemExistsDB(ID) {
			fmt.Println("No problem with this ID found in the database.")
			return "", nil
		}

		statement = GetStatementLocal(ID)
	}

	if statement == internal.NOLANG {
		err := fmt.Errorf("statement not available in %s", language)
		if useCase == 2 {
			return "", errors.New(internal.NOLANG)
		}
		internal.LogError(err)
		return "", err
	}

	text, err := internal.DecodeBase64Text(statement)
	if err != nil {
		internal.LogError(fmt.Errorf("failed to decode base64 text: %w", err))
		return "", err
	}

	DecodedText := formatText(text)

	if useCase == 2 {
		return DecodedText, nil
	}

	Rendered, err := renderStatement(ID, DecodedText, online)
	if err != nil {
		return "error", fmt.Errorf("failed to render statement: %w", err)
	}

	if err := runTUI(Rendered); err != nil {
		return "error", fmt.Errorf("failed to run TUI program: %w", err)
	}

	return DecodedText, nil
}

// Others

func renderStatement(ID, DecodedText string, online bool) (string, error) {
	ProblemInfoText := GetProblemInfoText(ID, online)
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
