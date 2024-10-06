package app

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func BlobClient(storageAccount string) *azblob.Client {
	urlSA :=  fmt.Sprintf("https://%s.blob.core.windows.net/", storageAccount)

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	client, err := azblob.NewClient(urlSA, credential, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	return client
}