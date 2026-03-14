package main

import (
	"log/slog"
	"os"

	"github.com/AntoninHuaut/EtuEDT-Back/domain"
	"github.com/AntoninHuaut/EtuEDT-Back/server"
)

func main() {
	if err := domain.LoadConfig(); err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	server.StartWebApp()
}
