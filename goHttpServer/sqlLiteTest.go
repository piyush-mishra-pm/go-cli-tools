package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDb() {

	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		log.Fatalf("Error opening DB: %s", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Error Pinging DB: %s", err)
	}

	fmt.Println("Connected Sql Lite DB")
}

func closeDb() {
	fmt.Println("Closed DB")
	db.Close()
}

func getQueryResult() string {
	defer db.Close()

	var version string
	err := db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)

	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("Version: %s", version)
}
