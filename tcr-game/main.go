// Fixed main.go - Remove unused import
package main

import (
	"flag"
	"log"

	"tcr-game/internal/server"
	"tcr-game/config"
)

func main() {
	var configPath = flag.String("config", "config/game_config.json", "Path to configuration file")
	var port = flag.String("port", "8080", "Server port")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create and start server
	srv := server.New(cfg)
	log.Printf("Starting TCR Game Server on port %s", *port)
	
	if err := srv.Start(*port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}