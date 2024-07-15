package main

import (
	"github.com/NicholasRodrigues/mini-db/internal/config"
	"github.com/NicholasRodrigues/mini-db/internal/server"
	"log"
)

func main() {
	config.LoadConfig()
	s := server.NewServer()

	go func() {
		if err := server.StartMetricsServer(); err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	s.Start()
}
