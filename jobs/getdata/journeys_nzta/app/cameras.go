package app

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/timotewb/cpu/jobs/getdata/common/config"
	"github.com/timotewb/cpu/jobs/getdata/common/helper"
)

func Cameras(allConfig config.AllConfig, jobConfig JobConfig) {

	// get sqlite db
	db, dbPath, err := helper.GetOrCreateSQLiteDB(allConfig, "journeys_nzta")
	if err != nil {
		log.Fatalf("function GetOrCreateSQLiteDB() failed: %v", err)
	}
	defer db.Close()
	fmt.Println(dbPath)

	// create target table if not exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS chargers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		class_name TEXT,
		last_edited TEXT,
		created TEXT,
		site_id INTEGER,
		name TEXT,
		operator TEXT,
		address TEXT,
		is_24_hours INTEGER,
		car_park_count INTEGER,
		has_carpark_cost INTEGER,
		max_time_limit TEXT,
		has_tourist_attraction INTEGER,
		provider_deleted INTEGER,
		hide_from_feed INTEGER,
		region_id INTEGER,
		cameras_id INTEGER,
		record_class_name TEXT,
		external_id INTEGER,
		uniq TEXT,
		type TEXT,
		cameras_id0 TEXT,
		last_updated INTEGER,
		region TEXT,
		owner_name TEXT,
		charging_cost TEXT,
		feature_type TEXT,
	)`)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
	// cameras_access_locations
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cameras_access_locations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		cameras_id INTEGER,
		FOREIGN KEY(cameras_id) REFERENCES cameras(id) ON DELETE CASCADE,
		lat REAL,
		lon REAL
	)`)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
	// cameras_connectors
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cameras_connectors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		cameras_id INTEGER,
		FOREIGN KEY(cameras_id) REFERENCES cameras(id) ON DELETE CASCADE,
		current TEXT,
		kw_rates INTEGER,
		connector_type TEXT,
		operation_status TEXT,
		next_planning_outage TEXT
	)`)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
	// cameras_regions
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cameras_regions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		cameras_id INTEGER,
		FOREIGN KEY(cameras_id) REFERENCES cameras(id) ON DELETE CASCADE,
		regions_id INTEGER
	)`)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	var result ChargersModel
	if jsonBytes, err := helper.GetXML(jobConfig.CamerasURL); err != nil {
		log.Fatal("failed to get json: ", err)
	} else {
		if err := json.Unmarshal(jsonBytes, &result); err != nil {
			log.Fatal("format = 1 unmarshal error: ", err)
		}

		// chargersSQL
		chargersSQL := `
		INSERT INTO chargers (
			class_name ,
			last_edited ,
			created ,
			site_id ,
			name ,
			operator ,
			address ,
			is_24_hours ,
			car_park_count ,
			has_carpark_cost ,
			max_time_limit ,
			has_tourist_attraction ,
			provider_deleted ,
			hide_from_feed ,
			region_id ,
			cameras_id ,
			record_class_name ,
			external_id ,
			uniq ,
			type ,
			cameras_id0 ,
			last_updated ,
			region ,
			owner_name ,
			charging_cost ,
			feature_type ,
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
		`
		chargersSQLPrep, err := db.Prepare(chargersSQL)
		if err != nil {
			log.Printf("failed to prepare statement: %s\n", err.Error())
		}
		defer chargersSQLPrep.Close()

		// accessLocationsSQL
		accessLocationsSQL := `
		INSERT INTO cameras_access_locations (
			cameras_id ,
			lat ,
			lon 
		) VALUES (?,?,?)
		`
		accessLocationsSQLPrep, err := db.Prepare(accessLocationsSQL)
		if err != nil {
			log.Printf("failed to prepare statement: %s\n", err.Error())
		}
		defer accessLocationsSQLPrep.Close()

		// connectorsSQL
		connectorsSQL := `
		INSERT INTO cameras_connectors (
			cameras_id ,
			current ,
			kw_rates ,
			connector_type ,
			operation_status ,
			next_planning_outage 
		) VALUES (?,?,?,?,?,?)
		`
		connectorsSQLPrep, err := db.Prepare(connectorsSQL)
		if err != nil {
			log.Printf("failed to prepare statement: %s\n", err.Error())
		}
		defer connectorsSQLPrep.Close()

		// regionsSQL
		regionsSQL := `
		INSERT INTO cameras_regions (
			cameras_id INTEGER,
			regions_id INTEGER
		) VALUES (?,?)
		`
		regionsSQLPrep, err := db.Prepare(regionsSQL)
		if err != nil {
			log.Printf("failed to prepare statement: %s\n", err.Error())
		}
		defer regionsSQLPrep.Close()

		for i := 0; i < len(result.Features); i++ {
			// Execute the statement with the charger record values
			_, err = chargersSQLPrep.Exec(
				result.Features[i].Properties.ClassName,
				result.Features[i].Properties.LastEdited,
				result.Features[i].Properties.Created,
				result.Features[i].Properties.SiteID,
				result.Features[i].Properties.Name,
				result.Features[i].Properties.Operator,
				result.Features[i].Properties.Address,
				result.Features[i].Properties.Is24Hours,
				result.Features[i].Properties.CarParkCount,
				result.Features[i].Properties.HasCarparkCost,
				result.Features[i].Properties.MaxTimeLimit,
				result.Features[i].Properties.HasTouristAttraction,
				result.Features[i].Properties.ProviderDeleted,
				result.Features[i].Properties.HideFromFeed,
				result.Features[i].Properties.RegionID,
				result.Features[i].Properties.ID,
				result.Features[i].Properties.RecordClassName,
				result.Features[i].Properties.ExternalID,
				result.Features[i].Properties.Uniq,
				result.Features[i].Properties.Type,
				result.Features[i].Properties.ID0,
				result.Features[i].Properties.LastUpdated,
				result.Features[i].Properties.Region,
				result.Features[i].Properties.OwnerName,
				result.Features[i].Properties.ChargingCost,
				result.Features[i].Properties.FeatureType,
			)
			if err != nil {
				log.Printf("failed to execute statement: %s\n", err.Error())
			}
			// Fetch the last inserted ID
			var lastInsertedID int64
			err = db.QueryRow("SELECT last_insert_rowid();").Scan(&lastInsertedID)
			if err != nil {
				log.Printf("failed to get last pk: %s\n", err.Error())
			}

			// cameras_access_locations
			for a := 0; a < len(result.Features[i].Properties.AccessLocations); a++ {

				// Execute the statement with the charger record values
				_, err = accessLocationsSQLPrep.Exec(
					lastInsertedID,
					result.Features[i].Properties.AccessLocations[a].Lat,
					result.Features[i].Properties.AccessLocations[a].Lon,
				)
				if err != nil {
					log.Printf("failed to execute statement: %s\n", err.Error())
				}
			}

			// cameras_connectors
			for a := 0; a < len(result.Features[i].Properties.Connectors); a++ {

				// Execute the statement with the charger record values
				_, err = connectorsSQLPrep.Exec(
					lastInsertedID,
					result.Features[i].Properties.Connectors[a].Current,
					result.Features[i].Properties.Connectors[a].KwRated,
					result.Features[i].Properties.Connectors[a].ConnectorType,
					result.Features[i].Properties.Connectors[a].OperationStatus,
					result.Features[i].Properties.Connectors[a].NextPlannedOutage,
				)
				if err != nil {
					log.Printf("failed to execute statement: %s\n", err.Error())
				}
			}

			// cameras_regions
			for a := 0; a < len(result.Features[i].Properties.Regions); a++ {

				// Execute the statement with the charger record values
				_, err = regionsSQLPrep.Exec(
					lastInsertedID,
					result.Features[i].Properties.Regions[a].ID,
				)
				if err != nil {
					log.Printf("failed to execute statement: %s\n", err.Error())
				}
			}
		}
		// remvoe duplicates from table
		_, err = db.Exec(`
		DELETE FROM chargers 
		WHERE id NOT IN (SELECT MIN(id) 
		FROM cameras 
		GROUP BY CONCAT(last_edited, created, uniq))
		`)
		if err != nil {
			log.Fatal("failed to remove duplicates from rss table: ", err)
			return
		}

		// remove orphaned
	}
}
