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
		checklangs(args[0])
	},
}

var printSubmissionsCmd = &cobra.Command{
	Use:   "submissions [ID]",
	Short: "View sent submissions to a problem.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		printSubmissions(args[0])
	},
}

var printSubmInfo = &cobra.Command{
	Use:   "submissioninfo [Problem ID] [Submission ID]",
	Short: "View a detailed description of a sent submission.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		printDetailsSubmissions(args[0], args[1])
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
			Id              int     `json:"id"`
			Created_at      string  `json:"created_at"`
			Language        string  `json:"language"`
			Score           int     `json:"score"`
			Max_memory      int     `json:"max_memory"`
			Max_time        float64 `json:"max_time"`
			Compile_error   bool    `json:"compile_error"`
			Compile_message string  `json:"compile_message"`
		}
		Count int `json:"count"`
	} `json:"data"`
}

func printSubmissions(problem_id string) {
	user_id := getUserID()

	//get submissions on problem

	jsonData := []byte(`{"key": "value"}`)
	url := fmt.Sprintf(URL_SUBMISSION_LIST, problem_id, user_id)

	body, err := makeRequest("GET", url, bytes.NewBuffer(jsonData), "1")
	if err != nil {
		logErr(err)
		return
	}

	var datasub submissionlist
	if err := json.Unmarshal(body, &datasub); err != nil {
		logErr(err)
		return
	}

	var rows []table.Row

	for i := range datasub.Data.Count {
		parsedTime, err := time.Parse(time.RFC3339Nano,
			datasub.Data.Submissions[i].Created_at)
		formattedTime := parsedTime.Format("2006-01-02 15:04:05")

		if err != nil {
			logErr(err)
			return
		}

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", datasub.Data.Submissions[i].Id),
			formattedTime,
			datasub.Data.Submissions[i].Language,
			fmt.Sprintf("%d", datasub.Data.Submissions[i].Score),
		})
	}

	columns := []table.Column{
		{Title: "ID", Width: 10},
		{Title: "TIme", Width: 20},
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
	if err := p.Start(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}

}

func printDetailsSubmissions(problem_id, submission_id string) {
	user_id := getUserID()

	jsonData := []byte(`{"key": "value"}`)
	url := fmt.Sprintf(URL_SUBMISSION_LIST, problem_id, user_id)

	body, err := makeRequest("GET", url, bytes.NewBuffer(jsonData), "1")
	if err != nil {
		logErr(err)
		return
	}

	var datasub submissionlist
	if err := json.Unmarshal(body, &datasub); err != nil {
		logErr(err)
		return
	}

	for i := range datasub.Data.Count {
		if nr, err := strconv.Atoi(submission_id); err == nil && datasub.Data.Submissions[i].Id == nr {
			parsedTime, err := time.Parse(time.RFC3339Nano,
				datasub.Data.Submissions[i].Created_at)
			formattedTime := parsedTime.Format("2006-01-02 15:04:05")

			if err != nil {
				logErr(err)
				return
			}
			problemInfo(problem_id)

			fmt.Printf("\nSubmission ID: #%d\nCreated: %s\nLanguage: %s\nScore: %d\n",
				datasub.Data.Submissions[i].Id, formattedTime,
				datasub.Data.Submissions[i].Language, datasub.Data.Submissions[i].Score)
			fmt.Printf("Max Memory: %dKB\nMax time: %.2fs\nCompile error: %t\nCompile message: %s\n\n",
				datasub.Data.Submissions[i].Max_memory, datasub.Data.Submissions[i].Max_time,
				datasub.Data.Submissions[i].Compile_error, datasub.Data.Submissions[i].Compile_message)
			return
		}
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

func checklangs(problem_id string) {
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

	for i := range langs.Data {
		fmt.Println(i+1, langs.Data[i].Name)
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
	fmt.Printf("Status: %s\nSubmission ID: %d\n", data.Status, data.Data)
}
