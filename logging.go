package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func logEndedSession(meetingID string, duration time.Duration) {
	logFile, err := os.OpenFile("logs/bbb-multicleaner-logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)
	logger.Printf("ENDED: %s | Duration: %s | Hours: %.2f", meetingID, duration.String(), duration.Hours())
}

func logError(errorType, details string) {
	errorFile, err := os.OpenFile("logs/bbb-multicleaner-errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening error log file: %v\n", err)
		return
	}
	defer errorFile.Close()

	logger := log.New(errorFile, "", log.LstdFlags)
	logger.Printf("ERROR: %s | Details: %s", errorType, details)
}

func logDebug(debugType, details string) {
	debugFile, err := os.OpenFile("logs/bbb-multicleaner-logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening debug log file: %v\n", err)
		return
	}
	defer debugFile.Close()
	logger := log.New(debugFile, "", log.LstdFlags)
	logger.Printf("DEBUG: %s | Details: %s", debugType, details)
}
