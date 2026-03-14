package main

import (
	"log/slog"
	"os"

	"github.com/AntoninHuaut/EtuEDT-Back/internal/config"
	"github.com/AntoninHuaut/EtuEDT-Back/internal/server"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	server.StartWebApp()
}
