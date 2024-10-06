package app

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/timotewb/cpu/jobs/ops/cpustat/models"
)

// ReadJobConfig reads and returns the application configuration from a JSON file.
// It returns an error if the file cannot be opened or if the JSON cannot be unmarshalled.
func ReadJobConfig(fullPath string) (models.JobConfig, error) {
	var config models.JobConfig

	configPath := filepath.Join(fullPath, "config.json")
	// Open the configuration file
	fileObj, err := os.Open(configPath)
	if err != nil {
		return config, err
	}
	defer fileObj.Close()

	// Read the file content
	content, err := io.ReadAll(fileObj)
	if err != nil {
		return config, err
	}

	// Decode the JSON content
	if err := json.Unmarshal(content, &config); err != nil {
		return config, err
	}
	return config, nil
}
