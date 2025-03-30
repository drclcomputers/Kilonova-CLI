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
	if len(os.Args) == 3 {
		problem_name = os.Args[2]
	}
	url := URL_SEARCH

	searchData := map[string]interface{}{
		"name_fuzzy": problem_name,
		"offset":     0,
	}

	jsonData, err := json.Marshal(searchData)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	body, err := makeRequest("POST", url, bytes.NewBuffer(jsonData), "2")
	if err != nil {
		logErr(err)
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
		fmt.Println("Id  |  Name  |  Source")
		for p := 0; p < pages; p++ {
			searchData["offset"] = p * 50

			jsonData, err := json.Marshal(searchData)
			if err != nil {
				logErr(err)
			}

			body, err = makeRequest("POST", url, bytes.NewBuffer(jsonData), "2")
			if err != nil {
				logErr(err)
			}

			err = json.Unmarshal(body, &data)
			if err != nil {
				logErr(err)
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
