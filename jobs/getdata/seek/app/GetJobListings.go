package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/timotewb/cpu/jobs/getdata/common/config"
	"github.com/timotewb/cpu/jobs/getdata/common/helper"
	m "github.com/timotewb/cpu/jobs/getdata/seek/models"
)

func GetJobListings(allConfig config.AllConfig, jobConfig m.JobConfig){

	// get sqlite db
	db, dbPath, err := helper.GetOrCreateSQLiteDB(allConfig, "alphavantage")
	if err != nil {
		log.Fatalf("from TimeSeriesDaily(): function GetOrCreateSQLiteDB() failed: %v", err)
	}
	defer db.Close()

	// job_listings
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS job_listings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		job_listing_id INTEGER,
		advertiser_id INTEGER,
		advertiser_desc TEXT,
		area_id INTEGER,
		area_desc TEXT,
		area_where_desc TEXT,
		classification_id INTEGER,
		classification_desc TEXT,
		company_name TEXT,
		company_psd_id INTEGER
	)`)
	if err != nil {
		log.Fatalf("from GetJobListings(): failed to create table 'job_listings': %v\n", err)
	}
	sqlStatement := `
	INSERT INTO job_listings (
		job_listing_id, advertiser_id, advertiser_desc, area_id, area_desc, area_where_desc, classification_id, classification_desc, company_name, company_psd_id
	) VALUES (?,?,?,?,?,?,?,?,?,?)
	`
	// Prepare the statement
	stmtJL, err := db.Prepare(sqlStatement)
	if err != nil {
		log.Fatalf("from TimeSeriesDaily(): failed to prepare insert statement: %s\n", err.Error())
	}
	defer stmtJL.Close()

	// bullet_points
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS bullet_points (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		job_listing_id INTEGER,
		point TEXT
	)`)
	if err != nil {
		log.Fatalf("from GetJobListings(): failed to create table 'bullet_points': %v\n", err)
	}
	sqlStatement = `
	INSERT INTO bullet_points (
		job_listing_id, point
	) VALUES (?,?)
	`
	// Prepare the statement
	stmtBP, err := db.Prepare(sqlStatement)
	if err != nil {
		log.Fatalf("from TimeSeriesDaily(): failed to prepare insert statement: %s\n", err.Error())
	}
	defer stmtBP.Close()
	
	var result m.JobsListing
	for _, url := range jobConfig.URLs {
		for i := 0; i < jobConfig.PageCount; i++{
			if jsonBytes, err := helper.GetURLData(url + fmt.Sprintf("%d", i)); err != nil {
				log.Fatalf("from GetJobListings(): failed to get json: %v\n", err)
			} else {
				if err := json.Unmarshal(jsonBytes, &result); err != nil {
					log.Fatalf("from GetJobListings(): unmarshal error: %v\n", err)
				}

				for j := 0; j < len(result.Data); j++ {
					_, err = stmtJL.Exec(
						result.Data[j].ID, result.Data[j].Advertiser.ID, result.Data[j].Advertiser.Description, result.Data[j].AreaID, result.Data[j].Area, 
						result.Data[j].AreaWhereValue, result.Data[j].Classification.ID, result.Data[j].Classification.Description, result.Data[j].CompanyName, 
						result.Data[j].CompanyProfileStructuredDataID,
					)
					if err != nil {
						log.Fatalf("from GetJobListings(): failed to execute insert stmtJL statement: %s\n", err.Error())
					}
					for b := 0; b < len(result.Data[j].BulletPoints); b++ {
						_, err = stmtBP.Exec(
							result.Data[j].ID, result.Data[j].BulletPoints[b],
						)
						if err != nil {
							log.Fatalf("from GetJobListings(): failed to execute insert stmtBP statement: %s\n", err.Error())
						}
					}
				}	
			}
		}
	}

	// udpate permissions
	err = os.Chmod(dbPath, 0777)
	if err != nil {
		log.Fatal("from GetJobListings(): failed to set permissions on db file: ", err)
		return
	}
}