package main

import (
	"github.com/NicholasRodrigues/mini-db/internal/config"
	"github.com/NicholasRodrigues/mini-db/internal/server"
)

func main() {
	config.LoadConfig()
	s := server.NewServer()
	s.Start()
}
