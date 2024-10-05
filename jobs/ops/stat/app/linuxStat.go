package app

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/timotewb/cpu/jobs/ops/stat/models"
)

func LinuxStat() models.ServerModel {
	var returnData models.ServerModel

	// load average
	cmd := exec.Command("cat", "/proc/loadavg")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("from cmd.CombinedOutput() error loadavg: %v", err)
		return returnData
	}
	parts := strings.Fields(string(output))
	returnData.LoadAverage = fmt.Sprintf("%s %s %s", parts[0], parts[1], parts[2])
	returnData.RunningProcs = parts[3]

	// uptime
	cmd = exec.Command("uptime", "-p")
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Printf("from cmd.CombinedOutput() error uptime: %v", err)
		return returnData
	}
	returnData.UpTime = strings.ReplaceAll(string(output), "\n", "")

	// hostname
	cmd = exec.Command("hostname")
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Printf("from cmd.CombinedOutput() error uptime: %v", err)
		return returnData
	}
	returnData.Name = strings.ReplaceAll(string(output), "\n", "")

	return returnData
} 