package app

import (
	"encoding/json"
	"log"
	"os"

	"github.com/timotewb/cpu/jobs/getdata/common/config"
	"github.com/timotewb/cpu/jobs/getdata/common/helper"
)

func Chargers(allConfig config.AllConfig, jobConfig JobConfig) {

	// get sqlite db
	db, dbPath, err := helper.GetOrCreateSQLiteDB(allConfig, "journeys_nzta")
	if err != nil {
		log.Fatalf("from Chargers(): function GetOrCreateSQLiteDB() failed: %v\n", err)
	}
	defer db.Close()

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
		feature_type TEXT
	)`)
	if err != nil {
		log.Fatalf("from Chargers(): failed to create table 'chargers': %v\n", err)
	}
	// chargers_access_locations
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS chargers_access_locations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chargers_id INTEGER,
		lat REAL,
		lon REAL,
		FOREIGN KEY(chargers_id) REFERENCES chargers(id) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Fatalf("from Chargers(): failed to create table 'chargers_access_locations': %v\n", err)
	}
	// chargers_connectors
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS chargers_connectors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chargers_id INTEGER,
		current TEXT,
		kw_rated INTEGER,
		connector_type TEXT,
		operation_status TEXT,
		next_planning_outage TEXT,
		FOREIGN KEY(chargers_id) REFERENCES chargers(id) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Fatalf("from Chargers(): failed to create table 'chargers_connectors': %v\n", err)
	}
	// chargers_regions
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS chargers_regions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chargers_id INTEGER,
		regions_id INTEGER,
		FOREIGN KEY(chargers_id) REFERENCES chargers(id) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Fatalf("from Chargers(): failed to create table 'chargers_regions': %v\n", err)
	}

	var result ChargersModel
	if jsonBytes, err := helper.GetURLData(jobConfig.ChargersURL); err != nil {
		log.Fatalf("from Chargers(): failed to get json: %v\n", err)
	} else {
		if err := json.Unmarshal(jsonBytes, &result); err != nil {
			log.Fatalf("from Chargers(): unmarshal error: %v\n", err)
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
			feature_type 
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
		`
		chargersSQLPrep, err := db.Prepare(chargersSQL)
		if err != nil {
			log.Fatalf("from Chargers(): failed to prepare 'INSERT INTO chargers' statement: %s\n", err.Error())
		}
		defer chargersSQLPrep.Close()

		// regionssSQL
		regionsSQL := `
		INSERT INTO chargers_regions (
			chargers_id ,
			regions_id
		) VALUES (?,?)
		`
		regionsSQLPrep, err := db.Prepare(regionsSQL)
		if err != nil {
			log.Fatalf("from Chargers(): failed to prepare 'INSERT INTO chargers_regions' statement: %s\n", err.Error())
		}
		defer regionsSQLPrep.Close()

		// accessLocationsSQL
		accessLocationsSQL := `
		INSERT INTO chargers_access_locations (
			chargers_id ,
			lat ,
			lon 
		) VALUES (?,?,?)
		`
		accessLocationsSQLPrep, err := db.Prepare(accessLocationsSQL)
		if err != nil {
			log.Fatalf("from Chargers(): failed to prepare 'INSERT INTO chargers_access_locations' statement: %s\n", err.Error())
		}
		defer accessLocationsSQLPrep.Close()

		// connectorsSQL
		connectorsSQL := `
		INSERT INTO chargers_connectors (
			chargers_id ,
			current ,
			kw_rated ,
			connector_type ,
			operation_status ,
			next_planning_outage 
		) VALUES (?,?,?,?,?,?)
		`
		connectorsSQLPrep, err := db.Prepare(connectorsSQL)
		if err != nil {
			log.Fatalf("from Chargers(): failed to prepare 'INSERT INTO chargers_connectors' statement: %s\n", err.Error())
		}
		defer connectorsSQLPrep.Close()

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
				log.Fatalf("from Chargers(): failed to execute 'INSERT INTO chargers' statement: %s\n", err.Error())
			}
			// Fetch the last inserted ID
			var lastInsertedID int64
			err = db.QueryRow("SELECT last_insert_rowid();").Scan(&lastInsertedID)
			if err != nil {
				log.Fatalf("from Chargers(): failed to get last pk: %s\n", err.Error())
			}

			// cameras_regions
			for a := 0; a < len(result.Features[i].Properties.Regions); a++ {

				// Execute the statement with the charger record values
				_, err = regionsSQLPrep.Exec(
					lastInsertedID,
					result.Features[i].Properties.Regions[a].ID,
				)
				if err != nil {
					log.Fatalf("from Chargers(): failed to execute 'INSERT INTO chargers_regions' statement: %s\n", err.Error())
				}
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
					log.Fatalf("from Chargers(): failed to execute 'INSERT INTO chargers_access_locations' statement: %s\n", err.Error())
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
					log.Fatalf("from Chargers(): failed to execute 'INSERT INTO chargers_connectors' statement: %s\n", err.Error())
				}
			}
		}
		// remvoe duplicates from table
		_, err = db.Exec(`DELETE FROM chargers WHERE id NOT IN (SELECT MIN(id) FROM chargers GROUP BY CONCAT(last_edited, created, uniq))`)
		if err != nil {
			log.Fatalf("from Chargers(): failed to remove duplicates from rss table: %v\n", err)
		}
		err = os.Chmod(dbPath, 0777)
		if err != nil {
			log.Fatal("from Chargers(): failed to set permissions on db file: ", err)
			return
		}
	}
}
