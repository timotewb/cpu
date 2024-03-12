package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/timotewb/cpu/api/app"
)

// Body represents the JSON structure of the request body for the API endpoint.
type Body struct {
	// Name is the name of the job to be executed.
	Name string   `json:"name"`
	Args []string `json:"args"`
}

// main initializes the application, reads the configuration, and starts the HTTP server.
// It listens for POST requests on "/api" and executes a job if the job name is found in the configuration.
func main() {

	// Read config file
	config, err := app.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	// POST /api handler for executing jobs.
	// It reads the job name from the request body, checks if it's in the configuration's job list, and executes the job if found.
	router.POST("/api", func(c *gin.Context) {

		// Read the contact from the request body
		var body Body
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check to see if body.Name is in config.JobList
		found := false
		for _, job := range config.JobList {
			if job == body.Name {
				found = true
				break
			}
		}
		if found {

			// Start building the command slice
			cmdArgs := []string{"-c", filepath.Join(config.AppPath, body.Name)}

			// Check if args is present and not empty
			if len(body.Args) > 0 {
				// Append each argument to the command slice
				cmdArgs = append(cmdArgs, body.Args...)
			}

			// Execute the command
			cmd := exec.Command("sh", cmdArgs...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error executing command: %s", err)})
				return
			}
			// Write the output back to the response
			c.String(http.StatusOK, string(output))
		} else {
			// Write the output back to the response
			c.JSON(http.StatusNotFound, gin.H{"error": "Name not found in job list"})
		}

	})

	router.Run(":3000")
}
