// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

// The database is used to store problem statements and other data related to the problems to reduce
// the number of API calls made to the Kilonova server. The database is not yet fully functional
// and is still being developed. Please, be patient until more features are added.

package database

import (
	"database/sql"
	utility "kncli/cmd/utility"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var RefreshDBCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refreshes the problem database",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		refreshDB()
	},
}

func returnConfigDir() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		utility.LogError(err)
		return "error"
	}
	configDir := filepath.Join(homedir, utility.CONFIGFOLDER, utility.KNCLIFOLDER)
	err = os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		utility.LogError(err)
		return "error"
	}
	return configDir
}

func dBExistsOrCreate() *sql.DB {
	DBFilename := filepath.Join(returnConfigDir(), utility.PROBLEMSDATABASE)
	db, err := sql.Open("sqlite3", DBFilename)
	if err != nil {
		utility.LogError(err)
		return nil
	}
	return db
}

func refreshDB() {
	db := dBExistsOrCreate()

	createTableSQL := `CREATE TABLE IF NOT EXISTS problems (
id INTEGER PRIMARY KEY,
name TEXT,
source TEXT,
timelimit INTEGER,
memorylimit INTEGER,
maxscore INTEGER,
credits TEXT,
statement TEXT
	);`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		utility.LogError(err)
	}

	println("Database refreshed successfully.")

	defer db.Close()
}
