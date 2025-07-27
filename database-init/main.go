package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func main() {
	databaseCreate()
	createTable()
}

func databaseCreate() {
	file, err := os.Create("../main.sqlite3")
	if err != nil {
		log.Fatalf("Error creating main.sqlite3: %v", err)
	}
	err = file.Close()
	if err != nil {
		log.Fatalf("Error closing main.sqlite3: %v", err)
	}

	fmt.Printf("File main.sqlite3 created")
}

func createTable() {
	db, err := sql.Open("sqlite3", "../main.sqlite3")
	if err != nil {
		log.Fatal(err)
	}

	sqlFile, err := os.ReadFile("./init.sql")
	if err != nil {
		log.Fatalf("Error reading file init.sql: %v", err)
	}

	sqlStatements := string(sqlFile)
	if _, err := db.Exec(sqlStatements); err != nil {
		log.Fatalf("SQL execution error: %v", err)
	}

	err = db.Close()
	if err != nil {
		log.Fatalf("Error closing main.sqlite3: %v", err)
	}

	fmt.Println("The database has been initialized successfully.")
}
