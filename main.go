package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	pathToZip := flag.String("dirs", "", "output directory")
	server := flag.String("server", "", "server url")
	user := flag.String("user", "", "user id")
	password := flag.String("password", "", "password or token")

	flag.Parse()

	if *pathToZip == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println(*pathToZip)

	if !fileExists(*pathToZip) {
		log.Fatalf("%s is not exist", *pathToZip)
	}

	writer, err := createDataZip("data.zip")

	if isDir(*pathToZip) {
		err = recursiveZip(*pathToZip, writer)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = closeDataZip(writer)
	if err != nil {
		log.Fatal(err)
	}

	err = sendFile(*server, "data.zip", *user, *password)
	if err != nil {
		log.Fatal(err)
	}

}

func createDataZip(name string) (*zip.Writer, error) {
	destinationFile, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	zipWriter := zip.NewWriter(destinationFile)
	return zipWriter, nil
}

func closeDataZip(file *zip.Writer) error {
	err := file.Close()
	if err != nil {
		return err
	}

	return nil
}

func addFileToZip(name string, writer *zip.Writer) {

}

func recursiveZip(pathToZip string, writer *zip.Writer) error {
	err := filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		zipFile, err := writer.Create(filePath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func sendFile(server, fileName, user, password string) error {
	dataFile, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer dataFile.Close()
	payload := io.MultiReader(dataFile)
	req, err := http.NewRequest("POST", server, payload)
	if err != nil {
		return err
	}
	req.SetBasicAuth(user, password)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func isDir(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
