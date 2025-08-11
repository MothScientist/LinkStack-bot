package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
)

func setupLogs() *os.File {
	err := createDirectory()
	if err != nil {
		log.Fatalf("Error creating logs directory: %v;", err)
	}

	logFile, err := os.OpenFile("logs/bot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0222)
	if err != nil {
		log.Fatalf("Error creating/opening file .log: %v;", err)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return logFile
}

// createDirectory Creates a logs directory in the current working directory if it does not exist.
func createDirectory() error {
	// Get current working directory
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Forming the full path to the logs directory
	logsDir := filepath.Join(dir, "logs")

	// Check for directory existence
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		err := os.Mkdir(logsDir, 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		// Handling other Stat errors
		return err
	}
	return nil
}

func logPanic() {
	if r := recover(); r != nil {
            log.Printf(
                "PANIC: %v\nStack trace:\n%s;",
                r,
                string(debug.Stack()),
            )
            os.Exit(1)
        }
}