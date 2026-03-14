package main

import (
	"log"

	"github.com/AntoninHuaut/EtuEDT-Back/api"
	"github.com/AntoninHuaut/EtuEDT-Back/domain"
)

func main() {
	if err := domain.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	api.StartWebApp()
}
