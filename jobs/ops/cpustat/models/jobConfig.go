package models

type JobConfig struct {
	Servers []ServerType `json:"servers"`
}

type ServerType {
	Name  string `json:"name"`
	IPAddress  string `json:"ip_address"`

}