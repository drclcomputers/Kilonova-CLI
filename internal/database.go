// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package internal

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"
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

var (
	dbInstance *sql.DB
	dbMutex    sync.Mutex
)

func DBOpen() *sql.DB {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if dbInstance != nil {
		return dbInstance
	}
	DBFilename := filepath.Join(GetConfigDir(), PROBLEMSDATABASE)
	db, err := sql.Open("sqlite3", DBFilename)
	if err != nil {
		LogError(fmt.Errorf("database open failed: %w", err))
		return nil
	}
	dbInstance = db
	return dbInstance
}

func DBClose() {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if dbInstance != nil {
		_ = dbInstance.Close()
		dbInstance = nil
	}
}

func CountProblemsDB() int {
	db := DBOpen()
	if db == nil {
		return 0
	}

	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM problems;`).Scan(&count)
	if err != nil {
		LogError(fmt.Errorf("count query failed: %w", err))
	}

	return count
}

func ProblemExistsDB(ID string) bool {
	db := DBOpen()
	if db == nil {
		return false
	}
	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM problems WHERE CAST(id as TEXT) = ?);`, ID).Scan(&exists)
	if err != nil {
		LogError(fmt.Errorf("existence query failed: %w", err))
	}
	return exists
}
