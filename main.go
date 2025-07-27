package main

import (
	"fmt"
	"log"
)

func main() {
	if !checkSqliteExists() {
		log.Fatal("The file .sqlite3 is missing")
	}

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	loadSqlQueries()

	fmt.Print("Compilation was successful\n")

	botProcess(cfg.Token)
}
