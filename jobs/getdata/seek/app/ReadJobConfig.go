package app

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	m "github.com/timotewb/cpu/jobs/getdata/seek/models"
)

func ReadJobConfig(configDir string) (m.JobConfig, error) {
	var config m.JobConfig

	configPath := filepath.Join(filepath.Join(configDir, "seek.json"))
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
