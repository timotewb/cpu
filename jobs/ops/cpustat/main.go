package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/timotewb/cpu/jobs/ops/cpustat/app"
)

func main(){
	// var storageAccount string
	// var containerName string
	// var blobName string

	// get path of code
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalln("unable to get the current file path")
	}
	fullPath := filepath.Dir(file)

	conf, err := app.ReadJobConfig(fullPath)
	if err != nil {
		log.Fatalf("from app.ReadJobConfig(): %s", err)
	}

	// read in blob data
	app.EnvVariables(".env")
	client := app.BlobClient(conf.Azure.StorageAccountName)
	ctx := context.Background()
	blob := app.ReadBlob(ctx, client, conf.Azure.ContainerName, conf.Azure.BlobName)

	// for each server
	for i := range conf.Servers{
		fmt.Println(conf.Servers[i].IPAddress)
		if app.Ping(conf.Servers[i].IPAddress){
			fmt.Println("alive")
		}
	}
	// ping, if true stat

	app.EnvVariables(".env")
	client := app.BlobClient(conf.Azure.StorageAccountName)
	ctx := context.Background()
	app.WriteToBlob(ctx, client, conf.Azure.ContainerName, conf.Azure.BlobName, []byte("Maybe some Blob!\n"))

}