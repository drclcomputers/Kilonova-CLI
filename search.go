package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// search problems
type search struct {
	Data struct {
		Count    int `json:"count"`
		Problems []struct {
			Id             int    `json:"id"`
			Name           string `json:"name"`
			Source_Credits string `json:"source_credits"`
		}
	} `json:"data"`
}

func searchProblems() {
	problem_name := ""
	if len(os.Args) == 3 {
		problem_name = os.Args[2]
	}
	url := "https://kilonova.ro/api/problem/search"

	searchData := map[string]interface{}{
		"name_fuzzy": problem_name,
		"offset":     0,
	}

	jsonData, err := json.Marshal(searchData)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "KilonovaCLIClient/1.0")
	tokenbyte, err := os.ReadFile("token")
	token := string(tokenbyte)
	if err != nil {
		token = "guest"
	}

	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	var data search
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		os.Exit(1)
	}

	pages := data.Data.Count / 50
	if data.Data.Count%50 != 0 {
		pages++
	}

	cnt := data.Data.Count
	if cnt == 0 {
		fmt.Println("No problems found")
	} else {
		fmt.Println(cnt, "problems found!")
		fmt.Println("Id | Name | Source")
		for p := 0; p < pages; p++ {
			searchData["offset"] = p * 50

			jsonData, err := json.Marshal(searchData)
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
				return
			}

			req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", "KilonovaCLIClient/1.0")
			req.Header.Set("Authorization", token)

			client = &http.Client{}
			resp, err = client.Do(req)
			if err != nil {
				fmt.Println("Error sending POST request:", err)
				os.Exit(1)
			}
			defer resp.Body.Close()

			body, err = io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				os.Exit(1)
			}

			err = json.Unmarshal(body, &data)
			if err != nil {
				fmt.Println("Error unmarshaling JSON:", err)
				os.Exit(1)
			}

			cntpag := len(data.Data.Problems)
			if p == pages-1 {
				cntpag = cnt % 50
			}
			for i := 0; i < cntpag; i++ {
				fmt.Println(data.Data.Problems[i].Id, "  |  ", data.Data.Problems[i].Name, "  |  ", data.Data.Problems[i].Source_Credits)
			}
		}
	}
}
