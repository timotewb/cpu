package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/timotewb/cpu/jobs/getdata/openweathermap/app"
)

func main() {
	var configDir string
	var cityIDs string
	var help bool
	// Define CLI flags in shrot and long form
	flag.StringVar(&configDir, "c", "", "Path where configuration files are stored (shorthand)")
	flag.StringVar(&configDir, "config", "", "Path where configuration files are stored")
	flag.StringVar(&cityIDs, "i", "", "Comma seperated list of Openweathermap City IDs")
	flag.BoolVar(&help, "h", false, "Show usage instructions (shorthand)")
	flag.BoolVar(&help, "help", false, "Show usage instructions")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "----------------------------------------------------------------------------------------")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Pass -c to specify where the configuration files are stored:")
		fmt.Fprintln(os.Stderr, "  -c\t\tstring\n  --config")
		fmt.Fprintln(os.Stderr, "  \tPath where configuration files are stored")
		fmt.Fprintln(os.Stderr, "\n  -i\t\t int")
		fmt.Fprintln(os.Stderr, "  \tComma seperated list of Openweathermap City IDs")
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
	allConfig, err := app.ReadAllConfig(configDir)
	if err != nil {
		log.Fatalf("function ReadAllConfig() failed: %v", err)
		return
	}

	// Read Job Config
	jobConfig, err := app.ReadJobConfig(configDir)
	if err != nil {
		log.Fatalf("function ReadJobConfig() failed: %v", err)
		return
	}

	fmt.Println(allConfig)

	// make call to api
	resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/group?id=%v&appid=%s", cityIDs, jobConfig.APIKey))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Error: Non 200 status code returned when attempting to retrieve file. Status Code was %v.\n", resp.StatusCode)
		os.Exit(1)
	}

	// convert respose to string then return
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		var data app.ResponseWrapper
		json.Unmarshal(bodyBytes, &data)
		if data.Cnt > 0 {
			// load data to db
			b, err := json.Marshal(&data.List)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Print(string(b))
		}
	}
}
