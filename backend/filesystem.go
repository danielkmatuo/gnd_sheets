package main

import (
 "fmt"
 "os"
)

func buildDataDir() error {
	fmt.Print("Creating new data dir...")

	return os.MkdirAll("../data", 0755)
}

func checkDataDirExist() error {
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
