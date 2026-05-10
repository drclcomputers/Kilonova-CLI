// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

// The database is used to store problem statements and other data related to the problems to reduce
// the number of API calls made to the Kilonova servers.

package database

import (
	"fmt"
	"github.com/charmbracelet/huh/spinner"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"kncli/cmd/problems"
	"kncli/internal"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var DatabaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Interact with the problem database",
}

var CreateDBCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates the problem database. (online)",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		action := func() { CreateDB() }
		if err := spinner.New().Title("Please wait...").Action(action).Run(); err != nil {
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
	Short: "Refreshes the problem database. (online)",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		action := func() { refreshDB() }
		if err := spinner.New().Title("Please wait...").Action(action).Run(); err != nil {
			internal.LogError(err)
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
	if db == nil {
		internal.LogError(fmt.Errorf("database is not available"))
		return
	}
	defer internal.DBClose()

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
		return
	}

	println("Database created successfully.")

	refreshDB()

}

func deleteDB() {
	if !internal.DBExists() {
		fmt.Println(`Database file does not exist.`)
	}

	lastrefresh := filepath.Join(internal.GetConfigDir(), internal.LASTREFRESHDB)
	_ = os.Remove(lastrefresh)

	dbFile := filepath.Join(internal.GetConfigDir(), internal.PROBLEMSDATABASE)
	_ = os.Remove(dbFile)

	println("Database deleted successfully.")
}

func refreshDB() {
	if !internal.DBExists() {
		fmt.Println(`Database file does not exist. Create it using 'database create'.`)
		return
	}

	url := fmt.Sprintf(internal.URL_PROBLEM, "get")
	data, err := internal.PostJSON[internal.ProblemList](url, struct{}{})
	if err != nil {
		internal.LogError(err)
		return
	}

	db := internal.DBOpen()
	if db == nil {
		internal.LogError(fmt.Errorf("database is not available"))
		return
	}
	defer internal.DBClose()

	tx, err := db.Begin()
	if err != nil {
		internal.LogError(fmt.Errorf("failed to begin database refresh transaction: %w", err))
		return
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	existsStmt, err := tx.Prepare(`SELECT EXISTS(SELECT 1 FROM problems WHERE id = ?)`)
	if err != nil {
		internal.LogError(fmt.Errorf("failed to prepare existence query: %w", err))
		return
	}
	defer existsStmt.Close()

	insertStmt, err := tx.Prepare(`INSERT INTO problems (id, name, sourcesize, timelimit, memorylimit, credits, statement)
 VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (id) DO NOTHING;`)
	if err != nil {
		internal.LogError(fmt.Errorf("failed to prepare insert query: %w", err))
		return
	}
	defer insertStmt.Close()

	for _, problem := range data.Data {
		var exists bool
		err := existsStmt.QueryRow(problem.Id).Scan(&exists)
		if err != nil {
			internal.LogError(err)
			continue
		}

		if exists {
			continue
		}

		statement := problems.GetStatementOnlineQuiet(strconv.Itoa(problem.Id), "RO", 2)
		if statement == internal.NOLANG {
			statement = problems.GetStatementOnlineQuiet(strconv.Itoa(problem.Id), "EN", 2)
		}

		_, err = insertStmt.Exec(problem.Id, problem.Name, problem.SourceSize,
			problem.Time, problem.MemoryLimit, problem.SourceCredits, statement)
		if err != nil {
			internal.LogError(fmt.Errorf("error inserting problem info: %v", err))
		}

	}

	if err := tx.Commit(); err != nil {
		internal.LogError(fmt.Errorf("failed to commit database refresh transaction: %w", err))
		return
	}
	committed = true

	currentTime := time.Now()
	filePath := filepath.Join(internal.GetConfigDir(), internal.LASTREFRESHDB)
	layout := time.RFC3339

	if err := os.WriteFile(filePath, []byte(currentTime.Format(layout)), 0644); err != nil {
		internal.LogError(err)
	}

	fmt.Println("Database refreshed successfully.")
}
