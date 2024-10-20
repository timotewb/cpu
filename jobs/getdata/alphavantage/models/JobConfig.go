package models

type JobConfig struct {
	// List of urls to be processed
	Symbols  []string `json:"symbols"`
	APIKey string `json:"api_key"`
}