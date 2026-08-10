package main

import (
	"fmt"
)

func main(){
	fmt.Print("Checking existence of data/ in the project folder...")

	err := checkDataDirExist()
	
	fmt.Print(err)
}

