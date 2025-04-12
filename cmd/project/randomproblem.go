// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package project

import (
	"fmt"
	"kncli/internal"
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

func ProblemCount() int {
	return internal.CountProblemsDB()
}

func getRandomProblemID() {
	count := ProblemCount()
	if count == 0 {
		fmt.Println("No problems available in the database.")
		return
	}

	randomID := rand.IntN(count) + 1
	fmt.Printf("Your random problem's ID: #%d\n", randomID)
}
