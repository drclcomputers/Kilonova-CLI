// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package cmd

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	u "net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

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
			logErr(fmt.Errorf("invalid first page number: %v", err))
			return
		}

		lastPage, err := strconv.Atoi(args[3])
		if err != nil {
			logErr(fmt.Errorf("invalid last page number: %v", err))
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

func parseSubmissionTime(timeStr string) string {
	parsedTime, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		logErr(fmt.Errorf("error parsing time: %v", err))
		return ""
	}
	return parsedTime.Format("2006-01-02 15:04:05")
}

func printSubmissions(problemId, userId string, firstPage, lastPage int) {
	if userId == "me" {
		userId = getUserID()
	}

	if firstPage > lastPage {
		logErr(fmt.Errorf("first page cannot be bigger than the last page"))
	}

	if firstPage < 0 || lastPage < 0 {
		logErr(fmt.Errorf("pages need to be positive numbers, different from 0"))
	}

	var datasub SubmissionList

	var rows []table.Row

	count := -1
	var startOffset = (firstPage - 1) * 50
	var endOffset = max((lastPage-1)*50, 50)

	for offset := max(startOffset, 0); offset < count && offset < endOffset; offset += 50 {
		var url = getSubmissionURL(userId, problemId, offset)

		body, err := makeRequest("GET", url, nil, "1")
		if err != nil {
			logErr(err)
		}

		err = json.Unmarshal(body, &datasub)
		if err != nil {
			logErr(err)
		}

		count = datasub.Data.Count

		for _, problem := range datasub.Data.Submissions {
			formattedTime := parseSubmissionTime(problem.CreatedAt)
			if formattedTime == "" {
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

	p := tea.NewProgram(model{table: t}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logErr(fmt.Errorf("error running program: %v", err))
	}
}

func downloadSource(submission_id, code string) {
	homedir, err := os.Getwd()
	if err != nil {
		logErr(err)
	}
	configDir := filepath.Join(homedir)
	downFile := filepath.Join(configDir, "source_"+submission_id+".txt")
	file, err := os.Create(downFile)
	if err != nil {
		logErr(fmt.Errorf("error creating file: %v", err))
	}
	defer file.Close()

	if err := os.WriteFile(downFile, []byte(code), 0644); err != nil {
		logErr(fmt.Errorf("error writing source code to file: %v", err))
	}
}

func printDetailsSubmission(submissionId string) {

	var datasub SubmissionDetails

	subsid, err := strconv.Atoi(submissionId)
	if err != nil {
		logErr(err)
	}

	url := fmt.Sprintf(URL_LATEST_SUBMISSION, subsid)

	formData := u.Values{
		"id": {submissionId},
	}

	body, err := makeRequest("GET", url, bytes.NewBufferString(formData.Encode()), "0")
	if err != nil {
		logErr(fmt.Errorf("error: %v", err))
	}

	err = json.Unmarshal(body, &datasub)
	if err != nil {
		logErr(err)
	}

	if datasub.Status != "success" {
		logErr(fmt.Errorf("error: %v", datasub.Status))
	}

	parsedTime, err := time.Parse(time.RFC3339Nano, datasub.Data.CreatedAt)
	if err != nil {
		logErr(err)
	}

	formattedTime := parsedTime.Format("2006-01-02 15:04:05")

	pbid := strconv.Itoa(datasub.Data.ProblemID)

	fmt.Println(problemInfo(pbid))

	code, _ := b64.StdEncoding.DecodeString(datasub.Data.Code)

	fmt.Printf("\nSubmission ID: #%d\nCreated: %s\nLanguage: %s\nScore: %.0f\n",
		datasub.Data.Id, formattedTime,
		datasub.Data.Language, datasub.Data.Score)
	fmt.Printf("Max Memory: %dKB\nMax time: %.2fs\nCompile error: %t\nCompile message: %s\n\nCode:\n%s\n",
		datasub.Data.MaxMemory, datasub.Data.MaxTime,
		datasub.Data.CompileError, datasub.Data.CompileMessage, code)

	action := func() { downloadSource(submissionId, string(code)) }
	if download {
		if err := spinner.New().Title("Waiting ...").Action(action).Run(); err != nil {
			logErr(err)
		}
	}
}

// upload code
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
	//get languages
	url := fmt.Sprintf(URL_LANGS_PB, problemId)
	body, err := makeRequest("GET", url, nil, "0")
	if err != nil {
		logErr(err)
	}

	var langs Languages
	if err := json.Unmarshal(body, &langs); err != nil {
		logErr(fmt.Errorf("error unmarshalling response: %s", err))
	}

	if useCase == 1 {
		for i := range langs.Data {
			fmt.Println(i+1, langs.Data[i].Name)
		}
		return nil
	} else {
		var listLangs []string
		for i := range langs.Data {
			listLangs = append(listLangs, langs.Data[i].Name)
		}
		return listLangs
	}
}

func uploadCode(id, language, file string) {
	//upload code

	codeFile, err := os.Open(file)
	if err != nil {
		logErr(err)
	}
	defer codeFile.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	_ = writer.WriteField("problem_id", id)

	_ = writer.WriteField("language", language)

	fileWriter, err := writer.CreateFormFile("code", file)
	if err != nil {
		logErr(err)
	}
	_, err = io.Copy(fileWriter, codeFile)
	if err != nil {
		logErr(err)
	}

	writer.Close()

	body, err := makeRequest("POST", URL_SUBMIT, io.Reader(&requestBody), string(writer.FormDataContentType()))
	if err != nil {
		logErr(err)
	}

	var data Submit
	err = json.Unmarshal(body, &data)
	if err != nil {
		var dataerr SubmitError
		err = json.Unmarshal(body, &dataerr)
		if err != nil {
			logErr(err)
		}
		logErr(fmt.Errorf("status: %s\nmessage: %s", dataerr.Status, dataerr.Data))
	}
	fmt.Printf("Submission sent: %s\nSubmission ID: %d\n", data.Status, data.Data)

	url := fmt.Sprintf(URL_LATEST_SUBMISSION, data.Data)

	body, err = makeRequest("GET", url, nil, "0")
	if err != nil {
		logErr(err)
	}

	var dataLatestSubmit LatestSubmission
	if err = json.Unmarshal(body, &dataLatestSubmit); err != nil {
		logErr(err)
	}

	action := func() {
		for dataLatestSubmit.Data.Status != "finished" {
			fmt.Print(".")
			body, err = makeRequest("GET", url, nil, "0")
			if err != nil {
				logErr(err)
			}

			if err = json.Unmarshal(body, &dataLatestSubmit); err != nil {
				logErr(err)
			}
		}
	}
	if err := spinner.New().Title("Waiting ...").Action(action).Run(); err != nil {
		logErr(err)
	}
	if !dataLatestSubmit.Data.CompileError {
		fmt.Println("\nSuccess! Score: ", dataLatestSubmit.Data.Score)
	} else {
		fmt.Println("\nCompilation failed! Score: 0")
	}
}
