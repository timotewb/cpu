package models

type CallCPUJobType struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}