package main

import (
	"fmt"
)

func main(){
	fmt.Print("Checking existence of data/ in the project folder...")

	dirExist := checkDataDirExist()
	
	if dirExist != nil {
		fmt.Printf("An error has occured while creating the data dir: %v", dirExist)
	}
	
	fmt.Println("Configuring server...")
	localPort := ":8080"
	configureServer()

	err := runServer(localPort)
	if err != nil {
		fmt.Errorf("Server error: %v", err)
	}
}

