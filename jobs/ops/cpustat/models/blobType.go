package models

type BlobType struct {
	Name  string `json:"name"`
	Platform  string `json:"platform"`
	LastUpdated  string `json:"last_updated"`
	LoadAverage  string `json:"load_average"`
	RunningProcs  int64 `json:"running_procs"`
	UpTime  string `json:"uptime"`
	MemoryUsedPct float64 `json:"mem_used_pct"`
	CPUCoreCount int64 `json:"cpu_core_count"`
	CPUModel string `json:"cpu_model"`
}