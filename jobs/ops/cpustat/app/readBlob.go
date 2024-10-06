package app

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/timotewb/cpu/jobs/ops/cpustat/models"
)

func ReadBlob(ctx context.Context, client *azblob.Client, containerName, blobName string) models.BlobType{
	var resp models.BlobType
	// Download the blob
	get, err := client.DownloadStream(ctx, containerName, blobName, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	downloadedData := bytes.Buffer{}
	retryReader := get.NewRetryReader(ctx, &azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(retryReader)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = retryReader.Close()
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := json.Unmarshal([]byte(downloadedData.String()), &resp); err != nil {
		return resp
	}
}