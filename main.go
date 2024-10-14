package main

import (
	"EtuEDT-Go/api"
	"EtuEDT-Go/cache"
	"EtuEDT-Go/domain"
	"log"
)

func main() {
	if err := domain.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	err := cache.StartScheduler()
	if err != nil {
		log.Fatal(err)
	}

	api.StartWebApp()
}
