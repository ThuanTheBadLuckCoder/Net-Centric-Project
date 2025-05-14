// Fixed internal/game/simple_rules.go - Add types to game package
package game

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
	
	"tcr-game/internal/models"
)

type SimpleGameManager struct {
	battleEngine *BattleEngine
	maxPlayers   int
	turnTime     int // seconds
}

// Define TurnAction in game package
type TurnAction struct {
	Type       string `json:"type"`        // "attack"
	TroopID    string `json:"troop_id"`    
	TargetTower int   `json:"target_tower"`
}

type TurnResult struct {
	Success        bool         `json:"success"`
	BattleResult   *BattleResult `json:"battle_result,omitempty"`
	CanContinue    bool         `json:"can_continue"`
	NextPlayer     string       `json:"next_player"`
	TurnRemaining  int          `json:"turn_remaining_seconds"`
	Error          string       `json:"error,omitempty"`
}

func NewSimpleGameManager(maxPlayers, turnTime int, critMultiplier float64) *SimpleGameManager {
	return &SimpleGameManager{
		battleEngine: NewBattleEngine(critMultiplier),
		maxPlayers:   maxPlayers,
		turnTime:     turnTime,
	}
}

// StartGame initializes a simple mode game
func (sgm *SimpleGameManager) StartGame(game *models.Game) error {
	if len(game.Players) != sgm.maxPlayers {
		return fmt.Errorf("need exactly %d players, got %d", sgm.maxPlayers, len(game.Players))
	}
	
	// Assign random troops to each player
	for _, player := range game.Players {
		if err := sgm.assignRandomTroops(player); err != nil {
			return fmt.Errorf("failed to assign troops to player %s: %v", player.ID, err)
		}
	}
	
	// Initialize game state
	game.Start()
	game.CurrentTurn = 0 // First player's turn
	
	// Initialize mana for players (not used in simple mode but kept for consistency)
	for _, player := range game.Players {
		player.Mana = 5
		player.MaxMana = 10
	}
	
	return nil
}

// assignRandomTroops gives each player 3 random troops from the available list
func (sgm *SimpleGameManager) assignRandomTroops(player *models.Player) error {
	// Load all available troops
	allTroops := player.AvailableTroops
	if len(allTroops) == 0 {
		return errors.New("no troops available")
	}
	
	// Randomly select 3 troops
	selected := make([]*models.Troop, 0, 3)
	usedIndices := make(map[int]bool)
	
	for len(selected) < 3 && len(selected) < len(allTroops) {
		index := rand.Intn(len(allTroops))
		if !usedIndices[index] {
			usedIndices[index] = true
			// Create a copy of the troop for this game
			troop := *allTroops[index]
			selected = append(selected, &troop)
		}
	}
	
	player.AvailableTroops = selected
	return nil
}

// ProcessTurn handles a player's turn in simple mode
func (sgm *SimpleGameManager) ProcessTurn(game *models.Game, playerID string, action TurnAction) (*TurnResult, error) {
	// Validate it's the player's turn
	currentPlayer := game.Players[game.CurrentTurn]
	if currentPlayer.ID != playerID {
		return &TurnResult{
			Success: false,
			Error:   "not your turn",
		}, nil
	}
	
	// Validate action type
	if action.Type != "attack" {
		return &TurnResult{
			Success: false,
			Error:   "invalid action type",
		}, nil
	}
	
	// Execute the attack
	battleResult, err := sgm.battleEngine.ExecuteAttack(game, playerID, action.TroopID, action.TargetTower)
	if err != nil {
		return &TurnResult{
			Success: false,
			Error:   err.Error(),
		}, nil
	}
	
	// Check if the game ended
	if battleResult.GameEnded {
		game.State = models.Finished
		endTime := time.Now()
		game.EndTime = &endTime
		game.Winner = sgm.findPlayerByID(game, battleResult.Winner)
		
		return &TurnResult{
			Success:      true,
			BattleResult: battleResult,
			CanContinue:  false,
			NextPlayer:   "",
		}, nil
	}
	
	// In simple mode, if a tower is destroyed, the same player continues
	canContinue := battleResult.TowerDestroyed
	nextPlayerTurn := game.CurrentTurn
	
	if !canContinue {
		// Switch to next player
		nextPlayerTurn = (game.CurrentTurn + 1) % len(game.Players)
		game.CurrentTurn = nextPlayerTurn
	}
	
	nextPlayer := game.Players[nextPlayerTurn]
	
	return &TurnResult{
		Success:       true,
		BattleResult:  battleResult,
		CanContinue:   canContinue,
		NextPlayer:    nextPlayer.ID,
		TurnRemaining: sgm.turnTime,
	}, nil
}

// ValidateTroopSelection checks if a troop can be used by the player
func (sgm *SimpleGameManager) ValidateTroopSelection(player *models.Player, troopID string) (*models.Troop, error) {
	for _, troop := range player.AvailableTroops {
		if troop.ID == troopID && troop.IsAlive() {
			return troop, nil
		}
	}
	return nil, errors.New("troop not available or already used")
}

// GetGameState returns the current state of the game for simple mode
func (sgm *SimpleGameManager) GetGameState(game *models.Game) map[string]interface{} {
	state := make(map[string]interface{})
	
	state["mode"] = game.Mode
	state["state"] = game.State
	state["current_turn"] = game.CurrentTurn
	
	// Player information
	players := make([]map[string]interface{}, len(game.Players))
	for i, player := range game.Players {
		playerData := map[string]interface{}{
			"id":       player.ID,
			"username": player.Username,
			"towers":   sgm.getTowerStates(player.Towers),
			"troops":   sgm.getTroopStates(player.AvailableTroops),
		}
		players[i] = playerData
	}
	state["players"] = players
	
	// Current player information
	if game.State == models.InProgress && game.CurrentTurn < len(game.Players) {
		currentPlayer := game.Players[game.CurrentTurn]
		state["current_player"] = map[string]interface{}{
			"id":            currentPlayer.ID,
			"username":      currentPlayer.Username,
			"valid_targets": sgm.battleEngine.GetValidTargets(game, sgm.getOpponentID(game, currentPlayer.ID)),
		}
	}
	
	// Game result if finished
	if game.State == models.Finished {
		state["winner"] = ""
		if game.Winner != nil {
			state["winner"] = game.Winner.ID
		}
	}
	
	return state
}

// getTowerStates converts towers to a serializable format
func (sgm *SimpleGameManager) getTowerStates(towers []*models.Tower) []map[string]interface{} {
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

// getTroopStates converts troops to a serializable format
func (sgm *SimpleGameManager) getTroopStates(troops []*models.Troop) []map[string]interface{} {
	states := make([]map[string]interface{}, len(troops))
	for i, troop := range troops {
		states[i] = map[string]interface{}{
			"id":     troop.ID,
			"name":   troop.Name,
			"hp":     troop.HP,
			"max_hp": troop.MaxHP,
			"attack": troop.Attack,
			"defense": troop.Defense,
			"alive":  troop.IsAlive(),
		}
	}
	return states
}

// Helper functions
func (sgm *SimpleGameManager) findPlayerByID(game *models.Game, playerID string) *models.Player {
	for _, player := range game.Players {
		if player.ID == playerID {
			return player
		}
	}
	return nil
}

func (sgm *SimpleGameManager) getOpponentID(game *models.Game, playerID string) string {
	for _, player := range game.Players {
		if player.ID != playerID {
			return player.ID
		}
	}
	return ""
}

// EndGame handles the end of a simple mode game
func (sgm *SimpleGameManager) EndGame(game *models.Game, reason string) error {
	if game.State == models.Finished {
		return errors.New("game already finished")
	}
	
	// Determine winner if not already set
	if game.Winner == nil {
		winnerID := sgm.battleEngine.GetGameWinner(game)
		if winnerID != "" {
			game.Winner = sgm.findPlayerByID(game, winnerID)
		}
	}
	
	game.State = models.Finished
	endTime := time.Now()
	game.EndTime = &endTime
	
	// Record end game event
	game.AddEvent("game_end", "", map[string]interface{}{
		"reason": reason,
		"winner": "",
	})
	if game.Winner != nil {
		game.Events[len(game.Events)-1].Data.(map[string]interface{})["winner"] = game.Winner.ID
	}
	
	return nil
}

// GetAvailableActions returns what actions the current player can take
func (sgm *SimpleGameManager) GetAvailableActions(game *models.Game, playerID string) ([]string, error) {
	if game.State != models.InProgress {
		return []string{}, errors.New("game not in progress")
	}
	
	currentPlayer := game.Players[game.CurrentTurn]
	if currentPlayer.ID != playerID {
		return []string{}, errors.New("not your turn")
	}
	
	actions := []string{}
	
	// Check available troops for attack
	for _, troop := range currentPlayer.AvailableTroops {
		if troop.IsAlive() {
			actions = append(actions, fmt.Sprintf("attack_with_%s", troop.ID))
		}
	}
	
	return actions, nil
}