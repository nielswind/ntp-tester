package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"os"
	"reflect"
	"strconv"

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
	fmt.Println("== Starting ntp-checker with prometheus metrics ==")
	metricsPort := getEnv("METRICS_PORT","2112")
	ntpServer := getEnv("NTP_SERVER","time.google.com") // cph-dc-2.corp.local
	checkEvery, _ := time.ParseDuration(getEnv("CHECK_DURATION","10s"))
	allowedDrift, _ := strconv.ParseFloat(getEnv("MAX_DRIFT_ALLOWED_SECONDS", "2"),64)
	fmt.Printf("Using ntp server %s and prometheus port %s\n", ntpServer, metricsPort)
	fmt.Println("Updates metrics every",checkEvery,"with allowed drift of",allowedDrift,"seconds")

	// Start a HTTP server for Prometheus to scrape
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":"+metricsPort, nil)) // Expose metrics on port 2112
	}()

	// Run the drift checker periodically
	for {
		checkTimeDrift(ntpServer, allowedDrift)
		time.Sleep(checkEvery) // Check drift every 60 seconds
	}
}

func checkTimeDrift(ntpServer string, allowedDrift float64) {
	// Get the current time from the NTP server
	// ntpTime, err := ntp.Time(ntpServer)
	ntpTime, err := ntp.Time(ntpServer)
	if err != nil {
		log.Printf("Failed to get time from NTP server %s: %v", ntpServer, err)
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
	threshold := allowedDrift;
	if drift > threshold || drift < -threshold {
		fmt.Println(os.Stderr, "Warning: Time drift exceeds threshold!")
	} else {
		fmt.Println("Time is in sync with NTP server.")
	}
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}
