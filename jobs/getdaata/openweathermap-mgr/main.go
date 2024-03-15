package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/timotewb/cpu/jobs/getdata/common"
	"github.com/timotewb/cpu/jobs/getdata/openweathermap-mgr/app"
)

func main() {
	var configDir string
	var help bool
	// Define CLI flags in shrot and long form
	flag.StringVar(&configDir, "c", "", "Path where configuration files are stored (shorthand)")
	flag.StringVar(&configDir, "config", "", "Path where configuration files are stored")
	flag.BoolVar(&help, "h", false, "Show usage instructions (shorthand)")
	flag.BoolVar(&help, "help", false, "Show usage instructions")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "----------------------------------------------------------------------------------------")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Pass -c to specify where the configuration files are stored:")
		fmt.Fprintln(os.Stderr, "  -c\t\tstring\n  --config")
		fmt.Fprintln(os.Stderr, "  \tPath where configuration files are stored")
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

	// Read All Config
	allConfig, err := common.ReadAllConfig(configDir)
	if err != nil {
		log.Fatalf("function ReadAllConfig() failed: %v", err)
		return
	}

	// Read Job Config
	cityList, err := app.ReadCityList(configDir)
	if err != nil {
		log.Fatalf("function ReadCityList() failed: %v", err)
		return
	}

	fmt.Println(allConfig)
	fmt.Println(cityList)
}
