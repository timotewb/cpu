package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

// Config holds the application configuration.
type AllConfig struct {
	// StagingPath is the data output path.
	StagingPath string `json:"staging_dir"`
	// LoadingPath is the path where data over te max size is moved to.
	LoadingPath string `json:"loading_dir"`
	// JobSQLiteMaxSizeMBList is the max size of a sqlite db in MB before a new db file is created.
	SQLiteMaxSizeMB int `json:"sqlite_max_size_mb"`
	//APIHost is the ip and path for calling api
	APIHost string `json:"api_host"`
}

// ReadAllConfig reads and returns the application configuration from a JSON file.
// It returns an error if the file cannot be opened or if the JSON cannot be unmarshalled.
func ReadAllConfig(configDir string) (AllConfig, error) {
	var config AllConfig
	// executablePath, err := os.Executable()
	// if err != nil {
	// 	return config, err
	// }
	// // Get the directory of the executable
	// executableDir := filepath.Dir(executablePath)
	// // Construct the path to the config.json file in the same directory
	// configPath := filepath.Join(executableDir, "all.json")
	configPath := filepath.Join(filepath.Join(configDir, "all.json"))

	// Open the configuration file
	file, err := os.Open(configPath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		return config, err
	}

	// Decode the JSON content
	if err := json.Unmarshal(content, &config); err != nil {
		return config, err
	}
	return config, nil
}
