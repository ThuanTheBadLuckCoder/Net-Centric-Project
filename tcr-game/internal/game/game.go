// Fixed internal/game/game.go - Remove unused fmt import
package game

import (
	"errors"
	"sync"
	"time"
	
	"tcr-game/internal/models"
)

type GameManager interface {
	StartGame(game *models.Game) error
	GetGameState(gameID string) (map[string]interface{}, error)
	EndGame(gameID string, reason string) error
}

type GameController struct {
	games       map[string]*models.Game
	simpleGM    *SimpleGameManager
	enhancedGM  *EnhancedGameManager
	eventManager *EventManager
	mutex       sync.RWMutex
}

func NewGameController(simpleGM *SimpleGameManager, enhancedGM *EnhancedGameManager) *GameController {
	return &GameController{
		games:       make(map[string]*models.Game),
		simpleGM:    simpleGM,
		enhancedGM:  enhancedGM,
		eventManager: NewEventManager(),
	}
}

func (gc *GameController) CreateGame(gameID string, mode models.GameMode) (*models.Game, error) {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	
	if _, exists := gc.games[gameID]; exists {
		return nil, errors.New("game already exists")
	}
	
	game := models.NewGame(gameID, mode)
	gc.games[gameID] = game
	
	return game, nil
}

func (gc *GameController) GetGame(gameID string) (*models.Game, error) {
	gc.mutex.RLock()
	defer gc.mutex.RUnlock()
	
	game, exists := gc.games[gameID]
	if !exists {
		return nil, errors.New("game not found")
	}
	
	return game, nil
}

func (gc *GameController) JoinGame(gameID string, player *models.Player) error {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	
	game, exists := gc.games[gameID]
	if !exists {
		return errors.New("game not found")
	}
	
	// Check if player already in game
	for _, p := range game.Players {
		if p.ID == player.ID {
			return errors.New("player already in game")
		}
	}
	
	if !game.AddPlayer(player) {
		return errors.New("game is full")
	}
	
	// Publish player joined event
	gc.eventManager.PublishPlayerJoined(gameID, player)
	
	// Start game if we have enough players
	if len(game.Players) == 2 && game.State == models.Waiting {
		return gc.startGame(game)
	}
	
	return nil
}

func (gc *GameController) startGame(game *models.Game) error {
	var err error
	
	switch game.Mode {
	case models.SimpleMode:
		err = gc.simpleGM.StartGame(game)
	case models.EnhancedMode:
		err = gc.enhancedGM.StartGame(game)
	default:
		err = errors.New("invalid game mode")
	}
	
	if err == nil {
		gc.eventManager.PublishGameStarted(game.ID, game)
	}
	
	return err
}

func (gc *GameController) ProcessSimpleAction(gameID, playerID string, action TurnAction) (*TurnResult, error) {
	game, err := gc.GetGame(gameID)
	if err != nil {
		return nil, err
	}
	
	if game.Mode != models.SimpleMode {
		return nil, errors.New("game is not in simple mode")
	}
	
	result, err := gc.simpleGM.ProcessTurn(game, playerID, action)
	if err != nil {
		return nil, err
	}
	
	// Publish attack event if successful
	if result.Success && result.BattleResult != nil {
		gc.eventManager.PublishAttackMade(gameID, result.BattleResult)
		
		if result.BattleResult.TowerDestroyed {
			gc.eventManager.PublishTowerDestroyed(
				gameID, 
				result.BattleResult.TargetTower, 
				result.BattleResult.DefenderID,
			)
		}
		
		if result.BattleResult.GameEnded {
			gc.eventManager.PublishGameEnded(gameID, result.BattleResult.Winner, "king_tower_destroyed")
		}
	}
	
	return result, nil
}

func (gc *GameController) ProcessEnhancedAction(gameID, playerID string, action EnhancedAction) (*EnhancedResult, error) {
	game, err := gc.GetGame(gameID)
	if err != nil {
		return nil, err
	}
	
	if game.Mode != models.EnhancedMode {
		return nil, errors.New("game is not in enhanced mode")
	}
	
	result, err := gc.enhancedGM.ProcessAction(gameID, playerID, action)
	if err != nil {
		return nil, err
	}
	
	// Publish attack event if successful
	if result.Success && result.BattleResult != nil {
		gc.eventManager.PublishAttackMade(gameID, result.BattleResult)
		
		if result.BattleResult.TowerDestroyed {
			gc.eventManager.PublishTowerDestroyed(
				gameID, 
				result.BattleResult.TargetTower, 
				result.BattleResult.DefenderID,
			)
		}
		
		if result.BattleResult.GameEnded {
			gc.eventManager.PublishGameEnded(gameID, result.BattleResult.Winner, "king_tower_destroyed")
		}
	}
	
	return result, nil
}

func (gc *GameController) GetGameState(gameID string) (map[string]interface{}, error) {
	game, err := gc.GetGame(gameID)
	if err != nil {
		return nil, err
	}
	
	switch game.Mode {
	case models.SimpleMode:
		return gc.simpleGM.GetGameState(game), nil
	case models.EnhancedMode:
		return gc.enhancedGM.GetGameState(gameID)
	default:
		return nil, errors.New("invalid game mode")
	}
}

func (gc *GameController) EndGame(gameID string, reason string) error {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	
	game, exists := gc.games[gameID]
	if !exists {
		return errors.New("game not found")
	}
	
	if game.State == models.Finished {
		return errors.New("game already finished")
	}
	
	// End the game based on mode
	var err error
	switch game.Mode {
	case models.SimpleMode:
		err = gc.simpleGM.EndGame(game, reason)
	case models.EnhancedMode:
		gc.enhancedGM.CleanupGame(gameID)
		game.State = models.Finished
		endTime := time.Now()
		game.EndTime = &endTime
	}
	
	if err == nil {
		// Determine winner if not set
		winner := ""
		if game.Winner != nil {
			winner = game.Winner.ID
		}
		
		gc.eventManager.PublishGameEnded(gameID, winner, reason)
		gc.eventManager.CleanupGame(gameID)
	}
	
	return err
}

func (gc *GameController) CleanupGame(gameID string) {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	
	if game, exists := gc.games[gameID]; exists {
		if game.Mode == models.EnhancedMode {
			gc.enhancedGM.CleanupGame(gameID)
		}
		gc.eventManager.CleanupGame(gameID)
		delete(gc.games, gameID)
	}
}

func (gc *GameController) GetActiveGames() map[string]*models.Game {
	gc.mutex.RLock()
	defer gc.mutex.RUnlock()
	
	games := make(map[string]*models.Game)
	for id, game := range gc.games {
		games[id] = game
	}
	return games
}

func (gc *GameController) GetGameCount() int {
	gc.mutex.RLock()
	defer gc.mutex.RUnlock()
	return len(gc.games)
}

func (gc *GameController) Subscribe(gameID string, ch chan GameEventData) error {
	if _, err := gc.GetGame(gameID); err != nil {
		return err
	}
	
	gc.eventManager.Subscribe(gameID, ch)
	return nil
}

func (gc *GameController) Unsubscribe(gameID string, ch chan GameEventData) {
	gc.eventManager.Unsubscribe(gameID, ch)
}

// GameStats provides statistics about ongoing games
type GameStats struct {
	TotalGames      int                    `json:"total_games"`
	ActiveGames     int                    `json:"active_games"`
	GamesByMode     map[string]int         `json:"games_by_mode"`
	GamesByState    map[string]int         `json:"games_by_state"`
	AverageGameTime time.Duration          `json:"average_game_time"`
}

func (gc *GameController) GetStats() GameStats {
	gc.mutex.RLock()
	defer gc.mutex.RUnlock()
	
	stats := GameStats{
		TotalGames:  len(gc.games),
		ActiveGames: 0,
		GamesByMode: make(map[string]int),
		GamesByState: make(map[string]int),
	}
	
	totalDuration := time.Duration(0)
	finishedGames := 0
	
	for _, game := range gc.games {
		// Count by mode
		stats.GamesByMode[string(game.Mode)]++
		
		// Count by state
		stats.GamesByState[string(game.State)]++
		
		// Count active games
		if game.State == models.InProgress {
			stats.ActiveGames++
		}
		
		// Calculate average game time for finished games
		if game.State == models.Finished && game.EndTime != nil {
			duration := game.EndTime.Sub(game.StartTime)
			totalDuration += duration
			finishedGames++
		}
	}
	
	if finishedGames > 0 {
		stats.AverageGameTime = totalDuration / time.Duration(finishedGames)
	}
	
	return stats
}