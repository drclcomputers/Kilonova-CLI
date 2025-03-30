package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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
	url := "https://kilonova.ro/api/user/self/"
	token, err := os.ReadFile("token")
	if err != nil || string(token) == "" {
		fmt.Println("Could not read session ID from file. Make sure you are logged in!")
		os.Exit(1)
	}

	jsonData := []byte(`{"key": "value"}`)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	req.Header.Set("User-Agent", "KilonovaCLIClient/1.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", string(token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s", err)
		os.Exit(1)
	}

	var data userid
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("error unmarshalling response: %s", err)
		os.Exit(1)
	}

	id := data.Data.ID

	//get submissions on problem
	if len(os.Args) < 3 {
		fmt.Println("Usage: <program> -submissions <problem_id>")
		os.Exit(1)
	}

	url = fmt.Sprintf("https://kilonova.ro/api/submissions/get?ascending=false&limit=500&offset=0&ordering=id&problem_id=%s&user_id=%d", os.Args[2], id)

	req, err = http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	req.Header.Set("User-Agent", "KilonovaCLIClient/1.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", string(token))

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s", err)
		os.Exit(1)
	}

	var datasub submissionlist
	if err := json.Unmarshal(body, &datasub); err != nil {
		fmt.Printf("error unmarshalling response: %s", err)
		os.Exit(1)
	}

	fmt.Println("Nr.  |  ID  |  Time  |  Language |  Score")

	for i := range datasub.Data.Count {
		parsedTime, err := time.Parse(time.RFC3339Nano, datasub.Data.Submissions[i].Created_at)
		formattedTime := parsedTime.Format("2006-01-02 15:04:05")
		if err != nil {
			fmt.Printf("Could not parse time %s", err)
			os.Exit(1)
		}
		fmt.Printf("%d.  |  %d  |  %s  |  %s  |  %d\n", i+1, datasub.Data.Submissions[i].Id, formattedTime, datasub.Data.Submissions[i].Language, datasub.Data.Submissions[i].Score)
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
	url := fmt.Sprintf("https://kilonova.ro/api/problem/%s/languages", os.Args[2])
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("error fetching data: %s", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading response body: %s", err)
		os.Exit(1)
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
	url := "https://kilonova.ro/api/submissions/submit"

	id := os.Args[2]
	lang := os.Args[3]
	file := os.Args[4]

	codeFile, err := os.Open(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	defer codeFile.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	_ = writer.WriteField("problem_id", id)

	_ = writer.WriteField("language", lang)

	fileWriter, err := writer.CreateFormFile("code", file)
	if err != nil {
		fmt.Println("Error creating form file:", err)
		os.Exit(1)
	}
	_, err = io.Copy(fileWriter, codeFile)
	if err != nil {
		fmt.Println("Error copying file:", err)
		os.Exit(1)
	}

	writer.Close()

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	req.Header.Set("User-Agent", "KilonovaCLIClient/1.0")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	token, err := os.ReadFile("token")
	if err != nil {
		fmt.Println("Not logged in! Please sign in!")
		os.Exit(1)
	}
	req.Header.Set("Authorization", string(token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	var data submit
	err = json.Unmarshal(body, &data)
	if err != nil {
		var dataerr submiterr
		err = json.Unmarshal(body, &dataerr)
		if err != nil {
			fmt.Println("Error unmarshalling json file. Err: ", err)
			os.Exit(1)
		}
		fmt.Printf("Status: %s\nMessage: %s", dataerr.Status, dataerr.Data)
	}
	fmt.Printf("Status: %s\nSubmission ID: %d", data.Status, data.Data)
}
