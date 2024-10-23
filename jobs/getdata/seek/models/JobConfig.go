package models

type JobConfig struct {
	// List of urls to be processed
	PageCount  int `json:"page_count"`
	URLs  []URLType `json:"urls"`
}

type URLType struct {
	// List of urls to be processed
	Classification  string `json:"classification"`
	URL  string `json:"url"`
}