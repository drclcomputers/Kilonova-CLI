package internal

import (
	"database/sql"
	"os"
	"path/filepath"
)

func DBExists() bool {
	DBFilename := filepath.Join(GetConfigDir(), PROBLEMSDATABASE)
	if _, err := os.Stat(DBFilename); err == nil {
		return true
	}
	return false
}

func DBOpen() *sql.DB {
	DBFilename := filepath.Join(GetConfigDir(), PROBLEMSDATABASE)
	db, err := sql.Open("sqlite3", DBFilename)
	if err != nil {
		LogError(err)
		return nil
	}
	return db
}

func DBClose(db *sql.DB) {
	_ = db.Close()
}

func CountProblemsDB() int {
	db := DBOpen()

	countSQL := `SELECT COUNT(*) FROM problems;`

	var count int
	err := db.QueryRow(countSQL).Scan(&count)
	if err != nil {
		LogError(err)
	}

	DBClose(db)

	return count
}
