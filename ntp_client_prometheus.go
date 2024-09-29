package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/beevik/ntp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define a Prometheus gauge metric to track the time drift
var (
	timeDriftGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ntp_time_drift_seconds",
		Help: "Time drift between local system and NTP server in seconds",
	})
)

func init() {
	// Register the gauge metric with Prometheus
	prometheus.MustRegister(timeDriftGauge)
}

func main() {
	// Start a HTTP server for Prometheus to scrape
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":2112", nil)) // Expose metrics on port 2112
	}()

	// NTP server to query (you can choose a different one if needed)
	ntpServer := "time.google.com"

	// Run the drift checker periodically
	for {
		checkTimeDrift(ntpServer)
		time.Sleep(60 * time.Second) // Check drift every 60 seconds
	}
}

func checkTimeDrift(ntpServer string) {
	// Get the current time from the NTP server
	ntpTime, err := ntp.Time(ntpServer)
	if err != nil {
		log.Printf("Failed to get time from NTP server: %v", err)
		return
	}

	// Get the current local system time
	systemTime := time.Now()

	// Calculate the time drift in seconds
	drift := ntpTime.Sub(systemTime).Seconds()

	// Update the Prometheus gauge metric with the drift value
	timeDriftGauge.Set(drift)

	fmt.Printf("NTP server time: %v\n", ntpTime)
	fmt.Printf("System time:     %v\n", systemTime)
	fmt.Printf("Time drift:      %.6f seconds\n", drift)

	// Optional: Print a warning if the drift exceeds a certain threshold
	threshold := 1.0 // 1 second threshold
	if drift > threshold || drift < -threshold {
		fmt.Println("Warning: Time drift exceeds threshold!")
	} else {
		fmt.Println("Time is in sync with NTP server.")
	}
}

