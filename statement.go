package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/glamour"
)

// print statement
type info struct {
	Data struct {
		Name          string  `json:"name"`
		Time          float64 `json:"time_limit"`
		MemoryLimit   int     `json:"memory_limit"`
		SourceSize    int     `json:"source_size"`
		SourceCredits string  `json:"source_credits"`
	} `json:"data"`
}

type statement struct {
	Status string `json:"status"`
	Data   struct {
		Data string `json:"data"`
	} `json:"data"`
}

func printStatement() {
	if len(os.Args) < 3 {
		fmt.Println("Enter problem id!")
		os.Exit(1)
	}
	id := os.Args[2]

	//info
	url := fmt.Sprintf("https://kilonova.ro/api/problem/%s/", id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}
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
		fmt.Println("Error: ", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading response body: %s", err)
		os.Exit(1)
	}

	var info info
	if err := json.Unmarshal(body, &info); err != nil {
		fmt.Printf("error unmarshalling response: %s", err)
		os.Exit(1)
	}
	fmt.Printf("\nName: %s\nID: #%s\nTime Limit: %.2fs\nMemory Limit: %dKB\nSource Size: %dKB\nCredits: %s\n", info.Data.Name, id, info.Data.Time, info.Data.MemoryLimit, info.Data.SourceSize, info.Data.SourceCredits)

	renderer, err := glamour.NewTermRenderer(glamour.WithStandardStyle("dark"))
	if err != nil {
		fmt.Printf("Error creating renderer: %v", err)
		os.Exit(1)
	}
	rendered, err := renderer.Render("# Statement")
	if err != nil {
		fmt.Printf("Error rendering markdown: %v", err)
		os.Exit(1)
	}
	fmt.Println(rendered)

	//statement
	url = fmt.Sprintf("https://kilonova.ro/api/problem/%s/get/attachmentByName/statement-ro.md", id)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	req.Header.Set("User-Agent", "KilonovaCLIClient/1.0")

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading response body: %s", err)
		os.Exit(1)
	}

	if strings.Contains(string(body), "\"status\":\"error\"") {
		url = fmt.Sprintf("https://kilonova.ro/api/problem/%s/get/attachmentByName/statement-en.md", id)
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		req.Header.Set("User-Agent", "KilonovaCLIClient/1.0")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("error reading response body: %s", err)
			os.Exit(1)
		}
	}

	var data statement
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("error unmarshalling response: %s", err)
		os.Exit(1)
	}
	text, err := b64.StdEncoding.DecodeString(data.Data.Data)
	if err != nil {
		fmt.Printf("error decoding base64 data: %s", err)
		os.Exit(1)
	}
	decodedtext := string(text)

	replacements := map[string]string{
		"$":              "",
		" \\ ":           "",
		"\\ldots":        "...",
		"\\leq":          "≤",
		"\\geq":          "≥",
		"\\el":           "",
		"\\in":           "∈",
		"\\le":           "≤",
		"\\qe":           "≥",
		"\\pm":           "±",
		"\\cdot":         "•",
		"\\sum_":         "Σ ",
		"\\displaystyle": "",
		"\\times":        "x",
		"\\%":            "%",
	}

	for old, new := range replacements {
		decodedtext = strings.ReplaceAll(decodedtext, old, new)
	}

	re := regexp.MustCompile(`\\text{(.*?)}`)
	decodedtext = re.ReplaceAllString(decodedtext, "$1")

	re = regexp.MustCompile(`\\texttt{(.*?)}`)
	decodedtext = re.ReplaceAllString(decodedtext, "$1")

	re = regexp.MustCompile(`\\bm{(.*?)}`)
	decodedtext = re.ReplaceAllString(decodedtext, "$1")

	re = regexp.MustCompile(`\\textit{(.*?)}`)
	decodedtext = re.ReplaceAllString(decodedtext, "$1")

	re = regexp.MustCompile(`\\rule\{[^}]+\}\{[^}]+\}`)
	decodedtext = re.ReplaceAllString(decodedtext, "$1")

	re = regexp.MustCompile(`~\[([^\]]+)\]`)
	//decodedtext = re.ReplaceAllString(decodedtext, "https://kilonova.ro/assets/problem/"+id+"/attachment/$1")
	decodedtext = re.ReplaceAllString(decodedtext, "Error: Unable to show photo! View in browser.")

	re = regexp.MustCompile(`^([^\|]+)`)
	decodedtext = re.ReplaceAllString(decodedtext, "$1")

	if err != nil {
		fmt.Printf("Error creating renderer: %v", err)
		os.Exit(1)
	}
	rendered, err = renderer.Render(decodedtext)
	if err != nil {
		fmt.Printf("Error rendering markdown: %v", err)
		os.Exit(1)
	}
	fmt.Println(rendered)
}
