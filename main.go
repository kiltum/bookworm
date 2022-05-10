package main

import (
	"archive/zip"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strconv"
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

func collectFiles(path string) []string {
	var fileList []string
	var extensions = []string{".zip", ".fb2"} // what we know to handle now

	files, _ := os.ReadDir(path)
	for _, f := range files {
		fileName := path + string(os.PathSeparator) + f.Name()
		if f.IsDir() == true {
			fileList = append(fileList, collectFiles(fileName)...)
		} else {
			if contains(extensions, filepath.Ext(fileName)) {
				if filepath.Ext(fileName) == ".zip" {

					read, err := zip.OpenReader(fileName)
					if err != nil {
						msg := "Failed to open: %s"
						log.Fatalf(msg, err)
					}
					defer read.Close()

					for _, file := range read.File {
						fileread, err := file.Open()
						if err != nil {
							//msg := "Failed to open zip %s for reading: %s"
							//return fmt.Errorf(msg, file.Name, err)
						}
						defer fileread.Close()

						fileList = append(fileList, fileName+"|"+file.Name+"|"+strconv.FormatUint(file.UncompressedSize64, 10)+"|"+file.ModTime().Format(time.RFC3339))

						//fmt.Fprintf(os.Stdout, "%s:", file.Name)

						if err != nil {
							//msg := "Failed to read zip %s for reading: %s"
							//return fmt.Errorf(msg, file.Name, err)
						}

						//fmt.Println()
					}

				} else { // no, its usual uncompressed file.fb2
					fi, err := os.Stat(fileName)
					if err != nil {
						//return err
					}

					fileList = append(fileList, fileName+"|"+strconv.FormatInt(fi.Size(), 10)+"|"+fi.ModTime().Format(time.RFC3339))
				}
			}
		}
	}
	return fileList
}

func main() {

	//viper.Set("sqldriver", "sqlite3")

	viper.SetConfigName("bookworm")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/home/kiltum/projects/bookworm/")
	//viper.AddConfigPath(".")
	//viper.AddConfigPath("/etc/bookworm/")
	//viper.AddConfigPath("$HOME/.bookworm")

	viper.SetEnvPrefix("bw")
	viper.BindEnv("sqldriver")

	//viper.WriteConfig()
	//viper.SafeWriteConfig()
	//viper.WriteConfigAs("/home/kiltum/projects/bookworm/")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("config not found")
			// Config file not found; ignore error if desired
		} else {
			fmt.Println("error in config")
			// Config file was found but another error was produced
		}
	}

	fmt.Println(viper.Get("sqldriver"))

	loadDatabase("db.sql")
	fmt.Println(s)
	//fl := collectFiles("/home/kiltum/Calibre Library")
	//for _, file := range fl {
	//	fmt.Println(file)
	//}

}
