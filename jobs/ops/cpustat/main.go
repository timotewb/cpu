package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/timotewb/cpu/jobs/ops/cpustat/app"
	"github.com/timotewb/cpu/jobs/ops/cpustat/models"
)

func main(){

	// get path of code
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("from os.Executable():", err)
		return
	}
	fullPath := filepath.Dir(exePath)

	// read config
	conf, err := app.ReadJobConfig(fullPath)
	if err != nil {
		log.Fatalf("from app.ReadJobConfig(): %s", err)
	}

	// get stats
	var resp models.BlobType
	// last update
	now := time.Now()
    resp.LastUpdated = now.Format("2006-01-02 15:04:05")

	// server details
	h, err := host.Info()
	if err != nil {
		log.Fatalf("from host.Info(): %s", err)
	}
	resp.Platform = fmt.Sprintf("%v, %v, %v", h.Platform, h.PlatformFamily, h.PlatformVersion)
	resp.UpTime = app.Uptime(h.Uptime)
    resp.Name = strings.ToLower(h.Hostname)
	resp.RunningProcs = int64(h.Procs)

	// resources
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalf("from mem.VirtualMemory(): %s", err)
	}
	resp.MemoryUsedPct = v.UsedPercent

	c, err := cpu.Info()
	if err != nil {
		log.Fatalf("from cpu.Info(): %s", err)
	}
	resp.CPUCoreCount = 0
	for i := range c {
		resp.CPUCoreCount = resp.CPUCoreCount + int64(c[i].Cores)
		resp.CPUModel = c[i].ModelName
	}
	
	l, err := load.Avg()
	if err != nil {
		log.Fatalf("from load.Avg(): %s", err)
	}
	resp.LoadAverage = fmt.Sprintf("%v, %v, %v", l.Load1, l.Load5, l.Load15)

	// prepare for output
	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatalf("from json.MarshalIndent(): %v", err)
	}

	// write to blob
	blobName := fmt.Sprintf("%v-latest.json", resp.Name)
	app.EnvVariables(fullPath)
	app.WriteToBlob(conf.StorageAccountName, conf.ContainerName, blobName, data)

}