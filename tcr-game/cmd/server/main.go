// cmd/server/main.go - Server command
package main

import (
	"log"
	
	"tcr-game/config"
	"tcr-game/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config/game_config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create and start server
	srv := server.New(cfg)
	log.Printf("Starting TCR Game Server on port %s", cfg.Server.Port)
	
	if err := srv.Start(cfg.Server.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}