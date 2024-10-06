package models

type BlobType struct {
	LastUpdated  string `json:"last_updated"`
	Servers []BlobServerType `json:"servers"`

}
type BlobServerType struct {
	Name  string `json:"name"`
	Ping  string `json:"ping"`
	LoadAverage  string `json:"load_average"`
	RunningProcs  string `json:"running_procs"`
	UpTime  string `json:"uptime"`

}