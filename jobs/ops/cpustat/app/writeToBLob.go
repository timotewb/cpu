package app

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func WriteToBlob(ctx context.Context, client *azblob.Client, containerName, blobName string, blobData []byte){

	_, err := client.UploadBuffer(ctx, containerName, blobName, blobData, &azblob.UploadBufferOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}
}