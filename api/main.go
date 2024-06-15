package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
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

	var configDir string
	var help bool
	// Define CLI flags in short and long form
	flag.StringVar(&configDir, "c", "", "Path where configuration file is stored (shorthand)")
	flag.StringVar(&configDir, "config", "", "Path where configuration file is stored")
	flag.BoolVar(&help, "h", false, "Show usage instructions (shorthand)")
	flag.BoolVar(&help, "help", false, "Show usage instructions")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "----------------------------------------------------------------------------------------")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Pass -c to specify where the configuration file is stored:")
		fmt.Fprintln(os.Stderr, "  -c\t\tstring\n  --config")
		fmt.Fprintln(os.Stderr, "  \tPath where configuration file is stored")
		fmt.Fprintln(os.Stderr, "\n  -h\t\tboolean\n  --help")
		fmt.Fprintln(os.Stderr, "  \tShow usage instructions")
		fmt.Fprintln(os.Stderr, "----------------------------------------------------------------------------------------")
	}
	flag.Parse()

	// Print the Help docuemntation to the terminal if user passes help flag
	if help {
		flag.Usage()
		return
	}

	router := gin.Default()

	// POST /api handler for executing jobs.
	// It reads the job name from the request body, checks if it's in the configuration's job list, and executes the job if found.
	router.POST("/api", func(c *gin.Context) {

		// Read config file each time api is called
		config, err := app.ReadConfig(configDir)
		if err != nil {
			log.Fatal(err)
		}

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
			cmdArgs := []string{filepath.Join(config.AppPath, body.Name)}

			// Check if args is present and not empty
			if len(body.Args) > 0 {
				// Append each argument to the command slice
				cmdArgs = append(cmdArgs, body.Args...)
			}

			// Execute the command
			cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error executing command: %s", err), "cmd": cmd.String()})
				return
			}
			// Write the output back to the response
			c.String(http.StatusOK, string(output))
		} else {
			// Write the output back to the response
			c.JSON(http.StatusNotFound, gin.H{"error": "Name not found in job list", "job": body.Name})
		}

	})

	router.Run(":3000")
}
