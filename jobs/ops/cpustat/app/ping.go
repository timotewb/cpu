package app

import (
	"log"
	"os/exec"
	"strings"
)

func Ping(ipAddr string) bool {
	if ipAddr == "" {
		log.Fatalf("No IP Address provided.")
	}
	out, _ := exec.Command("ping", ipAddr, "-c 5", "-i 3").Output()
	if strings.Contains(string(out), "Destination Host Unreachable") || strings.Contains(string(out), "No route to host") || strings.Contains(string(out), "Request timeout for icmp_seq") {
		return false
	} else {
		return true
	}
}