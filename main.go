package main

import (
	"EtuEDT-Go/api"
	"EtuEDT-Go/domain"
	"log"
)

func main() {
	if err := domain.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	api.StartWebApp()
}
