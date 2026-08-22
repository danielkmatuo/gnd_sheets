package main

import (
	"log"

	"github.com/danielkmatuo/gnd_sheets/backend"
)

func main() {
	err := backend.Run()
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}

