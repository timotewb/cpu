package helper

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/timotewb/cpu/jobs/getdata/common/config"
	"github.com/timotewb/cpu/jobs/getdata/common/models"
)

func CallCPUJob(configDir, name string, args []string) string{

	// Read All Config
	allConfig, err := config.ReadAllConfig(configDir)
	if err != nil {
		log.Fatalf("function ReadAllConfig() failed: %v", err)
	}

	var body models.CallCPUJobType
	body.Name = name
	body.Args = append(body.Args, args...)

	// Prepare the JSON body
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("error marshalling body to JSON: %v", err)
	}

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
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading response body: %v", err)
	}
	return string(bytes)

}