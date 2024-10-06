package models

type JobConfig struct {
	Servers []ServerType `json:"servers"`
	Azure AzureType `json:"azure"`
}

type ServerType struct {
	Name  string `json:"name"`
	IPAddress  string `json:"ip_address"`

}

type AzureType struct {
	StorageAccountName  string `json:"storage_account_name"`
	ContainerName  string `json:"container_name"`
	BlobName  string `json:"blob_name"`
}