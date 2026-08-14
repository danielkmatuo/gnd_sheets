package main

import (
 "fmt"
 "os"
 "path/filepath"
)

const (
	characters = "../data/characters.json"
)

func buildDataDir() error {
	fmt.Print("Creating new data dir...")

	return os.MkdirAll("../data", 0755)
}

func buildLocalDataDir() error {
	fmt.Print("Creating new local data dir...")

	return os.MkdirAll("../data/local", 0755)
}

func ensureDataDirExist() error {
	info, err := os.Stat("../data")

	if err == nil {
		if !info.IsDir(){
			return fmt.Errorf("data dir exists, but its not a directory...\n")
		}

		fmt.Print("data/ dir already exists... Skipping step...\n")
		return nil
	}
	
	if os.IsNotExist(err) {
		fmt.Print("data dir doesn't exist, creating it...\n")

		err = buildDataDir()
		if err != nil {
			return err
		}

		fmt.Print("data dir created succesfully!\n")
		return nil
	}
	return err
}

func ensureEssentialFilesExist() error {
	filesList, err := os.ReadDir("../data")
	if err != nil {
		return err
	}

	files := []string {characters}
	var foundFiles []string

	for _, file := range filesList {
		filePath := filepath.Join("../data", file.Name())

		for _, srcFile := range files {
			if filePath == srcFile {
				foundFiles = append(foundFiles, filePath)
				break
			} 		
		}
	}
	
	if len(files) != len(foundFiles) {
		return fmt.Errorf("necessary files not found")
	}

	return nil
}

func ensureChildDirExist() error {
	err := ensureDataDirExist()
	if err != nil {
		return err
	}

	info, err := os.Stat("../data/local")

	if err == nil {
		if !info.IsDir(){
			return fmt.Errorf("local child dir exists, but its not a directory...\n")
		}

		fmt.Print("data/local/ dir already exists... Skipping step...\n")
		return nil
	}
	
	if os.IsNotExist(err) {
		fmt.Print("data/local/ dir doesn't exist, creating it...\n")

		err = buildLocalDataDir()
		if err != nil {
			return err
		}

		fmt.Print("data/local/ dir created succesfully!\n")
		return nil
	}
	return err
}
