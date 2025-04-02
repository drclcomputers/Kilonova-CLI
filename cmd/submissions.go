// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

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
		checklangs(args[0], 1)
	},
}

var printSubmissionsCmd = &cobra.Command{
	Use:   "submissions [Problem ID or all (all problems)] [User ID, me (personal submissions), all (all users)] [1st page] [last page]",
	Short: "View sent submissions to a problem.",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		printSubmissions(args[0], args[1], args[2], args[3])
	},
}

var printSubmInfo = &cobra.Command{
	Use:   "submissioninfo [Problem ID] [User ID] [Submission ID]",
	Short: "View a detailed description of a sent submission.",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		printDetailsSubmissions(args[0], args[1], args[2])
	},
}

func init() {
	rootCmd.AddCommand(checkLangsCmd)
	rootCmd.AddCommand(uploadCodeCmd)
	rootCmd.AddCommand(printSubmissionsCmd)
	rootCmd.AddCommand(printSubmInfo)
}

// print submissions

type submissionlist struct {
	Data struct {
		Submissions []struct {
			UserID          int     `json:"user_id"`
			ProblemID       int     `json:"problem_id"`
			Id              int     `json:"id"`
			Created_at      string  `json:"created_at"`
			Language        string  `json:"language"`
			Score           float64 `json:"score"`
			Max_memory      int     `json:"max_memory"`
			Max_time        float64 `json:"max_time"`
			Compile_error   bool    `json:"compile_error"`
			Compile_message string  `json:"compile_message"`
		}
		Count int `json:"count"`
	} `json:"data"`
}

func printSubmissions(problem_id, user_id, fpag, lpag string) {
	if user_id == "me" {
		user_id = getUserID()
	}

	//get submissions on problem

	var datasub submissionlist

	var rows []table.Row

	cnt := -1
	fpagnr, err := strconv.Atoi(fpag)
	if err != nil {
		logErr(err)
		return
	}

	lpagnr, err := strconv.Atoi(lpag)
	if err != nil {
		logErr(err)
		return
	}

	for offset := max((fpagnr-1)*50, 0); cnt == -1 || (offset < cnt && offset < max((lpagnr-1)*50, 50)); offset += 50 {
		var url string
		switch {
		case user_id == "all" && problem_id == "all":
			url = fmt.Sprintf(URL_SUBMISSION_LIST_NO_FILTER, offset)
		case user_id == "all":
			url = fmt.Sprintf(URL_SUBMISSION_LIST_NO_USER, offset, problem_id)
		case problem_id == "all":
			url = fmt.Sprintf(URL_SUBMISSION_LIST_NO_PROBLEM, offset, user_id)
		default:
			url = fmt.Sprintf(URL_SUBMISSION_LIST, offset, problem_id, user_id)
		}

		body, err := makeRequest("GET", url, nil, "1")
		if err != nil {
			logErr(err)
			return
		}

		err = json.Unmarshal(body, &datasub)
		if err != nil {
			logErr(err)
			return
		}

		cnt = datasub.Data.Count

		for _, problem := range datasub.Data.Submissions {
			parsedTime, err := time.Parse(time.RFC3339Nano, problem.Created_at)
			formattedTime := parsedTime.Format("2006-01-02 15:04:05")

			if err != nil {
				logErr(err)
				return
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
		table.WithHeight(10),
	)

	t.SetStyles(table.DefaultStyles())

	p := tea.NewProgram(model{table: t})
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}

}

func printDetailsSubmissions(problem_id, user_id, submission_id string) {

	if user_id == "me" {
		user_id = getUserID()
	}

	//get submissions on problem

	var datasub submissionlist

	cnt := 50

	ok := false

	for offset := 0; offset < cnt; offset += 50 {
		url := fmt.Sprintf(URL_SUBMISSION_LIST, offset, problem_id, user_id)

		body, err := makeRequest("GET", url, nil, "1")
		if err != nil {
			logErr(err)
			return
		}

		err = json.Unmarshal(body, &datasub)
		if err != nil {
			logErr(err)
			return
		}

		cnt = datasub.Data.Count

		for _, problem := range datasub.Data.Submissions {
			if nr, err := strconv.Atoi(submission_id); err == nil && problem.Id == nr {
				ok = true

				parsedTime, err := time.Parse(time.RFC3339Nano, problem.Created_at)
				formattedTime := parsedTime.Format("2006-01-02 15:04:05")

				if err != nil {
					logErr(err)
					return
				}
				problemInfo(problem_id)

				fmt.Printf("\nSubmission ID: #%d\nCreated: %s\nLanguage: %s\nScore: %.0f\n",
					problem.Id, formattedTime,
					problem.Language, problem.Score)
				fmt.Printf("Max Memory: %dKB\nMax time: %.2fs\nCompile error: %t\nCompile message: %s\n\n",
					problem.Max_memory, problem.Max_time,
					problem.Compile_error, problem.Compile_message)
				return
			}
		}
	}

	if !ok {
		fmt.Println("No submission for these parameters were found!")
	}
}

// upload code
type submit struct {
	Status string `json:"status"`
	Data   int    `json:"data"`
}

type submiterr struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

type langs struct {
	Data []struct {
		Name string `json:"internal_name"`
	} `json:"data"`
}

type latestSubmission struct {
	Status string `json:"status"`
	Data   struct {
		Status       string `json:"status"`
		CompileError bool   `json:"compile_error"`
		Score        int    `json:"score"`
	}
}

func checklangs(problem_id string, use_case int) []string {
	//get languages
	url := fmt.Sprintf(URL_LANGS_PB, problem_id)
	body, err := makeRequest("GET", url, nil, "0")
	if err != nil {
		logErr(err)
	}

	var langs langs
	if err := json.Unmarshal(body, &langs); err != nil {
		fmt.Printf("error unmarshalling response: %s", err)
		os.Exit(1)
	}

	if use_case == 1 {
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

func uploadCode(id, lang, file string) {
	//upload code

	codeFile, err := os.Open(file)
	if err != nil {
		logErr(err)
	}
	defer codeFile.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	_ = writer.WriteField("problem_id", id)

	_ = writer.WriteField("language", lang)

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

	var data submit
	err = json.Unmarshal(body, &data)
	if err != nil {
		var dataerr submiterr
		err = json.Unmarshal(body, &dataerr)
		if err != nil {
			logErr(err)
		}
		fmt.Printf("Status: %s\nMessage: %s\n", dataerr.Status, dataerr.Data)
	}
	fmt.Printf("Submission sent: %s\nSubmission ID: %d\n", data.Status, data.Data)

	url := fmt.Sprintf(URL_LATEST_SUBMISSION, data.Data)

	body, err = makeRequest("GET", url, nil, "0")
	if err != nil {
		logErr(err)
	}

	var dataLatestSubmit latestSubmission
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
		log.Fatal(err)
	}

	if !dataLatestSubmit.Data.CompileError {
		fmt.Println("\nSucces! Score: ", dataLatestSubmit.Data.Score)
	} else {
		fmt.Println("\nCompilation failed! Score: 0")
	}

}
