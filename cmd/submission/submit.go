// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package submission

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"kncli/internal"
	"mime/multipart"
	"os"

	"github.com/charmbracelet/huh/spinner"
)

func uploadCode(id, language, file, contestID string) {
	codeFile, err := os.Open(file)
	if err != nil {
		internal.LogError(fmt.Errorf("failed to open code file %s: %w", file, err))
		return
	}
	defer codeFile.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	if err := writeFormFields(writer, id, language); err != nil {
		internal.LogError(err)
		return
	}

	if err := writeCodeFile(writer, file, codeFile); err != nil {
		internal.LogError(err)
		return
	}

	if contestID != "NO" {
		if err := writeFormContest(writer, contestID); err != nil {
			internal.LogError(err)
			return
		}
	}

	contentType := writer.FormDataContentType()

	if err := writer.Close(); err != nil {
		internal.LogError(fmt.Errorf("failed to close writer: %w", err))
		return
	}

	ResponseBody, err := internal.MakePostRequest(internal.URL_SUBMIT, &requestBody, internal.RequestMultipartForm, contentType)
	if err != nil {
		internal.LogError(fmt.Errorf("error submitting code: %w", err))
		return
	}

	var data Submit
	if err := json.Unmarshal(ResponseBody, &data); err != nil {
		handleSubmissionError(ResponseBody)
		return
	}

	fmt.Printf("Submission sent: %s\nSubmission ID: %d\n", data.Status, data.Data)
	checkSubmissionStatus(data.Data)
}

func writeFormFields(writer *multipart.Writer, id, language string) error {
	if err := writer.WriteField("problem_id", id); err != nil {
		return fmt.Errorf("failed to write problem_id field: %w", err)
	}
	if err := writer.WriteField("language", language); err != nil {
		return fmt.Errorf("failed to write language field: %w", err)
	}
	return nil
}

func writeFormContest(writer *multipart.Writer, id string) error {
	if err := writer.WriteField("contestID", id); err != nil {
		return fmt.Errorf("failed to write contestID field: %w", err)
	}
	return nil
}

func writeCodeFile(writer *multipart.Writer, file string, codeFile *os.File) error {
	fileWriter, err := writer.CreateFormFile("code", file)
	if err != nil {
		return fmt.Errorf("failed to create form file for code: %w", err)
	}
	if _, err := io.Copy(fileWriter, codeFile); err != nil {
		return fmt.Errorf("failed to copy code file content: %w", err)
	}
	return nil
}

func handleSubmissionError(ResponseBody []byte) {
	var dataerr SubmitError
	if err := json.Unmarshal(ResponseBody, &dataerr); err != nil {
		internal.LogError(fmt.Errorf("failed to parse error response: %w", err))
		return
	}
	internal.LogError(fmt.Errorf("status: %s\nmessage: %s", dataerr.Status, dataerr.Data))
}

func checkSubmissionStatus(submissionID int) {
	url := fmt.Sprintf(internal.URL_LATEST_SUBMISSION, submissionID)

	action := func() {
		var dataLatestSubmit LatestSubmission
		for {
			ResponseBody, err := internal.MakeGetRequest(url, nil, internal.RequestNone)
			if err != nil {
				internal.LogError(fmt.Errorf("failed to get submission status: %w", err))
				return
			}

			if err := json.Unmarshal(ResponseBody, &dataLatestSubmit); err != nil {
				internal.LogError(fmt.Errorf("failed to parse latest submission response: %w", err))
				return
			}

			if dataLatestSubmit.Data.Status == "finished" {
				break
			}

			fmt.Print(".")
		}

		if dataLatestSubmit.Data.CompileError {
			fmt.Println("\nCompilation failed! Score: 0")
		} else {
			fmt.Printf("\nSuccess! Score: %d\n", dataLatestSubmit.Data.Score)
		}
	}

	if err := spinner.New().Title("Please wait...").Action(action).Run(); err != nil {
		internal.LogError(fmt.Errorf("spinner error: %w", err))
		return
	}
}
