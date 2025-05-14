# Tower Castle Rush (TCR) Game

A real-time multiplayer tower defense game implemented in Go with WebSocket support.

## Features

### Game Modes

1. **Simple TCR Rules**
   - Turn-based gameplay
   - Each player gets 3 random troops
   - Must destroy guard towers before king tower
   - Winner destroys opponent's king tower

2. **Enhanced TCR Rules**
   - Real-time gameplay (3 minutes)
   - Mana system (starts at 5, regenerates 1/sec, max 10)
   - Critical hit system
   - Experience and leveling system
   - Winner has most towers destroyed or destroys king tower first

## Architecture

- **Models**: Game entities (Player, Troop, Tower, Game)
- **Auth**: Authentication and user management
- **Game Engine**: Core game logic with mode-specific managers
- **Storage**: JSON-based persistence
- **Server**: HTTP/WebSocket server with real-time updates
- **Protocol**: Standardized message formats

## Quick Start

```bash
# Setup project
./scripts/setup.sh

# Build
./scripts/build.sh

# Run tests
./scripts/test.sh

# Start server
go run main.go
```

## API Endpoints

- `POST /api/register` - Create new account
- `POST /api/login` - Login
- `POST /api/games` - Create game
- `POST /api/games/{id}/join` - Join game
- `GET /api/games/{id}/state` - Get game state
- `POST /api/games/{id}/action` - Make game action
- `WS /ws/{id}` - WebSocket connection

## Configuration

Edit `config/game_config.json` to customize:
- Server settings
- Game duration and rules
- Experience points
- Mana system parameters

## File Structure

```
tcr-game/
├── main.go
├── internal/         # Core application code
├── config/          # Configuration files
├── data/            # Game data (troops, towers, players)
├── web/             # Frontend assets
├── scripts/         # Build and utility scripts
└── tests/           # Test files
```

## Game Rules

### Simple Mode
1. Two players take turns
2. Each player has 3 random troops
3. Must destroy first guard tower before second guard tower
4. Must destroy both guard towers before king tower
5. First to destroy king tower wins

### Enhanced Mode
1. Real-time gameplay for 3 minutes
2. Mana system limits troop spawning
3. Can attack any living tower
4. Critical hits deal 120% damage
5. King tower destruction or most towers destroyed wins

## Development

### Prerequisites
- Go 1.21+
- Git

### Running Locally
```bash
git clone <repository>
cd tcr-game
go mod tidy
go run main.go
```

### Contributing
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request