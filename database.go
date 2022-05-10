package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var s string = "HHHHHHHHH"

func loadDatabase(filename string) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db1, err1 := sql.Open("sqlite3", ":memory:")
	if err1 != nil {
		log.Fatalf("cannot open an SQLite memory database: %v", err)
	}
	defer db1.Close()

	_, err = db.Exec("CREATE TABLE unix_time (time datetime); INSERT INTO unix_time (time) VALUES (strftime('%Y-%m-%dT%H:%MZ','now'))")
	if err != nil {
		log.Fatalf("cannot create schema: %v", err)
	}

}

func saveDatabase(filename string) {

}
