package helper

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/timotewb/cpu/jobs/getdata/common/config"
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

func GetOrCreateSQLiteDB(conf config.AllConfig, jobName string) (*sql.DB, string, error) {

	err := os.MkdirAll(conf.StagingPath, 0777)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create staging dir: %v", err)
	}

	// Check if the SQLite database file exists and its size
	files, err := os.ReadDir(conf.StagingPath)
	if err != nil {
		return nil, "", fmt.Errorf("unable to read staging dir: %v", err)
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
			return nil, "", fmt.Errorf("failed to get file info: %v", err)
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
					return nil, "", fmt.Errorf("failed to parse timestamp: %v", err)
				}

				if timestamp.After(mostRecentTime) {
					mostRecentTime = timestamp
					mostRecentDBPath = dbPath
				}
			} else {
				// move db file to loading dir
				err = MoveFile(dbPath, conf.LoadingPath)
				if err != nil {
					return nil, "", fmt.Errorf("failed to move db file to loading dir: %v", err)
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
			return nil, "", fmt.Errorf("unable to open new db: %v", err)
		}
		db.SetMaxOpenConns(1)

		return db, dbPath, nil
	} else {
		db, err := sql.Open("sqlite3", mostRecentDBPath)
		if err != nil {
			return nil, "", fmt.Errorf("unable to open existing db: %v", err)
		}
		db.SetMaxOpenConns(1)
		return db, mostRecentDBPath, nil
	}
}

func MoveFile(srcFilePath, destDirPath string) error {
	// Check if the destination directory exists
	_, err := os.Stat(destDirPath)
	if os.IsNotExist(err) {
		// Create the destination directory
		err = os.MkdirAll(destDirPath, 0777)
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

	// Set permissions to 777 for the moved file
	err = os.Chmod(destFilePath, 0777)
	if err != nil {
		return fmt.Errorf("failed to set permissions on moved file: %v", err)
	}

	return nil
}

func ParseDate(dateStr string) (string, error) {
	// Define the expected input formats
	formats := []string{
		"2006-01-02T15:04:05Z",            // ISO 8601 format
		"Mon, 02 Jan 2006 15:04:05 -0700", // RFC 1123 format
		// Add more formats as needed
	}

	var parsedTime time.Time
	var err error

	// Try parsing the date string with each format
	for _, format := range formats {
		parsedTime, err = time.Parse(format, dateStr)
		if err == nil {
			break // Successfully parsed, break the loop
		}
	}

	if err != nil {
		// None of the formats matched
		return "", fmt.Errorf("failed to parse date: %v", err)
	}

	// Format the parsed time into a consistent output format, including timezone
	consistentFormat := "2006-01-02 15:04:05 -0700"
	return parsedTime.Format(consistentFormat), nil
}
