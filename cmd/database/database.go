// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

// The database is used to store problem statements and other data related to the problems to reduce
// the number of API calls made to the Kilonova servers.

package database

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"kncli/cmd/problems"
	"kncli/internal"
	"os"
	"path/filepath"
	"strconv"

	"github.com/charmbracelet/huh/spinner"
)

var DatabaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Interact with the problem database",
}

var CreateDBCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates the problem database.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		action := func() { CreateDB() }
		if err := spinner.New().Title("Waiting ...").Action(action).Run(); err != nil {
			internal.LogError(err)
			return
		}
	},
}

var DeleteDBCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes the problem database.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		deleteDB()
	},
}

var RefreshDBCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refreshes the problem database.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		action := func() { CreateDB() }
		if err := spinner.New().Title("Waiting ...").Action(action).Run(); err != nil {
			refreshDB()
			return
		}
	},
}

func init() {
	DatabaseCmd.AddCommand(CreateDBCmd)
	DatabaseCmd.AddCommand(DeleteDBCmd)
	DatabaseCmd.AddCommand(RefreshDBCmd)
}

func CreateDB() {
	db := internal.DBOpen()

	createProblemTableSQL := `CREATE TABLE IF NOT EXISTS problems (
id INTEGER PRIMARY KEY,
name TEXT,
timelimit FLOAT,
memorylimit INTEGER,
sourcesize INTEGER,
credits TEXT,
statement TEXT
);`
	_, err := db.Exec(createProblemTableSQL)
	if err != nil {
		internal.LogError(err)
	}

	internal.DBClose(db)

	println("Database created successfully.")

	refreshDB()

}

func deleteDB() {
	if !internal.DBExists() {
		fmt.Println(`Database file does not exist.`)
	}
	dbFile := filepath.Join(internal.GetConfigDir(), internal.PROBLEMSDATABASE)
	err := os.Remove(dbFile)
	if err != nil {
		internal.LogError(err)
	}
	println("Database deleted successfully.")
}

func refreshDB() {
	if !internal.DBExists() {
		fmt.Println(`Database file does not exist.`)
	}

	url := fmt.Sprintf(internal.URL_PROBLEM, "get")
	data, err := internal.PostJSON[internal.ProblemList](url, nil)
	if err != nil {
		internal.LogError(err)
	}

	db := internal.DBOpen()

	for _, problem := range data.Data {
		query := `SELECT EXISTS(SELECT 1 FROM problems WHERE id = $1)`
		var exists bool
		err := db.QueryRow(query, problem.Id).Scan(&exists)
		if err != nil {
			internal.LogError(err)
		}

		if exists {
			continue
		}

		statement := problems.GetStatementOnline(strconv.Itoa(problem.Id), "RO")
		if statement == internal.NOLANG {
			statement = problems.GetStatementOnline(strconv.Itoa(problem.Id), "EN")
		}

		insertSQL := `INSERT INTO problems (id, name, sourcesize, timelimit, memorylimit, credits, statement)
 VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (id) DO NOTHING;`

		_, err = db.Exec(insertSQL, problem.Id, problem.Name, problem.SourceSize,
			problem.Time, problem.MemoryLimit, problem.SourceCredits, statement)
		if err != nil {
			internal.LogError(fmt.Errorf("error inserting problem info: %v", err))
		}

	}

	internal.DBClose(db)

	fmt.Println("Database refreshed successfully.")
}
