package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Print("Checking existence of data/ in the project folder...")

	err := ensureChildDirsExist()
	
	if err != nil {
		log.Fatalf("An error has occured while creating necessary data dirs: %v", err)
	}
	
	err = ensureEssentialFilesExist()
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("Configuring server...")
	localPort := ":8080"
	configureServer()

	err = runServer(localPort)
	if err != nil {
		log.Printf("Server error: %v", err)
	}
}

