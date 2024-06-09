package app

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type CityList []struct {
	Id      float64   `json:"id"`
	Name    string    `json:"name"`
	State   string    `json:"state"`
	Country string    `json:"country"`
	Coord   coordType `json:"coord"`
}
type coordType struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

func ReadCityList(configDir string) (CityList, error) {
	var cityList CityList

	configPath := filepath.Join(filepath.Join(configDir, "city.list.json"))
	// Open the configuration file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Decode the JSON content
	if err := json.Unmarshal(content, &cityList); err != nil {
		return nil, err
	}
	return cityList, nil
}

type JobConfig struct {
	// List of urls to be processed
	APIKey    string `json:"api_key"`
	GroupSize int    `json:"group_size"`
}

// ReadJobConfig reads and returns the application configuration from a JSON file.
// It returns an error if the file cannot be opened or if the JSON cannot be unmarshalled.
func ReadJobConfig(configDir string) (JobConfig, error) {
	var config JobConfig

	configPath := filepath.Join(filepath.Join(configDir, "openweathermap.json"))
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
