package app

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type JobConfig struct {
	// List of urls to be processed
	CamerasURL  string `json:"cameras_url"`
	ChargersURL string `json:"chargers_url"`
}

// ReadJobConfig reads and returns the application configuration from a JSON file.
// It returns an error if the file cannot be opened or if the JSON cannot be unmarshalled.
func ReadJobConfig(configDir string) (JobConfig, error) {
	var config JobConfig

	configPath := filepath.Join(filepath.Join(configDir, "cpustat.json"))
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
