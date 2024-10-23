package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/timotewb/cpu/jobs/getdata/common/config"
	"github.com/timotewb/cpu/jobs/getdata/common/helper"
	m "github.com/timotewb/cpu/jobs/getdata/seek/models"
)

func GetJobListings(allConfig config.AllConfig, jobConfig m.JobConfig){

	// get sqlite db
	db, dbPath, err := helper.GetOrCreateSQLiteDB(allConfig, "seek")
	if err != nil {
		log.Fatalf("from TimeSeriesDaily(): function GetOrCreateSQLiteDB() failed: %v", err)
	}
	defer db.Close()

	//----------------------------------------------------------------------------------------
	// job_listings
	//----------------------------------------------------------------------------------------
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
		company_psd_id INTEGER,
		location_name TEXT,
		location_id INTEGER,
		location_wv TEXT,
		job_location TEXT,
		job_location_country TEXT,
		listing_datetime TEXT,
		role_id TEXT,
		salary TEXT,
		sub_classification_id INTEGER,
		sub_classification_desc TEXT,
		suburb TEXT,
		suburb_id INTEGER,
		suburb_wv TEXT,
		teaser TEXT,
		title TEXT,
		work_type TEXT
	)`)
	if err != nil {
		log.Fatalf("from GetJobListings(): failed to create table 'job_listings': %v\n", err)
	}
	sqlStatement := `
	INSERT INTO job_listings (
		job_listing_id, advertiser_id, advertiser_desc, area_id, area_desc, area_where_desc, classification_id, classification_desc, company_name, company_psd_id,
		location_name, location_id, location_wv, job_location, job_location_country, listing_datetime, role_id, salary, sub_classification_id, sub_classification_desc,
		suburb, suburb_id, suburb_wv, teaser, title, work_type
	) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`
	// Prepare the statement
	stmtJL, err := db.Prepare(sqlStatement)
	if err != nil {
		log.Fatalf("from GetJobListings(): failed to prepare insert stmtJL statement: %s\n", err.Error())
	}
	defer stmtJL.Close()

	//----------------------------------------------------------------------------------------
	// bullet_points
	//----------------------------------------------------------------------------------------
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
		log.Fatalf("from TimeSeriesDaily(): failed to prepare insert stmtBP statement: %s\n", err.Error())
	}
	defer stmtBP.Close()

	//----------------------------------------------------------------------------------------
	// work_arrangements
	//----------------------------------------------------------------------------------------
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS work_arrangements (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		job_listing_id INTEGER,
		arrangement TEXT
	)`)
	if err != nil {
		log.Fatalf("from GetJobListings(): failed to create table 'work_arrangements': %v\n", err)
	}
	sqlStatement = `
	INSERT INTO work_arrangements (
		job_listing_id, arrangement
	) VALUES (?,?)
	`
	// Prepare the statement
	stmtWA, err := db.Prepare(sqlStatement)
	if err != nil {
		log.Fatalf("from TimeSeriesDaily(): failed to prepare insert stmtWA statement: %s\n", err.Error())
	}
	defer stmtWA.Close()

	//----------------------------------------------------------------------------------------
	// start
	//----------------------------------------------------------------------------------------
	var result m.JobsListing
	for _, url := range jobConfig.URLs {
		for i := 0; i < jobConfig.PageCount; i++{
			if jsonBytes, err := helper.GetURLData(url.URL + fmt.Sprintf("%d", i)); err != nil {
				log.Fatalf("from GetJobListings(): failed to get json: %v\n", err)
			} else {
				time.Sleep(5 * time.Second)
				if err := json.Unmarshal(jsonBytes, &result); err != nil {
					log.Fatalf("from GetJobListings(): unmarshal error: %v\n", err)
				}

				for j := 0; j < len(result.Data); j++ {
					_, err = stmtJL.Exec(
						result.Data[j].ID, result.Data[j].Advertiser.ID, result.Data[j].Advertiser.Description, result.Data[j].AreaID, result.Data[j].Area, 
						result.Data[j].AreaWhereValue, result.Data[j].Classification.ID, result.Data[j].Classification.Description, result.Data[j].CompanyName, 
						result.Data[j].CompanyProfileStructuredDataID, result.Data[j].Location, result.Data[j].LocationID, result.Data[j].LocationWhereValue,
						result.Data[j].JobLocation.Label, result.Data[j].JobLocation.CountryCode, result.Data[j].ListingDate, result.Data[j].RoleID,
						result.Data[j].Salary, result.Data[j].SubClassification.ID, result.Data[j].SubClassification.Description, result.Data[j].Suburb, result.Data[j].SuburbID,
						result.Data[j].SuburbWhereValue, result.Data[j].Teaser, result.Data[j].Title, result.Data[j].WorkType,
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
					for w := 0; w < len(result.Data[j].WorkArrangements.Data); w++ {
						_, err = stmtWA.Exec(
							result.Data[j].ID, result.Data[j].WorkArrangements.Data[w].Label.Text,
						)
						if err != nil {
							log.Fatalf("from GetJobListings(): failed to execute insert stmtWA statement: %s\n", err.Error())
						}

					}
				}	
			}
		}
	}

	//----------------------------------------------------------------------------------------
	// tidy up
	//----------------------------------------------------------------------------------------
	// remvoe duplicates from table
	_, err = db.Exec(`DELETE FROM job_listings WHERE id NOT IN (SELECT MIN(id) FROM job_listings GROUP BY CONCAT(job_listing_id, listing_datetime))`)
	if err != nil {
		log.Fatalf("from GetJobListings(): failed to remove duplicates from job_listings table: %v\n", err)
	}
	_, err = db.Exec(`DELETE FROM bullet_points WHERE id NOT IN (SELECT MIN(id) FROM bullet_points GROUP BY CONCAT(job_listing_id, point))`)
	if err != nil {
		log.Fatalf("from GetJobListings(): failed to remove duplicates from bullet_points table: %v\n", err)
	}
	_, err = db.Exec(`DELETE FROM work_arrangements WHERE id NOT IN (SELECT MIN(id) FROM work_arrangements GROUP BY CONCAT(job_listing_id, arrangement))`)
	if err != nil {
		log.Fatalf("from GetJobListings(): failed to remove duplicates from work_arrangements table: %v\n", err)
	}
	err = os.Chmod(dbPath, 0777)
	if err != nil {
		log.Fatal("from GetJobListings(): failed to set permissions on db file: ", err)
		return
	}
}