// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	problem "kilocli/cmd/problems"
	utility "kilocli/cmd/utility"
	"math/rand/v2"

	"github.com/spf13/cobra"
)

var GetRandPbCmd = &cobra.Command{
	Use:   "random",
	Short: "Get random problem to solve.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		getRandomProblemID()
	},
}

func problemCount() int {
	SearchData := map[string]interface{}{
		"name_fuzzy": "",
		"offset":     0,
	}

	JSONData, err := json.Marshal(SearchData)
	if err != nil {
		utility.LogError(fmt.Errorf("error marshaling JSON: %v", err))
		return 0
	}

	ResponseBody, err := utility.MakePostRequest(utility.URL_SEARCH, bytes.NewBuffer(JSONData), utility.RequestJSON)
	if err != nil {
		utility.LogError(err)
		return 0
	}

	var Data problem.SearchResponse
	err = json.Unmarshal(ResponseBody, &Data)
	if err != nil {
		utility.LogError(fmt.Errorf("error unmarshaling JSON: %v", err))
		return 0
	}

	return Data.Data.Count

}

func getRandomProblemID() {
	count := problemCount()
	if count == 0 {
		fmt.Println("No problems available.")
		return
	}

	randomID := rand.IntN(count) + 1
	fmt.Printf("Your random problem's ID: #%d\n", randomID)
}
