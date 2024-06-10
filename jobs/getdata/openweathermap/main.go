package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/timotewb/cpu/jobs/getdata/common/config"
	"github.com/timotewb/cpu/jobs/getdata/common/helper"
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
	allConfig, err := config.ReadAllConfig(configDir)
	if err != nil {
		log.Fatalf("function ReadAllConfig() failed: %v", err)
	}

	// Read Job Config
	jobConfig, err := app.ReadJobConfig(configDir)
	if err != nil {
		log.Fatalf("function ReadJobConfig() failed: %v", err)
	}

	fmt.Println(allConfig)

	// make call to api
	resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/group?id=%v&appid=%s", cityIDs, jobConfig.APIKey))
	if err != nil {
		log.Fatalf("function http.Get() failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("error reading response body: %v", err)
		}
		log.Fatalf("received non-200 response: %s - body: %s", resp.Status, string(bytes))
	}

	// get sqlite db
	db, dbPath, err := helper.GetOrCreateSQLiteDB(allConfig, "journeys_nzta")
	if err != nil {
		log.Fatalf("from Cameras(): function GetOrCreateSQLiteDB() failed: %v", err)
	}
	defer db.Close()

	fmt.Println(dbPath)

	// create target table if not exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS openweathermap (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		latitude REAL,
		longitude REAL,
		weather_details_id INTEGER,
		base STRING,
		teamp REAL,
		feels_like REAL,
		pressure INTEGER,
		humidity INTEGER,
		temp_min REAL,
		temp_max REAL,
		sea_level REAL,
		grnd_level REAL,
		visibility INTEGER,
		wind_speed REAL,
		wind_deg REAL,
		clouds INTEGER,
		rain1h INTEGER,
		rain3h INTEGER,
		snow1h INTEGER,
		snow3h INTEGER,
		dt INTEGER,
		sys_type INTEGER,
		sys_id INTEGER,
		sys_message STRING,
		sys_country STRING,
		sunrise INTEGER,
		sunset INTEGER,
		timezone INTEGER,
		id0 INTEGER,
		name STRING,
		cod INTEGER
	)`)
	if err != nil {
		log.Fatalf("failed to create table 'openweathermap': %v\n", err)
	}
	// chargers_access_locations
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS weather_details (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		weather_details_id INTEGER,
		id0 INTEGER,
		main STRING,
		description STRING,
		icon STRING,
		FOREIGN KEY(weather_details_id) REFERENCES openweathermap(id) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Fatalf("from Chargers(): failed to create table 'chargers_access_locations': %v\n", err)
	}

	// convert respose to string then return
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("error reading response body (StatusOK): %v", err)
		}
		var data app.ResponseWrapper
		json.Unmarshal(bodyBytes, &data)
		if data.Cnt > 0 {
			// load data to db
			// b, err := json.Marshal(&data.List)
			// if err != nil {
			// 	fmt.Println(err)
			// 	log.Fatalf("function json.Marshal() failed: %v", err)
			// }

			for i := 0; i < len(data.List); i++ {
				fmt.Println(data.List[i].Clouds)
			}
		}
	}
}
