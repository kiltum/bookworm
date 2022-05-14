package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var db *sql.DB

var sl3_create string = `CREATE TABLE IF NOT EXISTS files
                         (file_id INTEGER PRIMARY KEY,
                         path TEXT NOT NULL,
						 size INTEGER,
                         modtime TEXT NOT NULL);
`

func loadDatabase(filename string) {
	var err error
	db, err = sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(sl3_create)
	if err != nil {
		log.Fatalf("cannot create schema: %v", err)
	}
	fmt.Println("DB opened")

}

func closeDatabase() {
	db.Close() // TODO: need ? check for closing error
	fmt.Println("DB closed")
}
