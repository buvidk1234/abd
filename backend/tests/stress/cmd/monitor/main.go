package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	interval = flag.Int("i", 5, "Monitoring interval (seconds)")
	duration = flag.Int("d", 60, "Monitoring duration (seconds)")
	logFile  = flag.String("log", "system_monitor.log", "Log file path")
	targets  = flag.String("targets", "abd-server,mysql,redis,kafka", "Process names to monitor (comma separated)")
)

func main() {
	flag.Parse()

	f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, f))

	log.Printf("Starting system monitor: interval %ds, duration %ds, targets: %s", *interval, *duration, *targets)
	
	targetList := strings.Split(*targets, ",")
	startTime := time.Now()
	endTime := startTime.Add(time.Duration(*duration) * time.Second)

	for time.Now().Before(endTime) {
		log.Printf("--- %s ---", time.Now().Format("15:04:05"))
		for _, target := range targetList {
			stats := getProcessStats(target)
			log.Printf("%-15s: %s", target, stats)
		}
		time.Sleep(time.Duration(*interval) * time.Second)
	}

	log.Printf("System monitor finished.")
}

func getProcessStats(name string) string {
	// Simple Windows tasklist/powershell based stats
	// This is a rough estimation for the stress test report
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Get-Process -Name %s* | Select-Object -Property Name, CPU, WorkingSet | ConvertTo-Json", name))
	out, err := cmd.Output()
	if err != nil {
		return "Not Running"
	}
	// We just return raw or simplified json for logging, can be parsed later
	return string(out)
}
