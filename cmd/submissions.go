// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package cmd

import (
	"bytes"
	"embed"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	u "net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zyedidia/highlight"
)

//go:embed highlight/*.yaml
var highlightDir embed.FS

var download bool = false

var uploadCodeCmd = &cobra.Command{
	Use:   "submit [ID] [LANGUAGE] [FILENAME]",
	Short: "Submit solution to problem.",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		uploadCode(args[0], args[1], args[2])
	},
}

var checkLangsCmd = &cobra.Command{
	Use:   "langs [ID]",
	Short: "View available languages for solutions.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checkLanguages(args[0], 1)
	},
}

var printSubmissionsCmd = &cobra.Command{
	Use:   "submissions [Problem ID or all (all problems)] [User ID, me (personal submissions), all (all users)] [1st page] [last page]",
	Short: "View sent submissions to a problem.",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		firstPage, err := strconv.Atoi(args[2])
		if err != nil {
			logError(fmt.Errorf("invalid first page number: %v", err))
			return
		}

		lastPage, err := strconv.Atoi(args[3])
		if err != nil {
			logError(fmt.Errorf("invalid last page number: %v", err))
			return
		}

		printSubmissions(args[0], args[1], firstPage, lastPage)
	},
}

var printSubmInfoCmd = &cobra.Command{
	Use:   "submissioninfo [Submission ID]",
	Short: "View a detailed description of a sent submission.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		printDetailsSubmission(args[0])
	},
}

func init() {
	rootCmd.AddCommand(checkLangsCmd)
	rootCmd.AddCommand(uploadCodeCmd)
	rootCmd.AddCommand(printSubmissionsCmd)
	rootCmd.AddCommand(printSubmInfoCmd)

	printSubmInfoCmd.Flags().BoolVarP(&download, "download_source", "d", true, "Download the source code of a submission.")
}

// print submissions

type SubmissionData struct {
	UserID         int     `json:"user_id"`
	ProblemID      int     `json:"problem_id"`
	Id             int     `json:"id"`
	CreatedAt      string  `json:"created_at"`
	Language       string  `json:"language"`
	Score          float64 `json:"score"`
	MaxMemory      int     `json:"max_memory"`
	MaxTime        float64 `json:"max_time"`
	CompileError   bool    `json:"compile_error"`
	CompileMessage string  `json:"compile_message"`
	Code           string  `json:"code,omitempty"`
}

type SubmissionList struct {
	Data struct {
		Submissions []SubmissionData `json:"submissions"`
		Count       int              `json:"count"`
	} `json:"data"`
}

type SubmissionDetails struct {
	Status string         `json:"status"`
	Data   SubmissionData `json:"data"`
}

func getSubmissionURL(userId, problemId string, offset int) string {
	switch {
	case userId == "all" && problemId == "all":
		return fmt.Sprintf(URL_SUBMISSION_LIST_NO_FILTER, offset)
	case userId == "all":
		return fmt.Sprintf(URL_SUBMISSION_LIST_NO_USER, offset, problemId)
	case problemId == "all":
		return fmt.Sprintf(URL_SUBMISSION_LIST_NO_PROBLEM, offset, userId)
	default:
		return fmt.Sprintf(URL_SUBMISSION_LIST, offset, problemId, userId)
	}
}

func parseSubmissionTime(timeStr string) (string, error) {
	parsedTime, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse time %q: %w", timeStr, err)
	}
	return parsedTime.Format("2006-01-02 15:04:05"), nil
}

func printSubmissions(problemId, userId string, firstPage, lastPage int) {
	if userId == "me" {
		userId = getUserID()
	}

	if firstPage <= 0 || lastPage <= 0 {
		logError(fmt.Errorf("invalid pages: both firstPage and lastPage must be positive integers"))
		return
	}

	if firstPage > lastPage {
		logError(fmt.Errorf("firstPage cannot be greater than lastPage"))
		return
	}

	var datasub SubmissionList
	var rows []table.Row
	var count = -1
	startOffset := (firstPage - 1) * 50
	endOffset := max((lastPage-1)*50, 50)

	for offset := max(startOffset, 0); offset < count && offset < endOffset; offset += 50 {
		url := getSubmissionURL(userId, problemId, offset)

		body, err := MakeGetRequest(url, nil, RequestFormAuth)
		if err != nil {
			logError(err)
			continue
		}

		if err := json.Unmarshal(body, &datasub); err != nil {
			logError(err)
			continue
		}

		count = datasub.Data.Count

		for _, problem := range datasub.Data.Submissions {
			formattedTime, err := parseSubmissionTime(problem.CreatedAt)
			if err != nil {
				logError(err)
				continue
			}

			rows = append(rows, table.Row{
				fmt.Sprintf("%d", problem.ProblemID),
				fmt.Sprintf("%d", problem.UserID),
				fmt.Sprintf("%d", problem.Id),
				formattedTime,
				problem.Language,
				fmt.Sprintf("%.0f", problem.Score),
			})
		}
	}

	columns := []table.Column{
		{Title: "Pb ID", Width: 5},
		{Title: "User ID", Width: 7},
		{Title: "Submission ID", Width: 12},
		{Title: "Time", Width: 25},
		{Title: "Language", Width: 10},
		{Title: "Score", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	t.SetStyles(table.DefaultStyles())

	p := tea.NewProgram(&Model{table: t}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logError(fmt.Errorf("error running program: %w", err))
	}
}

func downloadSource(submissionId, code string) {
	homedir, err := os.Getwd()
	if err != nil {
		logError(fmt.Errorf("failed to get current working directory: %w", err))
		return
	}

	downFile := filepath.Join(homedir, "source_"+submissionId+".txt")
	if err := os.WriteFile(downFile, []byte(code), 0644); err != nil {
		logError(fmt.Errorf("failed to write source code to file %q: %w", downFile, err))
		return
	}

	fmt.Printf("Source code for submission #%s saved to %q\n", submissionId, downFile)
}

func printDetailsSubmission(submissionId string) {
	var details SubmissionDetails

	id, err := strconv.Atoi(submissionId)
	if err != nil {
		logError(fmt.Errorf("invalid submission ID %q: %w", submissionId, err))
		return
	}

	url := fmt.Sprintf(URL_LATEST_SUBMISSION, id)

	formData := u.Values{
		"id": {submissionId},
	}

	body, err := MakeGetRequest(url, bytes.NewBufferString(formData.Encode()), RequestNone)
	if err != nil {
		logError(fmt.Errorf("error fetching submission details: %w", err))
		return
	}

	if err := json.Unmarshal(body, &details); err != nil {
		logError(fmt.Errorf("error unmarshalling response: %w", err))
		return
	}

	if details.Status != "success" {
		logError(fmt.Errorf("submission fetch failed with status: %v", details.Status))
		return
	}

	formattedTime, err := parseSubmissionTime(details.Data.CreatedAt)
	if err != nil {
		logError(err)
		return
	}

	problemId := strconv.Itoa(details.Data.ProblemID)
	fmt.Println(problemInfo(problemId))

	code, err := b64.StdEncoding.DecodeString(details.Data.Code)
	if err != nil {
		logError(fmt.Errorf("error decoding source code: %w", err))
		return
	}

	printSubmissionDetails(details, formattedTime, code)

	if download {
		action := func() { downloadSource(submissionId, string(code)) }
		if err := spinner.New().Title("Waiting for download...").Action(action).Run(); err != nil {
			logError(fmt.Errorf("error during source code download for submission #%s: %w", submissionId, err))
		}
	}
}

type SubmissionDetailsTemplate struct {
	ID             int
	CreatedAt      string
	Language       string
	Score          float64
	MaxMemory      int
	MaxTime        float64
	CompileError   bool
	CompileMessage string
	Code           string
}

func printSubmissionDetails(details SubmissionDetails, formattedTime string, code []byte) {
	submissionData := SubmissionDetailsTemplate{
		ID:             details.Data.Id,
		CreatedAt:      formattedTime,
		Language:       details.Data.Language,
		Score:          details.Data.Score,
		MaxMemory:      details.Data.MaxMemory,
		MaxTime:        details.Data.MaxTime,
		CompileError:   details.Data.CompileError,
		CompileMessage: details.Data.CompileMessage,
		Code:           formatCodeOutput(string(code), details.Data.Language),
	}

	const submissionTemplate = `
Submission ID: #{{.ID}}
Created: {{.CreatedAt}}
Language: {{.Language}}
Score: {{.Score}}

Max memory: {{.MaxMemory}}KB
Max time: {{.MaxTime}}s
Compile error: {{.CompileError}}
Compile message: {{.CompileMessage}}

Code:
{{.Code}}
`

	tmpl, err := template.New("submissionDetails").Parse(submissionTemplate)
	if err != nil {
		logError(fmt.Errorf("failed to parse template: %w", err))
		return
	}

	if err := tmpl.Execute(os.Stdout, submissionData); err != nil {
		logError(fmt.Errorf("failed to execute template: %w", err))
	}
}

func formatCodeOutput(code string, lang string) string {
	if len(code) > 1000 {
		code = code[:1000] + "...\n"
	}

	// More syntax files: https://github.com/zyedidia/highlight
	syntaxFile, err := highlightDir.ReadFile("highlight/" + lang + ".yaml")
	if err != nil {
		logError(fmt.Errorf("no syntax file for lang (%w)", err))
	}

	syntaxDef, err := highlight.ParseDef(syntaxFile)
	if err != nil {
		return code
	}
	h := highlight.NewHighlighter(syntaxDef)
	matches := h.HighlightString(code)

	var highlightedCode string

	lines := strings.Split(code, "\n")
	var printHl = color.New(color.Reset).SprintFunc()
	for lineN, l := range lines {
		for colN, c := range l {

			if group, ok := matches[lineN][colN]; ok {
				if group == highlight.Groups["statement"] {
					printHl = color.New(color.FgGreen).SprintFunc()
				} else if group == highlight.Groups["preproc"] {
					printHl = color.New(color.FgHiRed).SprintFunc()
				} else if group == highlight.Groups["identifier"] {
					printHl = color.New(color.FgRed).SprintFunc()
				} else if group == highlight.Groups["function"] {
					printHl = color.New(color.FgBlue).SprintFunc()
				} else if group == highlight.Groups["constant.string"] {
					printHl = color.New(color.FgHiCyan).SprintFunc()
				} else if group == highlight.Groups["constant.specialChar"] {
					printHl = color.New(color.FgHiMagenta).SprintFunc()
				} else if group == highlight.Groups["constant.number"] {
					printHl = color.New(color.FgHiBlue).SprintFunc()
				} else if group == highlight.Groups["constant.bool"] {
					printHl = color.New(color.FgHiBlue).SprintFunc()
				} else if group == highlight.Groups["symbol.brackets"] {
					printHl = color.New(color.FgRed).SprintFunc()
				} else if group == highlight.Groups["type"] {
					printHl = color.New(color.FgYellow).SprintFunc()
				} else if group == highlight.Groups["comment"] {
					printHl = color.New(color.FgHiBlack).SprintFunc()
				} else {
					printHl = color.New(color.Reset).SprintFunc()
				}
			}

			highlightedCode += printHl(string(c))
		}

		highlightedCode += "\n"
	}

	return highlightedCode
}

type Submit struct {
	Status string `json:"status"`
	Data   int    `json:"data"`
}

type SubmitError struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

type Languages struct {
	Data []struct {
		Name string `json:"internal_name"`
	} `json:"data"`
}

type LatestSubmission struct {
	Status string `json:"status"`
	Data   struct {
		Status       string `json:"status"`
		CompileError bool   `json:"compile_error"`
		Score        int    `json:"score"`
	}
}

func checkLanguages(problemId string, useCase int) []string {
	url := fmt.Sprintf(URL_LANGS_PB, problemId)
	body, err := MakeGetRequest(url, nil, RequestNone)
	if err != nil {
		logError(fmt.Errorf("failed to make request for problem ID %s: %w", problemId, err))
		return nil
	}

	var langs Languages
	if err := json.Unmarshal(body, &langs); err != nil {
		logError(fmt.Errorf("error unmarshalling languages response: %w", err))
		return nil
	}

	switch useCase {
	case 1:
		printLanguages(langs)
		return nil
	default:
		return extractLanguageNames(langs)
	}
}

func printLanguages(langs Languages) {
	for i, lang := range langs.Data {
		fmt.Printf("%d: %s\n", i+1, lang.Name)
	}
}

func extractLanguageNames(langs Languages) []string {
	var listLangs []string
	for _, lang := range langs.Data {
		listLangs = append(listLangs, lang.Name)
	}
	return listLangs
}

func uploadCode(id, language, file string) {
	codeFile, err := os.Open(file)
	if err != nil {
		logError(fmt.Errorf("failed to open code file %s: %w", file, err))
		return
	}
	defer codeFile.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	defer writer.Close()

	if err := writeFormFields(writer, id, language); err != nil {
		logError(err)
		return
	}

	if err := writeCodeFile(writer, file, codeFile); err != nil {
		logError(err)
		return
	}

	contentType := writer.FormDataContentType()

	body, err := MakePostRequest(URL_SUBMIT, &requestBody, RequestMultipartForm, contentType)
	if err != nil {
		logError(fmt.Errorf("error submitting code: %w", err))
		return
	}

	var data Submit
	if err := json.Unmarshal(body, &data); err != nil {
		handleSubmissionError(body)
		return
	}

	fmt.Printf("Submission sent: %s\nSubmission ID: %d\n", data.Status, data.Data)
	checkSubmissionStatus(data.Data)
}

func writeFormFields(writer *multipart.Writer, id, language string) error {
	if err := writer.WriteField("problem_id", id); err != nil {
		return fmt.Errorf("failed to write problem_id field: %w", err)
	}
	if err := writer.WriteField("language", language); err != nil {
		return fmt.Errorf("failed to write language field: %w", err)
	}
	return nil
}

func writeCodeFile(writer *multipart.Writer, file string, codeFile *os.File) error {
	fileWriter, err := writer.CreateFormFile("code", file)
	if err != nil {
		return fmt.Errorf("failed to create form file for code: %w", err)
	}
	if _, err := io.Copy(fileWriter, codeFile); err != nil {
		return fmt.Errorf("failed to copy code file content: %w", err)
	}
	return nil
}

func handleSubmissionError(body []byte) {
	var dataerr SubmitError
	if err := json.Unmarshal(body, &dataerr); err != nil {
		logError(fmt.Errorf("failed to parse error response: %w", err))
		return
	}
	logError(fmt.Errorf("status: %s\nmessage: %s", dataerr.Status, dataerr.Data))
}

func checkSubmissionStatus(submissionID int) {
	url := fmt.Sprintf(URL_LATEST_SUBMISSION, submissionID)

	action := func() {
		var dataLatestSubmit LatestSubmission
		for {
			body, err := MakeGetRequest(url, nil, RequestNone)
			if err != nil {
				logError(fmt.Errorf("failed to get submission status: %w", err))
				return
			}

			if err := json.Unmarshal(body, &dataLatestSubmit); err != nil {
				logError(fmt.Errorf("failed to parse latest submission response: %w", err))
				return
			}

			if dataLatestSubmit.Data.Status == "finished" {
				break
			}

			fmt.Print(".")
		}

		// Handle success or compilation failure
		if dataLatestSubmit.Data.CompileError {
			fmt.Println("\nCompilation failed! Score: 0")
		} else {
			fmt.Printf("\nSuccess! Score: %d\n", dataLatestSubmit.Data.Score)
		}
	}

	if err := spinner.New().Title("Waiting ...").Action(action).Run(); err != nil {
		logError(fmt.Errorf("spinner error: %w", err))
	}
}
