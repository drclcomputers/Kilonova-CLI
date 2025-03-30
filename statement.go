package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/glamour"
)

func formatText(decodedtext string) string {
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
	decodedtext = re.ReplaceAllString(decodedtext, "Error: Unable to show photo! View in browser.")

	re = regexp.MustCompile(`^([^\|]+)`)
	decodedtext = re.ReplaceAllString(decodedtext, "$1")

	return decodedtext
}

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
	if len(os.Args) < 4 {
		fmt.Println("Usage: <program> -statement <ID> <RO or EN>")
		os.Exit(1)
	}
	id := os.Args[2]
	language := os.Args[3]

	//info
	url := fmt.Sprintf(URL_PROBLEM, id)
	body, err := makeRequest("GET", url, nil, "0")
	if err != nil {
		logErr(err)
	}

	var info info
	if err := json.Unmarshal(body, &info); err != nil {
		logErr(err)
	}
	fmt.Printf("\nName: %s\nID: #%s\nTime Limit: %.2fs\nMemory Limit: %dKB\nSource Size: %dKB\nCredits: %s\n", info.Data.Name, id, info.Data.Time, info.Data.MemoryLimit, info.Data.SourceSize, info.Data.SourceCredits)

	renderer, err := glamour.NewTermRenderer(glamour.WithStandardStyle("dark"))
	if err != nil {
		logErr(err)
	}
	rendered, err := renderer.Render("# Statement")
	if err != nil {
		logErr(err)
	}
	fmt.Println(rendered)

	//statement
	if language == "RO" {
		url = fmt.Sprintf(URL_STATEMENT, id, STAT_FILENAME_RO)
	} else {
		url = fmt.Sprintf(URL_STATEMENT, id, STAT_FILENAME_EN)
	}
	body, err = makeRequest("GET", url, nil, "0")
	if err != nil {
		logErr(err)
	}

	if strings.Contains(string(body), "\"status\":\"error\"") {
		log.Fatal("Error: Problem statement is not available in the chosen language!")
	}

	var data statement
	if err := json.Unmarshal(body, &data); err != nil {
		logErr(err)
	}
	text, err := b64.StdEncoding.DecodeString(data.Data.Data)
	if err != nil {
		logErr(err)
	}
	decodedtext := string(text)

	decodedtext = formatText(decodedtext)

	if err != nil {
		logErr(err)
	}
	rendered, err = renderer.Render(decodedtext)
	if err != nil {
		logErr(err)
	}
	fmt.Println(rendered)
}
