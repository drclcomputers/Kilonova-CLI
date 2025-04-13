// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package submission

type Subtest struct {
	ID         int     `json:"id"`
	Done       bool    `json:"done"`
	Skipped    bool    `json:"skipped"`
	Verdict    string  `json:"verdict"`
	Time       float64 `json:"time"`
	Memory     int     `json:"memory"`
	Percentage int     `json:"percentage"`
	TestID     int     `json:"test_id"`
	Score      int     `json:"score"`
}
type SubmissionData struct {
	Id             int       `json:"id"`
	CreatedAt      string    `json:"created_at"`
	UserID         int       `json:"user_id"`
	ProblemID      int       `json:"problem_id"`
	Language       string    `json:"language"`
	CompileError   bool      `json:"compile_error"`
	ContestID      *int      `json:"contest_id"`
	MaxTime        float64   `json:"max_time"`
	MaxMemory      int       `json:"max_memory"`
	Score          float64   `json:"score"`
	CompileMessage string    `json:"compile_message"`
	Code           string    `json:"code"`
	Subtests       []Subtest `json:"subtests"`
}

type SubmissionDetails struct {
	Status string         `json:"status"`
	Data   SubmissionData `json:"data"`
}

type SubmissionList struct {
	Data struct {
		Submissions []SubmissionData `json:"submissions"`
		Count       int              `json:"count"`
	} `json:"data"`
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
	ContestID      *int
	Code           string
}
