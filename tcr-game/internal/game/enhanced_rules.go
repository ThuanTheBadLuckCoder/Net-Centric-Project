// Fixed internal/game/enhanced_rules.go - Properly declare types
package game

import (
	"errors"
	"fmt"
	"sync"
	"time"
	
	"tcr-game/internal/models"
)

type EnhancedGameManager struct {
	battleEngine    *BattleEngine
	manaRegenRate   float64
	gameDuration    int // seconds
	expWin          int
	expDraw         int
	activeGames     map[string]*EnhancedGameState
	mutex           sync.RWMutex
}

type EnhancedGameState struct {
	Game          *models.Game
	StartTime     time.Time
	LastManaUpdate time.Time
	GameTimer     *time.Timer
	ManaTimer     *time.Ticker  // Correctly typed as Ticker
	GameEnded     bool
	mutex         sync.RWMutex
}

type EnhancedAction struct {
	Type         string    `json:"type"`        // "spawn_troop"
	TroopID      string    `json:"troop_id"`    
	TargetTower  int       `json:"target_tower"`
	Timestamp    time.Time `json:"timestamp"`
}

type EnhancedResult struct {
	Success       bool          `json:"success"`
	BattleResult  *BattleResult `json:"battle_result,omitempty"`
	PlayerMana    int           `json:"player_mana"`
	GameTimeLeft  int           `json:"game_time_left_seconds"`
	GameEnded     bool          `json:"game_ended"`
	Winner        string        `json:"winner,omitempty"`
	Error         string        `json:"error,omitempty"`
}

func NewEnhancedGameManager(manaRegenRate float64, gameDuration, expWin, expDraw int, critMultiplier float64) *EnhancedGameManager {
	return &EnhancedGameManager{
		battleEngine:  NewBattleEngine(critMultiplier),
		manaRegenRate: manaRegenRate,
		gameDuration:  gameDuration,
		expWin:        expWin,
		expDraw:       expDraw,
		activeGames:   make(map[string]*EnhancedGameState),
	}
}

// StartGame initializes an enhanced mode game
func (egm *EnhancedGameManager) StartGame(game *models.Game) error {
	if len(game.Players) != 2 {
		return fmt.Errorf("need exactly 2 players, got %d", len(game.Players))
	}
	
	// Initialize game state
	game.Start()
	game.Duration = egm.gameDuration
	
	// Initialize players for enhanced mode
	for _, player := range game.Players {
		player.Mana = 5
		player.MaxMana = 10
		player.LastManaUpdate = time.Now()
	}
	
	// Create enhanced game state
	gameState := &EnhancedGameState{
		Game:          game,
		StartTime:     time.Now(),
		LastManaUpdate: time.Now(),
		GameEnded:     false,
	}
	
	// Start game timer
	gameState.GameTimer = time.AfterFunc(time.Duration(egm.gameDuration)*time.Second, func() {
		egm.endGameByTimeout(game.ID)
	})
	
	// Start mana regeneration timer - Now correctly using Ticker
	gameState.ManaTimer = time.NewTicker(time.Second)
	go egm.manageManaRegeneration(gameState)
	
	// Store the game state
	egm.mutex.Lock()
	egm.activeGames[game.ID] = gameState
	egm.mutex.Unlock()
	
	return nil
}

// Rest of the methods remain the same...
func (egm *EnhancedGameManager) ProcessAction(gameID, playerID string, action EnhancedAction) (*EnhancedResult, error) {
	egm.mutex.RLock()
	gameState, exists := egm.activeGames[gameID]
	egm.mutex.RUnlock()
	
	if !exists {
		return nil, errors.New("game not found")
	}
	
	gameState.mutex.Lock()
	defer gameState.mutex.Unlock()
	
	if gameState.GameEnded {
		return &EnhancedResult{
			Success: false,
			Error:   "game has ended",
		}, nil
	}
	
	// Find the player
	var player *models.Player
	for _, p := range gameState.Game.Players {
		if p.ID == playerID {
			player = p
			break
		}
	}
	
	if player == nil {
		return &EnhancedResult{
			Success: false,
			Error:   "player not found",
		}, nil
	}
	
	// Update player's mana
	player.UpdateMana(egm.manaRegenRate)
	
	// Process the action
	switch action.Type {
	case "spawn_troop":
		return egm.processSpawnTroop(gameState, player, action)
	default:
		return &EnhancedResult{
			Success: false,
			Error:   "invalid action type",
		}, nil
	}
}

func (egm *EnhancedGameManager) processSpawnTroop(gameState *EnhancedGameState, player *models.Player, action EnhancedAction) (*EnhancedResult, error) {
	// Find the troop template
	var troopTemplate *models.Troop
	for _, troop := range player.AvailableTroops {
		if troop.ID == action.TroopID {
			troopTemplate = troop
			break
		}
	}
	
	if troopTemplate == nil {
		return &EnhancedResult{
			Success:   false,
			Error:     "troop not available",
			PlayerMana: player.Mana,
		}, nil
	}
	
	// Check mana cost
	if !player.CanSpendMana(troopTemplate.ManaCost) {
		return &EnhancedResult{
			Success:   false,
			Error:     fmt.Sprintf("insufficient mana: need %d, have %d", troopTemplate.ManaCost, player.Mana),
			PlayerMana: player.Mana,
		}, nil
	}
	
	// Execute the attack
	battleResult, err := egm.battleEngine.ExecuteAttack(gameState.Game, player.ID, action.TroopID, action.TargetTower)
	if err != nil {
		return &EnhancedResult{
			Success:    false,
			Error:      err.Error(),
			PlayerMana: player.Mana,
		}, nil
	}
	
	// Calculate remaining time
	timeLeft := egm.gameDuration - int(time.Since(gameState.StartTime).Seconds())
	if timeLeft < 0 {
		timeLeft = 0
	}
	
	result := &EnhancedResult{
		Success:      true,
		BattleResult: battleResult,
		PlayerMana:   player.Mana,
		GameTimeLeft: timeLeft,
		GameEnded:    battleResult.GameEnded,
	}
	
	// Check if game ended due to king tower destruction
	if battleResult.GameEnded {
		egm.endGame(gameState, battleResult.Winner)
		result.Winner = battleResult.Winner
	}
	
	return result, nil
}

func (egm *EnhancedGameManager) manageManaRegeneration(gameState *EnhancedGameState) {
	for range gameState.ManaTimer.C {
		gameState.mutex.Lock()
		if gameState.GameEnded {
			gameState.mutex.Unlock()
			break
		}
		
		// Update mana for all players
		for _, player := range gameState.Game.Players {
			player.UpdateMana(egm.manaRegenRate)
		}
		gameState.mutex.Unlock()
	}
}

func (egm *EnhancedGameManager) endGameByTimeout(gameID string) {
	egm.mutex.RLock()
	gameState, exists := egm.activeGames[gameID]
	egm.mutex.RUnlock()
	
	if !exists || gameState.GameEnded {
		return
	}
	
	gameState.mutex.Lock()
	defer gameState.mutex.Unlock()
	
	// Determine winner based on tower destruction count
	winner := egm.battleEngine.GetGameWinner(gameState.Game)
	egm.endGame(gameState, winner)
}

func (egm *EnhancedGameManager) endGame(gameState *EnhancedGameState, winnerID string) {
	if gameState.GameEnded {
		return
	}
	
	gameState.GameEnded = true
	gameState.Game.State = models.Finished
	endTime := time.Now()
	gameState.Game.EndTime = &endTime
	
	// Set winner
	if winnerID != "" {
		for _, player := range gameState.Game.Players {
			if player.ID == winnerID {
				gameState.Game.Winner = player
				break
			}
		}
	}
	
	// Stop timers
	if gameState.GameTimer != nil {
		gameState.GameTimer.Stop()
	}
	if gameState.ManaTimer != nil {
		gameState.ManaTimer.Stop()
	}
	
	// Award experience points
	egm.awardExperience(gameState.Game)
	
	// Add end game event
	gameState.Game.AddEvent("game_end", "", map[string]interface{}{
		"reason": "time_up",
		"winner": winnerID,
	})
}

func (egm *EnhancedGameManager) awardExperience(game *models.Game) {
	if len(game.Players) != 2 {
		return
	}
	
	player1 := game.Players[0]
	player2 := game.Players[1]
	
	if game.Winner == nil {
		// Draw - both players get draw experience
		player1.AddExperience(egm.expDraw)
		player2.AddExperience(egm.expDraw)
	} else {
		// Winner gets win experience, loser gets 0
		if game.Winner.ID == player1.ID {
			player1.AddExperience(egm.expWin)
		} else {
			player2.AddExperience(egm.expWin)
		}
	}
}

func (egm *EnhancedGameManager) GetGameState(gameID string) (map[string]interface{}, error) {
	egm.mutex.RLock()
	gameState, exists := egm.activeGames[gameID]
	egm.mutex.RUnlock()
	
	if !exists {
		return nil, errors.New("game not found")
	}
	
	gameState.mutex.RLock()
	defer gameState.mutex.RUnlock()
	
	state := make(map[string]interface{})
	game := gameState.Game
	
	state["mode"] = game.Mode
	state["state"] = game.State
	state["duration"] = game.Duration
	
	// Calculate remaining time
	timeLeft := egm.gameDuration - int(time.Since(gameState.StartTime).Seconds())
	if timeLeft < 0 {
		timeLeft = 0
	}
	state["time_remaining"] = timeLeft
	
	// Player information
	players := make([]map[string]interface{}, len(game.Players))
	for i, player := range game.Players {
		// Update mana
		player.UpdateMana(egm.manaRegenRate)
		
		playerData := map[string]interface{}{
			"id":       player.ID,
			"username": player.Username,
			"mana":     player.Mana,
			"max_mana": player.MaxMana,
			"towers":   egm.getTowerStates(player.Towers),
			"troops":   egm.getAvailableTroops(player.AvailableTroops),
		}
		players[i] = playerData
	}
	state["players"] = players
	
	// Game result if finished
	if game.State == models.Finished {
		state["winner"] = ""
		if game.Winner != nil {
			state["winner"] = game.Winner.ID
		}
		
		// Tower destruction counts
		if len(game.Players) == 2 {
			state["tower_scores"] = map[string]int{
				game.Players[0].ID: 3 - egm.battleEngine.CountDestroyedTowers(game.Players[0]),
				game.Players[1].ID: 3 - egm.battleEngine.CountDestroyedTowers(game.Players[1]),
			}
		}
	}
	
	return state, nil
}

func (egm *EnhancedGameManager) getTowerStates(towers []*models.Tower) []map[string]interface{} {
	states := make([]map[string]interface{}, len(towers))
	for i, tower := range towers {
		states[i] = map[string]interface{}{
			"type":     tower.Type,
			"name":     tower.Name,
			"hp":       tower.HP,
			"max_hp":   tower.MaxHP,
			"position": tower.Position,
			"alive":    tower.IsAlive(),
		}
	}
	return states
}

func (egm *EnhancedGameManager) getAvailableTroops(troops []*models.Troop) []map[string]interface{} {
	states := make([]map[string]interface{}, len(troops))
	for i, troop := range troops {
		states[i] = map[string]interface{}{
			"id":        troop.ID,
			"name":      troop.Name,
			"attack":    troop.Attack,
			"defense":   troop.Defense,
			"hp":       troop.HP,
			"mana_cost": troop.ManaCost,
			"crit_chance": troop.CritChance,
		}
	}
	return states
}

func (egm *EnhancedGameManager) CleanupGame(gameID string) {
	egm.mutex.Lock()
	defer egm.mutex.Unlock()
	
	if gameState, exists := egm.activeGames[gameID]; exists {
		gameState.mutex.Lock()
		gameState.GameEnded = true
		if gameState.GameTimer != nil {
			gameState.GameTimer.Stop()
		}
		if gameState.ManaTimer != nil {
			gameState.ManaTimer.Stop()
		}
		gameState.mutex.Unlock()
		
		delete(egm.activeGames, gameID)
	}
}