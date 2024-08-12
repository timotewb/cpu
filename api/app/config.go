package app

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

// Config holds the application configuration.
type Config struct {
	// AppPath is the path to the application directory.
	AppPath string `json:"app_path"`
	// JobList is a list of valid job names.
	JobList []string `json:"job_list"`
	// Notificaitons ettings
	Notification Notification`json:"notification"`
}
type Notification struct {
    To         string   `json:"to"`
    From       string   `json:"from"`
    SMTP       SMTP     `json:"smtp"`
}
type SMTP struct {
    Host      string  `json:"host"`
    Port      int     `json:"port"`
    Username  string  `json:"username"`
    Password  string  `json:"password"`
}

// ReadConfig reads and returns the application configuration from a JSON file.
// It returns an error if the file cannot be opened or if the JSON cannot be unmarshalled.
func ReadConfig(configDir string) (Config, error) {
	var config Config
	// Open the configuration file
	file, err := os.Open(filepath.Join(configDir, "config.json"))
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
