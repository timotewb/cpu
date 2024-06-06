package app

import (
	"encoding/json"
	"log"

	"github.com/timotewb/cpu/jobs/getdata/common/config"
	"github.com/timotewb/cpu/jobs/getdata/common/helper"
)

func Cameras(allConfig config.AllConfig, jobConfig JobConfig) {

	// get sqlite db
	db, _, err := helper.GetOrCreateSQLiteDB(allConfig, "journeys_nzta")
	if err != nil {
		log.Fatalf("from Cameras(): function GetOrCreateSQLiteDB() failed: %v", err)
	}
	defer db.Close()

	// create target table if not exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cameras (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		class_name TEXT,
		last_edited TEXT,
		created TEXT,
		external_id INTEGER,
		name TEXT,
		description TEXT,
		offline INTEGER,
		under_maintenance INTEGER,
		image_url TEXT,
		latitude TEXT,
		longitude TEXT,
		direction TEXT,
		sort_group TEXT,
		tas_journey_id INTEGER,
		region_id INTEGER,
		tas_region_id INTEGER,
		property_id INTEGER,
		uniq TEXT,
		property_type TEXT,
		last_updated INTEGER
	)`)
	if err != nil {
		log.Fatalf("from Cameras(): failed to create table 'cameras': %v\n", err)
	}

	var result CamerasModel
	if jsonBytes, err := helper.GetURLData(jobConfig.CamerasURL); err != nil {
		log.Fatalf("from Cameras(): failed to get json: %v\n", err)
	} else {
		if err := json.Unmarshal(jsonBytes, &result); err != nil {
			log.Fatalf("from Cameras(): unmarshal error: %v\n", err)
		}
		sqlStatement := `
		INSERT INTO cameras (
			class_name, last_edited, created, external_id, name, description, offline, under_maintenance, image_url, latitude, longitude, direction, sort_group, tas_journey_id, region_id, tas_region_id, property_id, uniq, property_type, last_updated
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
		`
		// Prepare the statement
		stmt, err := db.Prepare(sqlStatement)
		if err != nil {
			log.Fatalf("from Cameras(): failed to prepare insert statement: %s\n", err.Error())
		}
		defer stmt.Close()
		for i := 0; i < len(result.Features); i++ {
			// Execute the statement with the camera record values
			_, err = stmt.Exec(
				result.Features[i].Properties.ClassName, result.Features[i].Properties.LastEdited, result.Features[i].Properties.Created, result.Features[i].Properties.ExternalID, result.Features[i].Properties.Name, result.Features[i].Properties.Description, result.Features[i].Properties.Offline, result.Features[i].Properties.UnderMaintenance, result.Features[i].Properties.ImageURL, result.Features[i].Properties.Latitude, result.Features[i].Properties.Longitude, result.Features[i].Properties.Direction, result.Features[i].Properties.SortGroup, result.Features[i].Properties.TasJourneyID, result.Features[i].Properties.RegionID, result.Features[i].Properties.TasRegionID, result.Features[i].Properties.ID, result.Features[i].Properties.Uniq, result.Features[i].Properties.Type, result.Features[i].Properties.LastUpdated,
			)
			if err != nil {
				log.Fatalf("from Cameras(): failed to execute insert statement: %s\n", err.Error())
			}
		}
		// remvoe duplicates from table
		_, err = db.Exec(`DELETE FROM cameras WHERE id NOT IN (SELECT MIN(id) FROM cameras GROUP BY CONCAT(last_edited, created, uniq))`)
		if err != nil {
			log.Fatalf("from Cameras(): failed to remove duplicates from cameras table: %v\n", err)
		}
	}
}
