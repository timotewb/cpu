package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/timotewb/cpu/jobs/getdata/common/config"
	"github.com/timotewb/cpu/jobs/getdata/openweathermap-mgr/app"
)

type BodyType struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

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
	allConfig, err := config.ReadAllConfig(configDir)
	if err != nil {
		log.Fatalf("function ReadAllConfig() failed: %v", err)
	}

	// Read Job Config
	jobConfig, err := app.ReadJobConfig(configDir)
	if err != nil {
		log.Fatalf("function ReadJobConfig() failed: %v", err)
	}

	// ReadCity List
	cityList, err := app.ReadCityList(configDir)
	if err != nil {
		log.Fatalf("function ReadCityList() failed: %v", err)
	}

	var body BodyType
	var cityIDs string
	var groupSize int
	var g int

	body.Name = "openweathermap"
	groupSize = jobConfig.GroupSize
	g = 0

	for i := 0; i < len(cityList); i++ {

		if g == groupSize || i == len(cityList) {

			body.Args = append(body.Args, "-c", configDir, "-i", cityIDs)

			// Prepare the JSON body
			jsonBody, err := json.Marshal(body)
			if err != nil {
				log.Fatalf("error marshalling body to JSON: %v", err)
			}

			fmt.Println(string(jsonBody))
			// Create a new HTTP request
			req, err := http.NewRequest("POST", allConfig.APIHost, bytes.NewBuffer(jsonBody))
			if err != nil {
				log.Fatalf("error creating new HTTP request: %v", err)
			}

			// Set the Content-Type header to application/json
			req.Header.Set("Content-Type", "application/json")

			// Send the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatalf("error sending POST request: %v", err)
			}
			defer resp.Body.Close()

			// Check the response status
			if resp.StatusCode != http.StatusOK {
				bytes, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Fatalf("error reading response body: %v", err)
				}
				log.Fatalf("received non-200 response: %s - body: %s", resp.Status, string(bytes))
			}

			body.Args = make([]string, 0)
			cityIDs = ""
			g = 0
			time.Sleep(time.Duration(jobConfig.SleepTimeSec) * time.Second)

		} else {
			if cityIDs == "" {
				cityIDs = strconv.Itoa(int(cityList[i].Id))
			} else {
				cityIDs = cityIDs + "," + strconv.Itoa(int(cityList[i].Id))
			}
			g += 1
		}
	}

}
