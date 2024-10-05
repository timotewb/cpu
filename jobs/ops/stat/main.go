package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/timotewb/cpu/jobs/ops/stat/app"
)

func main(){
	var opSysType string
	var help bool
	// Define CLI flags in shrot and long form
	flag.StringVar(&opSysType, "o", "", "Operating System type (shorthand)")
	flag.StringVar(&opSysType, "op-sys-type", "", "Operating System type")
	flag.BoolVar(&help, "h", false, "Show usage instructions (shorthand)")
	flag.BoolVar(&help, "help", false, "Show usage instructions")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "----------------------------------------------------------------------------------------")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Pass -o to specify the Operating System type:")
		fmt.Fprintln(os.Stderr, "  -o\t\tstring\n  --op-sys-type")
		fmt.Fprintln(os.Stderr, "  valid values:\t'linux'")
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

	if opSysType == "linux"{
		returnData := app.LinuxStat()
		jsonData, err := json.Marshal(returnData)
		if err != nil {
			log.Printf("from main() error json.Marshal(): %v", err)
		}
		fmt.Println(string(jsonData))
	}
	
}