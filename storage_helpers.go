package main

import (
	"log"
	"os"
    "path/filepath"
)

var queries = &SqlQueries{}

func loadSqlQueries() {
    queryLoad := []struct {
        queriesField *string
        sqlFilename string
    }{
        {&queries.addRecord, "add_record"},
        {&queries.delRecord, "del_record"},
        {&queries.getRecord, "get_record"},
        {&queries.getListRecords, "get_list_records"},
        {&queries.getRandomRecord, "get_random_record"},
        {&queries.recordIsExists, "record_is_exists"},
    }

    for _, load := range queryLoad {
		var err error
		*load.queriesField, err = readSqlFile(load.sqlFilename)
        if err != nil {
            log.Fatalf("Error reading file %s: %w", load.sqlFilename, err)
        }
    }
}

func readSqlFile(filename string) (string, error) {
	sqlFile, err := os.ReadFile("sql/" + filename + ".sql")
	if err != nil {
		return "", err
	}
	return string(sqlFile), nil
}

func checkSqliteFileExists() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %w", err)
	}

	dbPath := filepath.Join(currentDir, "main.sqlite3")

	fileInfo, err := os.Stat(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("File main.sqlite3 does not exist")
		}
		log.Fatalf("File verification error: %w", err)
	}

	if fileInfo.IsDir() {
		log.Fatal("main.sqlite3 is a directory, not a file")
	}

    if fileInfo.Mode().Perm()&0400 == 0 {
        log.Fatal("No permission to read file main.sqlite3")
    }
}