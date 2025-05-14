// Fixed internal/game/engine.go - Remove unused fmt import
package game

import (
	"errors"
	"sync"
	
	"tcr-game/config"
	"tcr-game/internal/models"
	"tcr-game/internal/storage"
)

type GameEngine struct {
	storage         *storage.JSONStorage
	simpleManager   *SimpleGameManager
	enhancedManager *EnhancedGameManager
	activeGames     map[string]*models.Game
	mutex           sync.RWMutex
	config          *config.Config
}

func NewGameEngine(cfg *config.Config, storage *storage.JSONStorage) *GameEngine {
	simpleManager := NewSimpleGameManager(
		cfg.Game.Simple.MaxPlayers,
		cfg.Game.Simple.TurnTime,
		cfg.Game.Enhanced.CritMultiplier,
	)
	
	enhancedManager := NewEnhancedGameManager(
		cfg.Game.Enhanced.ManaRegen,
		cfg.Game.Enhanced.GameDuration,
		cfg.Game.Enhanced.ExpWin,
		cfg.Game.Enhanced.ExpDraw,
		cfg.Game.Enhanced.CritMultiplier,
	)
	
	return &GameEngine{
		storage:         storage,
		simpleManager:   simpleManager,
		enhancedManager: enhancedManager,
		activeGames:     make(map[string]*models.Game),
		config:          cfg,
	}
}

func (ge *GameEngine) CreateGame(gameID string, mode models.GameMode) (*models.Game, error) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()
	
	if _, exists := ge.activeGames[gameID]; exists {
		return nil, errors.New("game already exists")
	}
	
	game := models.NewGame(gameID, mode)
	ge.activeGames[gameID] = game
	
	return game, nil
}

func (ge *GameEngine) JoinGame(gameID string, player *models.Player) error {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()
	
	game, exists := ge.activeGames[gameID]
	if !exists {
		return errors.New("game not found")
	}
	
	if !game.AddPlayer(player) {
		return errors.New("game is full")
	}
	
	// Load available troops for the player
	if err := ge.loadPlayerTroops(player); err != nil {
		return err
	}
	
	// Start game if we have enough players
	if len(game.Players) == 2 {
		return ge.startGame(game)
	}
	
	return nil
}

func (ge *GameEngine) startGame(game *models.Game) error {
	switch game.Mode {
	case models.SimpleMode:
		return ge.simpleManager.StartGame(game)
	case models.EnhancedMode:
		return ge.enhancedManager.StartGame(game)
	default:
		return errors.New("invalid game mode")
	}
}

func (ge *GameEngine) loadPlayerTroops(player *models.Player) error {
	troops, err := ge.storage.LoadTroops()
	if err != nil {
		return err
	}
	
	player.AvailableTroops = make([]*models.Troop, len(troops))
	for i, troop := range troops {
		playerTroop := *troop
		playerTroop.ApplyLevel(player.GetTroopLevel(troop.ID))
		player.AvailableTroops[i] = &playerTroop
	}
	
	return nil
}

func (ge *GameEngine) ProcessSimpleAction(gameID, playerID string, action TurnAction) (*TurnResult, error) {
	ge.mutex.RLock()
	game, exists := ge.activeGames[gameID]
	ge.mutex.RUnlock()
	
	if !exists {
		return nil, errors.New("game not found")
	}
	
	if game.Mode != models.SimpleMode {
		return nil, errors.New("game is not in simple mode")
	}
	
	return ge.simpleManager.ProcessTurn(game, playerID, action)
}

func (ge *GameEngine) ProcessEnhancedAction(gameID, playerID string, action EnhancedAction) (*EnhancedResult, error) {
	ge.mutex.RLock()
	game, exists := ge.activeGames[gameID]
	ge.mutex.RUnlock()
	
	if !exists {
		return nil, errors.New("game not found")
	}
	
	if game.Mode != models.EnhancedMode {
		return nil, errors.New("game is not in enhanced mode")
	}
	
	return ge.enhancedManager.ProcessAction(gameID, playerID, action)
}

func (ge *GameEngine) GetGame(gameID string) (*models.Game, error) {
	ge.mutex.RLock()
	defer ge.mutex.RUnlock()
	
	game, exists := ge.activeGames[gameID]
	if !exists {
		return nil, errors.New("game not found")
	}
	
	return game, nil
}

func (ge *GameEngine) GetGameState(gameID string) (map[string]interface{}, error) {
	game, err := ge.GetGame(gameID)
	if err != nil {
		return nil, err
	}
	
	switch game.Mode {
	case models.SimpleMode:
		return ge.simpleManager.GetGameState(game), nil
	case models.EnhancedMode:
		return ge.enhancedManager.GetGameState(gameID)
	default:
		return nil, errors.New("invalid game mode")
	}
}

func (ge *GameEngine) EndGame(gameID string, reason string) error {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()
	
	game, exists := ge.activeGames[gameID]
	if !exists {
		return errors.New("game not found")
	}
	
	switch game.Mode {
	case models.SimpleMode:
		return ge.simpleManager.EndGame(game, reason)
	case models.EnhancedMode:
		ge.enhancedManager.CleanupGame(gameID)
		return nil
	default:
		return errors.New("invalid game mode")
	}
}

func (ge *GameEngine) CleanupGame(gameID string) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()
	
	if game, exists := ge.activeGames[gameID]; exists {
		if game.Mode == models.EnhancedMode {
			ge.enhancedManager.CleanupGame(gameID)
		}
		delete(ge.activeGames, gameID)
	}
}

func (ge *GameEngine) GetActiveGames() map[string]*models.Game {
	ge.mutex.RLock()
	defer ge.mutex.RUnlock()
	
	games := make(map[string]*models.Game)
	for id, game := range ge.activeGames {
		games[id] = game
	}
	return games
}