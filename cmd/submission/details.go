// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package submission

import (
	"bytes"
	"embed"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	problem "kncli/cmd/problems"
	"kncli/internal"
	u "net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/charmbracelet/huh/spinner"
	"github.com/fatih/color"
	"github.com/zyedidia/highlight"
)

//go:embed highlight/*.yaml
var highlightDir embed.FS

func downloadSource(submissionId, code string) {
	homedir, err := os.Getwd()
	if err != nil {
		internal.LogError(fmt.Errorf("failed to get current working directory: %w", err))
		return
	}

	downFile := filepath.Join(homedir, "source_"+submissionId+".txt")
	if err := os.WriteFile(downFile, []byte(code), 0644); err != nil {
		internal.LogError(fmt.Errorf("failed to write source code to file %q: %w", downFile, err))
		return
	}

	fmt.Printf("Source code for submission #%s saved to %q\n", submissionId, downFile)
}

func printDetailsSubmission(submissionId string) {
	var details SubmissionDetails

	url := fmt.Sprintf(internal.URL_LATEST_SUBMISSION, submissionId)

	formData := u.Values{
		"id": {submissionId},
	}

	ResponseBody, err := internal.MakeGetRequest(url, bytes.NewBufferString(formData.Encode()), internal.RequestNone)
	if err != nil {
		internal.LogError(fmt.Errorf("error fetching submission details: %w", err))
		return
	}

	if err := json.Unmarshal(ResponseBody, &details); err != nil {
		internal.LogError(fmt.Errorf("error unmarshalling response: %w", err))
		return
	}

	if details.Status != internal.SUCCESS {
		internal.LogError(fmt.Errorf("submission fetch failed with status: %v", details.Status))
		return
	}

	formattedTime, err := internal.ParseTime(details.Data.CreatedAt)
	if err != nil {
		internal.LogError(err)
		return
	}

	ProblemID := strconv.Itoa(details.Data.ProblemID)
	fmt.Println(problem.GetProblemInfoText(ProblemID))

	code, err := b64.StdEncoding.DecodeString(details.Data.Code)
	if err != nil {
		internal.LogError(fmt.Errorf("error decoding source code: %w", err))
		return
	}

	printTemplateSubmission(details, formattedTime, code)

	if shouldDownload {
		action := func() { downloadSource(submissionId, string(code)) }
		if err := spinner.New().Title("Downloading...").Action(action).Run(); err != nil {
			internal.LogError(fmt.Errorf("error during downloading source code for submission #%s: %w", submissionId, err))
			return
		}
	}
}

func printTemplateSubmission(details SubmissionDetails, formattedTime string, code []byte) {
	submissionData := SubmissionDetailsTemplate{
		ID:             details.Data.Id,
		CreatedAt:      formattedTime,
		Language:       details.Data.Language,
		Score:          details.Data.Score,
		MaxMemory:      details.Data.MaxMemory,
		MaxTime:        details.Data.MaxTime,
		CompileError:   details.Data.CompileError,
		CompileMessage: details.Data.CompileMessage,
		ContestID:      details.Data.ContestID,
		Code:           formatCodeOutput(string(code), details.Data.Language),
	}

	if submissionData.Code == internal.ERROR {
		submissionData.Code = string(code)
	}

	tmpl, err := template.New("submissionDetails").Parse(internal.SubmissionTemplate)
	if err != nil {
		internal.LogError(fmt.Errorf("failed to parse template: %w", err))
		return
	}

	if err := tmpl.Execute(os.Stdout, submissionData); err != nil {
		internal.LogError(fmt.Errorf("failed to execute template: %w", err))
		return
	}
}

func formatCodeOutput(code string, lang string) string {
	if len(code) > 600 {
		code = code[:600] + "...\n"
	}

	// More syntax files: https://github.com/zyedidia/highlight
	syntaxFile, err := highlightDir.ReadFile("highlight/" + lang + ".yaml")
	if err != nil {
		return internal.ERROR
	}

	syntaxDef, err := highlight.ParseDef(syntaxFile)
	if err != nil {
		return code
	}
	h := highlight.NewHighlighter(syntaxDef)
	matches := h.HighlightString(code)

	var highlightedCode string

	lines := strings.Split(code, "\n")
	var printHl = color.New(color.Reset).SprintFunc()
	for lineN, l := range lines {
		for colN, c := range l {

			if group, ok := matches[lineN][colN]; ok {
				if group == highlight.Groups["statement"] {
					printHl = color.New(color.FgGreen).SprintFunc()
				} else if group == highlight.Groups["preproc"] {
					printHl = color.New(color.FgHiRed).SprintFunc()
				} else if group == highlight.Groups["identifier"] {
					printHl = color.New(color.FgRed).SprintFunc()
				} else if group == highlight.Groups["function"] {
					printHl = color.New(color.FgBlue).SprintFunc()
				} else if group == highlight.Groups["constant.string"] {
					printHl = color.New(color.FgHiCyan).SprintFunc()
				} else if group == highlight.Groups["constant.specialChar"] {
					printHl = color.New(color.FgHiMagenta).SprintFunc()
				} else if group == highlight.Groups["constant.number"] {
					printHl = color.New(color.FgHiBlue).SprintFunc()
				} else if group == highlight.Groups["constant.bool"] {
					printHl = color.New(color.FgHiBlue).SprintFunc()
				} else if group == highlight.Groups["symbol.brackets"] {
					printHl = color.New(color.FgRed).SprintFunc()
				} else if group == highlight.Groups["type"] {
					printHl = color.New(color.FgYellow).SprintFunc()
				} else if group == highlight.Groups["comment"] {
					printHl = color.New(color.FgHiBlack).SprintFunc()
				} else {
					printHl = color.New(color.Reset).SprintFunc()
				}
			}

			highlightedCode += printHl(string(c))
		}

		highlightedCode += "\n"
	}

	return highlightedCode
}
