package app

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func GetXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("get error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("status error: %v", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("read body: %v", err)
	}

	return data, nil
}

func GetOrCreateSQLiteDB(conf AllConfig, jobName string) (*sql.DB, error) {

	err := os.MkdirAll(conf.StagingPath, 0777) // 0755 is the permission for the directory
	if err != nil {
		return nil, fmt.Errorf("unable to create staging dir: %v", err)
	}

	// Check if the SQLite database file exists and its size
	files, err := os.ReadDir(conf.StagingPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read staging dir: %v", err)
	}

	var mostRecentDBPath string
	var mostRecentTime time.Time
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		dbPath := filepath.Join(conf.StagingPath, file.Name())
		// Assuming dbPath is a string containing the path to your file
		fileInfo, err := os.Stat(dbPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get file info: %v", err)
		}

		// Get the size of the file in bytes
		sizeInBytes := fileInfo.Size()

		// Convert the size to megabytes
		sizeInMegabytes := int(sizeInBytes) / (1024 * 1024)

		if filepath.Ext(file.Name()) == ".db" && filepath.Base(file.Name())[:3] == jobName {

			if sizeInMegabytes < conf.SQLiteMaxSizeMB {
				timestampStr := file.Name()[4:18] // Assuming the format is "rss_YYYYMMDDHHMMSS.db"
				timestamp, err := time.Parse("20060102150405", timestampStr)
				if err != nil {
					return nil, fmt.Errorf("failed to parse timestamp: %v", err)
				}

				if timestamp.After(mostRecentTime) {
					mostRecentTime = timestamp
					mostRecentDBPath = dbPath
				}
			} else {
				// move db file to loading dir
				err = MoveFile(dbPath, conf.LoadingPath)
				if err != nil {
					return nil, fmt.Errorf("failed to move db file to loading dir: %v", err)
				}
			}
		}
	}

	if mostRecentDBPath == "" {
		timestamp := time.Now().Format("20060102150405")
		newDBName := fmt.Sprintf("%v_%s.db", jobName, timestamp)
		dbPath := filepath.Join(conf.StagingPath, newDBName)
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			return nil, fmt.Errorf("unable to open new db: %v", err)
		}
		db.SetMaxOpenConns(1)
		return db, nil
	} else {
		db, err := sql.Open("sqlite3", mostRecentDBPath)
		if err != nil {
			return nil, fmt.Errorf("unable to open existing db: %v", err)
		}
		db.SetMaxOpenConns(1)
		return db, nil
	}
}

func MoveFile(srcFilePath, destDirPath string) error {
	// Check if the destination directory exists
	_, err := os.Stat(destDirPath)
	if os.IsNotExist(err) {
		// Create the destination directory
		err = os.MkdirAll(destDirPath, 0755) // 0755 is the permission for the directory
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check if directory exists: %v", err)
	}

	// Move the file
	destFilePath := filepath.Join(destDirPath, filepath.Base(srcFilePath))
	err = os.Rename(srcFilePath, destFilePath)
	if err != nil {
		return fmt.Errorf("failed to move file: %v", err)
	}

	return nil
}
