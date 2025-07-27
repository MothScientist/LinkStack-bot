package main

import (
	"fmt"
	"log"
)

func main() {
	checkSqliteFileExists() // Checking for the presence of a DB file
	cfg := loadConfig() // Loading configuration
	logFile := setupLogs() // Setting up logging
	loadSqlQueries() // Loading SQL queries into memory

	defer logFile.Close()

	fmt.Print("Compilation was successful\n")
	log.Print("Launching the bot...")

	botProcess(cfg.Token)
}
