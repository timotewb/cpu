package app

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func WriteToBlob(storageAccount, containerName, blobName string, blobData []byte){

	url :=  fmt.Sprintf("https://%s.blob.core.windows.net/", storageAccount)

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	client, err := azblob.NewClient(url, credential, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	

	_, err = client.UploadBuffer(context.Background(), containerName, blobName, blobData, &azblob.UploadBufferOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}
}