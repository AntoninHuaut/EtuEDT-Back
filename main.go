package main

import (
	"EtuEDT-Go/api"
	"EtuEDT-Go/cache"
	"EtuEDT-Go/config"
	"log"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	err := cache.StartScheduler()
	if err != nil {
		log.Fatal(err)
	}

	api.StartWebApp()
}
