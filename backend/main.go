package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Print("Checking existence of data/ in the project folder...")

	dirErr := checkDataDirExist()
	
	if dirErr != nil {
		log.Fatalf("An error has occured while creating the data dir: %v", dirErr)
	}
	
	fmt.Println("Configuring server...")
	localPort := ":8080"
	configureServer()

	err := runServer(localPort)
	if err != nil {
		log.Printf("Server error: %v", err)
	}
}

