package main

import (
	"archive/zip"
	"fmt"
	"golang.org/x/exp/slices"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func listZipFiles(file *zip.File) error {
	fileread, err := file.Open()
	if err != nil {
		msg := "Failed to open zip %s for reading: %s"
		return fmt.Errorf(msg, file.Name, err)
	}
	defer fileread.Close()

	fmt.Fprintf(os.Stdout, "%s:", file.Name)

	if err != nil {
		msg := "Failed to read zip %s for reading: %s"
		return fmt.Errorf(msg, file.Name, err)
	}

	fmt.Println()

	return nil
}

func collectFiles(path string) []string {
	var fileList []string
	var extensions = []string{".zip", ".fb2"}

	files, _ := os.ReadDir(path)
	for _, f := range files {
		fileName := path + string(os.PathSeparator) + f.Name()
		if f.IsDir() == true {
			fileList = append(fileList, collectFiles(fileName)...)
		} else {
			if slices.Contains(extensions, filepath.Ext(fileName)) {
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

						//fileList = append(fileList, fileName+"@"+file.Name+strconv.FormatUint(file.UncompressedSize64, 10))

						//fmt.Fprintf(os.Stdout, "%s:", file.Name)

						if err != nil {
							//msg := "Failed to read zip %s for reading: %s"
							//return fmt.Errorf(msg, file.Name, err)
						}

						fmt.Println()
					}

				} else {
					fi, err := os.Stat(fileName)
					if err != nil {
						//return err
					}

					fileList = append(fileList, fileName+"@"+strconv.FormatInt(fi.Size(), 10)+"@"+fi.ModTime().Format("2006-01-02 15:04:05"))
				}
			}
		}
	}
	return fileList
}

func main() {

	fl := collectFiles("/home/kiltum/Calibre Library")
	for _, file := range fl {
		fmt.Println(file)
	}

}
