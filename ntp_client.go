package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	// NTP server to query (you can choose a different one if needed)
	ntpServer := "time.google.com"

	// Get the current time from the NTP server
	ntpTime, err := ntp.Time(ntpServer)
	if err != nil {
		log.Fatalf("Failed to get time from NTP server: %v", err)
	}

	// Get the current local system time
	systemTime := time.Now()

	// Calculate the time drift
	drift := ntpTime.Sub(systemTime)

	fmt.Printf("NTP server time: %v\n", ntpTime)
	fmt.Printf("System time:     %v\n", systemTime)
	fmt.Printf("Time drift:      %v\n", drift)

	// Optional: If the drift is greater than a certain threshold, print a warning
	threshold := time.Second * 1
	if drift > threshold || drift < -threshold {
		fmt.Println("Warning: Time drift exceeds threshold!")
	} else {
		fmt.Println("Time is in sync with NTP server.")
	}
}

