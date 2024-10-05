package app

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func WriteToBlob(storageAccount, containerName, blobName string, blobData []byte){
	urlSA :=  fmt.Sprintf("https://%s.blob.core.windows.net/", storageAccount)
	ctx := context.Background()

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	client, err := azblob.NewClient(urlSA, credential, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = client.UploadBuffer(ctx, containerName, blobName, blobData, &azblob.UploadBufferOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}
}