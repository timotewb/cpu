package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/timotewb/cpu/jobs/ops/cpustat/app"
)

func main(){
	var storageAccount string
	var containerName string
	var blobName string
	var help bool
	// Define CLI flags in shrot and long form
	flag.StringVar(&storageAccount, "s", "", "Storage Account name (shorthand)")
	flag.StringVar(&storageAccount, "storage-accounts", "", "Storage Account name")
	flag.StringVar(&containerName, "c", "", "Container name (shorthand)")
	flag.StringVar(&containerName, "container", "", "Container name")
	flag.StringVar(&blobName, "b", "", "Blob name (shorthand)")
	flag.StringVar(&blobName, "blob", "", "Blob name")
	flag.BoolVar(&help, "h", false, "Show usage instructions (shorthand)")
	flag.BoolVar(&help, "help", false, "Show usage instructions")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "----------------------------------------------------------------------------------------")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Pass -s to specify the Storage Account name:")
		fmt.Fprintln(os.Stderr, "  -s\t\tstring\n  --storage-account")
		fmt.Fprintln(os.Stderr, "Pass -c to specify the Container name:")
		fmt.Fprintln(os.Stderr, "  -c\t\tstring\n  --container")
		fmt.Fprintln(os.Stderr, "Pass -b to specify the Blob name:")
		fmt.Fprintln(os.Stderr, "  -b\t\tstring\n  --blob")
		fmt.Fprintln(os.Stderr, "\n  -h\n  --help")
		fmt.Fprintln(os.Stderr, "  \tShow usage instructions")
		fmt.Fprintln(os.Stderr, "----------------------------------------------------------------------------------------")
	}
	flag.Parse()

	// Print the Help docuemntation to the terminal if user passes help flag
	if help {
		flag.Usage()
		return
	}

	if storageAccount == "" {
		log.Fatalf("No Storage Account name provided.")
	}
	if containerName == "" {
		log.Fatalf("No Container name provided.")
	}
	if blobName == "" {
		log.Fatalf("No Blob name provided.")
	}

	app.EnvVariables(".env")

	// for each server
	// ping, if true stat

	app.WriteToBlob(storageAccount, containerName, blobName, []byte("Maybe some Blob!\n"))

}