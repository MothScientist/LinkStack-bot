package main

import (
	"fmt"
	"log"
)

func main() {
	checkSqliteFileExists() // Checking for the presence of a DB file
	cfg := loadConfig() // Loading configuration
	loadSqlQueries() // Loading SQL queries into memory
	loadLocaleJson() // Loading translation memory for the /help command

	logFile := setupLogs() // Setting up logging (load last to use log.Fatal() in functions above)
	defer logFile.Close()
	defer logPanic()

	fmt.Print("Compilation was successful\n")
	log.Print("Launching the bot...")

	botProcess(cfg.Token)
}
