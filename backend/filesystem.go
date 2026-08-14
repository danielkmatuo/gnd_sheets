package main

import (
 "fmt"
 "os"
 "path/filepath"
)

const (
	characters = "../data/reference/characters.json"
)

func buildDataDir() error {
	fmt.Print("Creating new data dir...")

	return os.MkdirAll("../data", 0755)
}

func buildCharactersDataDir() error {
	fmt.Print("Creating new characters data dir...")
	return os.MkdirAll("../data/characters", 0755)
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

func ensureChildDirsExist() error {
    err := ensureDataDirExist()
    if err != nil {
        return err
    }

    info, err := os.Stat("../data/reference")

    if err == nil {
        if !info.IsDir() {
            return fmt.Errorf("reference is not a directory")
        }
    } else if os.IsNotExist(err) {
        return fmt.Errorf("reference dir does not exist; therefore no reference file exists")
    } else {
        return err
    }

    info, err = os.Stat("../data/characters")

    if err == nil {
        if !info.IsDir() {
            return fmt.Errorf("characters exists, but is not a directory")
        }

        fmt.Print("data/characters/ dir already exists... Skipping step...\n")
        return nil
    }

    if os.IsNotExist(err) {
        fmt.Print("data/characters/ dir doesn't exist, creating it...\n")

        err = buildCharactersDataDir()
        if err != nil {
            return err
        }

        fmt.Print("data/characters/ dir created successfully!\n")
        return nil
    }

    return err
}
