package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"
)

// print submissions
type userid struct {
	Data struct {
		ID int `json:"id"`
	} `json:"data"`
}

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

func printSubmissions() {
	//get user id

	jsonData := []byte(`{"key": "value"}`)
	body, err := makeRequest("GET", URL_SELF, bytes.NewBuffer(jsonData), "1")
	if err != nil {
		logErr(err)
		return
	}

	var data userid
	if err := json.Unmarshal(body, &data); err != nil {
		logErr(err)
		return
	}

	id := data.Data.ID

	//get submissions on problem
	if len(os.Args) < 3 {
		fmt.Println("Usage: <program> -submissions <problem_id>")
		os.Exit(1)
	}

	url := fmt.Sprintf(URL_SUBMISSION_LIST, os.Args[2], id)

	body, err = makeRequest("GET", url, bytes.NewBuffer(jsonData), "1")
	if err != nil {
		logErr(err)
		return
	}

	var datasub submissionlist
	if err := json.Unmarshal(body, &datasub); err != nil {
		logErr(err)
		return
	}

	fmt.Println("Nr.  |  ID  |  Time  |  Language |  Score")

	for i := range datasub.Data.Count {
		parsedTime, err := time.Parse(time.RFC3339Nano,
			datasub.Data.Submissions[i].Created_at)
		formattedTime := parsedTime.Format("2006-01-02 15:04:05")

		if err != nil {
			logErr(err)
			return
		}

		fmt.Printf("%d.  |  %d  |  %s  |  %s  |  %d\n",
			i+1, datasub.Data.Submissions[i].Id, formattedTime,
			datasub.Data.Submissions[i].Language,
			datasub.Data.Submissions[i].Score)
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

func checklangs() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: <program> -langs <problem_id>")
		os.Exit(1)
	}
	//get languages
	url := fmt.Sprintf(URL_LANGS_PB, os.Args[2])
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

func uploadCode() {

	if len(os.Args) < 5 {
		fmt.Println("Usage: <program> -upload <problem_id> <language> <filename>")
		os.Exit(1)
	}

	//upload code
	url := URL_SUBMIT

	id := os.Args[2]
	lang := os.Args[3]
	file := os.Args[4]

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

	body, err := makeRequest("POST", url, io.Reader(&requestBody), string(writer.FormDataContentType()))
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
		fmt.Printf("Status: %s\nMessage: %s", dataerr.Status, dataerr.Data)
	}
	fmt.Printf("Status: %s\nSubmission ID: %d", data.Status, data.Data)
}
