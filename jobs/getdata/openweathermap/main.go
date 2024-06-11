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
		log.Fatalf("from Openweathermap(): function ReadAllConfig() failed: %v", err)
	}

	// Read Job Config
	jobConfig, err := app.ReadJobConfig(configDir)
	if err != nil {
		log.Fatalf("from Openweathermap(): function ReadJobConfig() failed: %v", err)
	}

	fmt.Println(allConfig)

	// make call to api
	resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/group?id=%v&appid=%s", cityIDs, jobConfig.APIKey))
	if err != nil {
		log.Fatalf("from Openweathermap(): function http.Get() failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("from Openweathermap(): error reading response body: %v", err)
		}
		log.Fatalf("from Openweathermap(): received non-200 response: %s - body: %s", resp.Status, string(bytes))
	}

	// get sqlite db
	db, dbPath, err := helper.GetOrCreateSQLiteDB(allConfig, "journeys_nzta")
	if err != nil {
		log.Fatalf("from Openweathermap(): function GetOrCreateSQLiteDB() failed: %v", err)
	}
	defer db.Close()

	fmt.Println(dbPath)
	//----------------------------------------------------------------------------------------
	// Create tables code
	//----------------------------------------------------------------------------------------
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS openweathermap (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		coord_latitude REAL,
		coord_longitude REAL,
		weather_id INTEGER,
		weather_main STRING,
		weather_description STRING,
		weather_icon STRING,
		base STRING,
		main_temp REAL,
		main_feels_like REAL,
		main_pressure INTEGER,
		main_humidity INTEGER,
		main_temp_min REAL,
		main_temp_max REAL,
		main_sea_level REAL,
		main_grnd_level REAL,
		visibility INTEGER,
		wind_speed REAL,
		wind_deg REAL,
		cloud_all INTEGER,
		rain_1h INTEGER,
		rain_3h INTEGER,
		snow_1h INTEGER,
		snow_3h INTEGER,
		dt INTEGER,
		sys_type INTEGER,
		sys_id INTEGER,
		sys_message STRING,
		sys_country STRING,
		sys_sunrise INTEGER,
		sys_sunset INTEGER,
		timezone INTEGER,
		id0 INTEGER,
		name STRING,
		cod INTEGER
	)`)
	if err != nil {
		log.Fatalf("from Openweathermap(): failed to create table 'openweathermap': %v\n", err)
	}
	// chargers_access_locations
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS weather_details (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		openweathermap_id INTEGER,
		id0 INTEGER,
		main STRING,
		description STRING,
		icon STRING,
		FOREIGN KEY(openweathermap_id) REFERENCES openweathermap(id) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Fatalf("from Openweathermap(): failed to create table 'chargers_access_locations': %v\n", err)
	}

	//----------------------------------------------------------------------------------------
	// Insert Into code
	//----------------------------------------------------------------------------------------
	sqlInsertOWM := `
	INSERT INTO openweathermap (
		coord_latitude, coord_longitude, weather_id, weather_main, weather_description, weather_icon, 
		base, main_temp, main_feels_like, main_pressure, main_humidity, main_temp_min, main_temp_max, 
		main_sea_level, main_grnd_level, visibility, wind_speed, wind_deg, cloud_all, rain_1h, rain_3h, 
		snow_1h, snow_3h, dt, sys_type, sys_id, sys_message, sys_country, sys_sunrise, sys_sunset, 
		timezone, id0, name, cod
	) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);
	`
	// Prepare the statement
	stmtOWM, err := db.Prepare(sqlInsertOWM)
	if err != nil {
		log.Fatalf("from Openweathermap(): failed to prepare insert into openweathermap statement: %s\n", err.Error())
	}
	sqlInsertWD := `
	INSERT INTO weather_details (
		openweathermap_id, id0, main, description, icon
	) VALUES (?,?,?,?,?);
	`
	// Prepare the statement
	stmtWD, err := db.Prepare(sqlInsertWD)
	if err != nil {
		log.Fatalf("from Openweathermap(): failed to prepare insert into weather_details statement: %s\n", err.Error())
	}

	// convert respose to string then return
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("from Openweathermap(): error reading response body (StatusOK): %v", err)
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

				_, err = stmtOWM.Exec(
					data.List[i].Coord.Lat,
					data.List[i].Coord.Lon,
					data.List[i].Base,
					data.List[i].Main.Temp,
					data.List[i].Main.Feels_like,
					data.List[i].Main.Pressure,
					data.List[i].Main.Humidity,
					data.List[i].Main.Temp_min,
					data.List[i].Main.Temp_max,
					data.List[i].Main.Sea_level,
					data.List[i].Main.Grnd_level,
					data.List[i].Visibility,
					data.List[i].Wind.Speed,
					data.List[i].Wind.Deg,
					data.List[i].Clouds.All,
					data.List[i].Rain.Rain1h,
					data.List[i].Rain.Rain3h,
					data.List[i].Snow.Snow1h,
					data.List[i].Snow.Snow3h,
					data.List[i].Dt,
					data.List[i].Sys.Type,
					data.List[i].Sys.Id,
					data.List[i].Sys.Message,
					data.List[i].Sys.Country,
					data.List[i].Sys.Sunrise,
					data.List[i].Sys.Sunset,
					data.List[i].Timezone,
					data.List[i].Id,
					data.List[i].Name,
					data.List[i].Cod,
				)
				if err != nil {
					log.Fatalf("from Openweathermap(): failed to execute insert into openweathermap statement: %s\n", err.Error())
				}
				// Fetch the last inserted ID
				var lastInsertedID int64
				err = db.QueryRow("SELECT last_insert_rowid();").Scan(&lastInsertedID)
				if err != nil {
					log.Fatalf("from Openweathermap(): failed to get last pk: %s\n", err.Error())
				}

				// weather_details
				for a := 0; a < len(data.List[i].Weather); a++ {

					// Execute the statement with the charger record values
					_, err = stmtWD.Exec(
						lastInsertedID,
						data.List[i].Weather[a].Id,
						data.List[i].Weather[a].Main,
						data.List[i].Weather[a].Description,
						data.List[i].Weather[a].Icon,
					)
					if err != nil {
						log.Fatalf("from Openweathermap(): failed to execute 'INSERT INTO weather_details' statement: %s\n", err.Error())
					}
				}
			}
		}
	}
	// remvoe duplicates from table
	_, err = db.Exec(`DELETE FROM openweathermap WHERE id NOT IN (SELECT MIN(id) FROM openweathermap GROUP BY CONCAT(last_edited, created, uniq))`)
	if err != nil {
		log.Fatalf("from Openweathermap(): failed to remove duplicates from openweathermap table: %v\n", err)
	}
	err = os.Chmod(dbPath, 0777)
	if err != nil {
		log.Fatal("from Openweathermap(): failed to set permissions on db file: ", err)
		return
	}
}
