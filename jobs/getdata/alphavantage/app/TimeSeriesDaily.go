package app

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gocarina/gocsv"
	m "github.com/timotewb/cpu/jobs/getdata/alphavantage/models"
	"github.com/timotewb/cpu/jobs/getdata/common/config"
	"github.com/timotewb/cpu/jobs/getdata/common/helper"
)

func TimeSeriesDaily(allConfig config.AllConfig, jobConfig m.JobConfig){

	// get sqlite db
	db, dbPath, err := helper.GetOrCreateSQLiteDB(allConfig, "alphavantage")
	if err != nil {
		log.Fatalf("from TimeSeriesDaily(): function GetOrCreateSQLiteDB() failed: %v", err)
	}
	defer db.Close()

	// create target table if not exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS time_series_daily (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		symbol TEXT,
		timestamp TEXT,
		open REAL,
		high REAL,
		low REAL,
		close REAL,
		volume INTEGER
	)`)
	if err != nil {
		log.Fatalf("from TimeSeriesDaily(): failed to create table 'time_series_daily': %v\n", err)
	}

	var result []*m.TimeSeriesDaily
	var url string
	for _, s := range jobConfig.Symbols{
		url = fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&outputsize=full&apikey=%s&datatype=csv", s, jobConfig.APIKey)
		if csvBytes, err := helper.GetURLData(url); err != nil {
			log.Fatalf("from TimeSeriesDaily(): failed to get json: %v\n", err)
		} else {
			var reader io.Reader = bytes.NewReader(csvBytes)
			if err := gocsv.Unmarshal(reader, &result); err != nil {
				log.Fatalf("from TimeSeriesDaily(): unmarshal error: %v\n", err)
			}
			sqlStatement := `
			INSERT INTO time_series_daily (
				symbol,timestamp,open,high,low,close,volume
			) VALUES (?,?,?,?,?,?,?)
			`
			// Prepare the statement
			stmt, err := db.Prepare(sqlStatement)
			if err != nil {
				log.Fatalf("from TimeSeriesDaily(): failed to prepare insert statement: %s\n", err.Error())
			}
			defer stmt.Close()
			for i := 0; i < len(result); i++ {
				// Execute the statement with the camera record values
				_, err = stmt.Exec(
					s, result[i].Timestamp, result[i].Open, result[i].High, result[i].Low, result[i].Close, result[i].Volume,
				)
				if err != nil {
					log.Fatalf("from TimeSeriesDaily(): failed to execute insert statement: %s\n", err.Error())
				}
			}
			// remvoe duplicates from table
			_, err = db.Exec(`DELETE FROM time_series_daily WHERE id NOT IN (SELECT MIN(id) FROM time_series_daily GROUP BY CONCAT(symbol, timestamp, close))`)
			if err != nil {
				log.Fatalf("from TimeSeriesDaily(): failed to remove duplicates from cameras table: %v\n", err)
			}
			err = os.Chmod(dbPath, 0777)
			if err != nil {
				log.Fatal("from TimeSeriesDaily(): failed to set permissions on db file: ", err)
				return
			}
		}
	}
}