// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package submission

type SubmissionData struct {
	UserID         int     `json:"user_id"`
	ProblemID      int     `json:"problem_id"`
	Id             int     `json:"id"`
	CreatedAt      string  `json:"created_at"`
	Language       string  `json:"language"`
	Score          float64 `json:"score"`
	MaxMemory      int     `json:"max_memory"`
	MaxTime        float64 `json:"max_time"`
	CompileError   bool    `json:"compile_error"`
	CompileMessage string  `json:"compile_message"`
	Code           string  `json:"code,omitempty"`
}

type SubmissionList struct {
	Data struct {
		Submissions []SubmissionData `json:"submissions"`
		Count       int              `json:"count"`
	} `json:"data"`
}

type SubmissionDetails struct {
	Status string         `json:"status"`
	Data   SubmissionData `json:"data"`
}

type Submit struct {
	Status string `json:"status"`
	Data   int    `json:"data"`
}

type SubmitError struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

type Languages struct {
	Data []struct {
		Name string `json:"internal_name"`
	} `json:"data"`
}

type LatestSubmission struct {
	Status string `json:"status"`
	Data   struct {
		Status       string `json:"status"`
		CompileError bool   `json:"compile_error"`
		Score        int    `json:"score"`
	}
}

type SubmissionDetailsTemplate struct {
	ID             int
	CreatedAt      string
	Language       string
	Score          float64
	MaxMemory      int
	MaxTime        float64
	CompileError   bool
	CompileMessage string
	Code           string
}
