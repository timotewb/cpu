package app

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/timotewb/cpu/jobs/getdata/common/config"
	"github.com/timotewb/cpu/jobs/getdata/common/helper"
	m "github.com/timotewb/cpu/jobs/getdata/seek/models"
)

func GetData(allConfig config.AllConfig, jobConfig m.JobConfig){

	var result m.JobsListing
	for _, url := range jobConfig.URLs {
		for i := 0; i < jobConfig.PageCount; i++{
			if jsonBytes, err := helper.GetURLData(url + fmt.Sprintf("%d", i)); err != nil {
				log.Fatalf("from Cameras(): failed to get json: %v\n", err)
			} else {
				if err := json.Unmarshal(jsonBytes, &result); err != nil {
					log.Fatalf("from Cameras(): unmarshal error: %v\n", err)
				}
				fmt.Println(result)
			}
		}
	}
}