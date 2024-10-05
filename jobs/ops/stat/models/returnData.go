package models

type ServerModel struct {
	Name     string `json:"name"`
	LoadAverage     string `json:"load_average"`
	RunningProcs     string `json:"running_procs"`
	UpTime     string `json:"uptime"`
}