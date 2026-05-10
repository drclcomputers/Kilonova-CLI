// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package project

import (
	"database/sql"
	"fmt"
	"kncli/internal"
	"math/rand/v2"

	"github.com/spf13/cobra"
)

var GetRandPbCmd = &cobra.Command{
	Use:   "random",
	Short: "Get random problem to solve.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		getRandomProblemID()
	},
}

func getRandomProblemID() {
	db := internal.DBOpen()
	if db == nil {
		fmt.Println("Problem database not available. Run 'database create' first.")
		return
	}
	defer internal.DBClose()

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM problems").Scan(&count); err != nil || count == 0 {
		fmt.Println("No problems available in the database.")
		return
	}

	offset := rand.IntN(count)
	var id int
	var name string
	if err := db.QueryRow("SELECT id, name FROM problems LIMIT 1 OFFSET ?", offset).Scan(&id, &name); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No problems available.")
		} else {
			internal.LogError(fmt.Errorf("failed to get random problem: %w", err))
		}
		return
	}

	fmt.Printf("Your random problem: #%d — %s\n", id, name)
}
