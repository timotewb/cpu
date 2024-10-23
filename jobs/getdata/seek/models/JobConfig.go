package models

type JobConfig struct {
	// List of urls to be processed
	PageCount  int `json:"page_count"`
	URLs  []string `json:"urls"`
}