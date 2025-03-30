package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	if len(os.Args) >= 3 {
		problem_name = os.Args[2]
	}

	searchData := map[string]interface{}{
		"name_fuzzy": problem_name,
		"offset":     0,
	}

	jsonData, err := json.Marshal(searchData)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	body, err := makeRequest("POST", URL_SEARCH, bytes.NewBuffer(jsonData), "2")
	if err != nil {
		logErr(err)
	}

	var data search
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		os.Exit(1)
	}

	if data.Data.Count == 0 {
		fmt.Println("No problems found.")
		return
	}

	pages := data.Data.Count / 50
	if data.Data.Count%50 != 0 {
		pages++
	}

	cnt := data.Data.Count
	fmt.Println(cnt, "problems found!")
	fmt.Println("Id  |  Name  |  Source")
	for offset := 0; offset < data.Data.Count; offset += 50 {
		searchData["offset"] = offset

		jsonData, err := json.Marshal(searchData)
		if err != nil {
			logErr(err)
		}

		body, err = makeRequest("POST", URL_SEARCH, bytes.NewBuffer(jsonData), "2")
		if err != nil {
			logErr(err)
		}

		err = json.Unmarshal(body, &data)
		if err != nil {
			logErr(err)
		}

		for _, problem := range data.Data.Problems {
			fmt.Printf("%-4d | %-30s | %s\n", problem.Id, problem.Name, problem.Source_Credits)
		}
	}
}
