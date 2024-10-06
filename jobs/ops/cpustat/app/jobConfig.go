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
func ReadJobConfig() (models.JobConfig, error) {
	var config models.JobConfig

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Unable to get the current file path")
		return
	}
	fullPath := filepath.Dir(file)

	configPath := filepath.join(fullPath, "config.json")
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
