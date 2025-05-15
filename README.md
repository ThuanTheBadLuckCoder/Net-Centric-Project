# Net-Centric-Project

# Final Demo Structure 
// Project Structure for Tower Castle Rush (TCR) Game

/*

tcr-game/

├── main.go                    // Entry point

├── go.mod                     // Go module file

├── go.sum                     // Go dependencies

├── README.md                  // Project documentation

├── .gitignore                 // Git ignore file

├── config/

│   ├── config.go              // Configuration management

│   └── game_config.json       // Game configuration settings

├── data/

│   ├── troops.json            // Troop specifications

│   ├── towers.json            // Tower specifications

│   └── players/               // Player data directory

│       ├── player1.json       // Individual player data

│       └── player2.json       // Individual player data

├── internal/

│   ├── auth/

│   │   ├── auth.go            // Authentication logic

│   │   └── user.go            // User management

│   ├── game/

│   │   ├── game.go            // Core game logic

│   │   ├── engine.go          // Game engine (Enhanced mode)

│   │   ├── simple_rules.go    // Simple TCR implementation

│   │   ├── enhanced_rules.go  // Enhanced TCR implementation

│   │   ├── battle.go          // Battle mechanics

│   │   └── events.go          // Game events

│   ├── models/

│   │   ├── player.go          // Player model

│   │   ├── tower.go           // Tower model

│   │   ├── troop.go           // Troop model

│   │   ├── battle.go          // Battle result model

│   │   └── gamestate.go       // Game state model

│   ├── server/

│   │   ├── server.go          // HTTP/WebSocket server

│   │   ├── handlers.go        // Request handlers

│   │   ├── websocket.go       // WebSocket connection management

│   │   └── middleware.go      // Server middleware

│   ├── storage/

│   │   ├── json_storage.go    // JSON file storage

│   │   ├── player_storage.go  // Player data persistence

│   │   └── game_storage.go    // Game data persistence

│   └── utils/

│       ├── calculator.go      // Damage calculations

│       ├── validator.go       // Input validation

│       └── logger.go          // Logging utilities

├── cmd/

│   ├── server/

│   │   └── main.go            // Server command

│   └── client/

│       └── main.go            // Test client (optional)

├── pkg/

│   ├── protocol/

│   │   ├── messages.go        // Game protocol messages

│   │   └── constants.go       // Game constants

│   └── errors/

│       └── errors.go          // Custom error types

├── web/

│   ├── static/

│   │   ├── css/

│   │   └── styles.css             

│   │   ├── js/

│   │   └── scripts.js       

│   │   └── images/

│   └── templates/

│       ├── index.html         // Game interface

│       └── game.html          // Game board
├── tests/

│   ├── unit/

│   │   ├── game_test.go       // Game logic tests

│   │   ├── battle_test.go     // Battle mechanics tests

│   │   └── models_test.go     // Model tests

│   ├── integration/

│   │   ├── server_test.go     // Server integration tests

│   │   └── game_flow_test.go  // End-to-end game tests

│   └── testdata/

│       ├── test_troops.json   // Test troop data

│       └── test_players.json  // Test player data

└── scripts/

    ├── build.sh               // Build script

    ├── test.sh                // Test script

    └── setup.sh               // Initial setup script

*/


# How to Run?
Windows
1. Find Project Folder: Right click to tcr-game then choose Git Base Here

2. Input: ./scripts/setup.sh
