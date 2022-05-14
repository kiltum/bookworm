package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"time"
)

// check  s contains e
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

var collectedBefore uint64
var collected uint64

func printCollected() {
	for range time.Tick(time.Second * 1) {
		fmt.Printf("Collected %d/%d files      \r", collected, collectedBefore)
	}
}

func collectFiles(path string) {
	var extensions = []string{".zip", ".fb2"} // what we know to handle now

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	stmtInsertFiles, err := db.PrepareContext(ctx, "INSERT INTO files(path,size,modtime) VALUES (?,?,?)")
	if err != nil {
		log.Fatal(err)
	}

	stmtCheckFiles, err := db.PrepareContext(ctx, "SELECT size, modtime FROM files WHERE PATH = ?")
	if err != nil {
		log.Fatal(err)
	}

	db.QueryRow("SELECT count(*) FROM files").Scan(&collectedBefore)

	go printCollected()

	files, _ := os.ReadDir(path)
	for _, f := range files {
		fileName := path + string(os.PathSeparator) + f.Name()
		if f.IsDir() == true {
			collectFiles(fileName)
		} else {
			if contains(extensions, filepath.Ext(fileName)) {
				fi, err := os.Stat(fileName)
				if err != nil {
					//return err
				}

				var fs int64
				var md string
				err = stmtCheckFiles.QueryRow(fileName).Scan(&fs, &md)
				if err != nil {
					if err == sql.ErrNoRows {
						// insert as new file, that need to be parsed
						stmtInsertFiles.Exec(fileName, fi.Size(), fi.ModTime().Format(time.RFC3339))
						//fmt.Println("NONE  " + fileName)
					}
				} else {
					if fi.Size() != fs {
						// Change to re-parse
					}
					if md != fi.ModTime().Format(time.RFC3339) {
						// Change to re-parse
					}
					// Ok, file is already in database with same time and size, so do nothing
				}
				collected = collected + 1
			}
		}
	}
	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	viper.SetConfigName("bookworm")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/home/kiltum/projects/bookworm/")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/bookworm/")
	viper.AddConfigPath("$HOME/.bookworm")

	viper.SetEnvPrefix("bw")
	viper.BindEnv("sqldriver")
	viper.BindEnv("path")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("config not found")
			viper.Set("sqldriver", "sqlite3")
			viper.Set("path", []string{"."})
			viper.SafeWriteConfig()
		} else {
			fmt.Println(viper.ConfigFileUsed(), err)
			os.Exit(1)
		}
	}

	fmt.Println(viper.Get("sqldriver"))
	fmt.Println(viper.Get("path"))

	loadDatabase("db.sql")
	defer closeDatabase()

	collectFiles("/home/kiltum/Calibre Library")
	collectFiles("/mnt/swamp/Torrent")
	//for _, file := range fl {
	//	fmt.Println(file)
	//}

}
