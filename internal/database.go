package internal

import (
	"database/sql"
	"os"
	"path"
	"path/filepath"
	"time"
)

func RefreshOrNotDB() bool {
	currentTime := time.Now()
	filePath := path.Join(GetConfigDir(), LASTREFRESHDB)
	layout := time.RFC3339

	if !FileExists(LASTREFRESHDB) {
		file, err := os.Create(filePath)
		if err != nil {
			LogError(err)
			return false
		}
		defer file.Close()

		_, err = file.WriteString(currentTime.Format(layout))
		if err != nil {
			LogError(err)
			return false
		}
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		LogError(err)
		return false
	}

	parsedTime, err := time.Parse(layout, string(data))
	if err != nil {
		LogError(err)
		return false
	}

	if currentTime.Sub(parsedTime) > 7*24*time.Hour {
		return true
	}

	return false
}

func DBExists() bool {
	return FileExists(PROBLEMSDATABASE)
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

func ProblemExistsDB(ID string) bool {
	db := DBOpen()
	defer db.Close()
	query := `SELECT EXISTS(SELECT 1 FROM problems WHERE CAST(id as TEXT) LIKE ?);`
	var exists bool
	err := db.QueryRow(query, ID).Scan(&exists)
	if err != nil {
		LogError(err)
	}

	if exists {
		return true
	}
	return false
}
