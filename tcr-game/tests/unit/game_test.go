// Fixed tests/unit/game_test.go - Fix printf format issue
package unit

import (
	"testing"
	
	"tcr-game/internal/game"
	"tcr-game/internal/models"
)

func TestBattleEngine_CalculateDamage(t *testing.T) {
	engine := game.NewBattleEngine(1.2)
	
	// Create test troop and tower
	troop := &models.Troop{
		ID:         "test_troop",
		Attack:     30,
		CritChance: 0.0, // No crit for predictable testing
	}
	
	tower := &models.Tower{
		Defense: 10,
		HP:      100,
		MaxHP:   100,
	}
	
	damage, crit := engine.CalculateDamage(troop, tower, false)
	
	if damage != 20 {
		t.Errorf("Expected damage 20, got %d", damage)
	}
	
	if crit {
		t.Errorf("Expected no crit with 0%% chance")  // Fixed: escaped the % character
	}
}

func TestSimpleGameManager_StartGame(t *testing.T) {
	manager := game.NewSimpleGameManager(2, 30, 1.2)
	gameObj := models.NewGame("test_game", models.SimpleMode)
	
	// Add two players
	player1 := models.NewPlayer("p1", "player1", "pass1")
	player2 := models.NewPlayer("p2", "player2", "pass2")
	
	// Add some troops to players
	troop1 := &models.Troop{ID: "goblin", Name: "Goblin", HP: 100, MaxHP: 100}
	troop2 := &models.Troop{ID: "archer", Name: "Archer", HP: 80, MaxHP: 80}
	troop3 := &models.Troop{ID: "knight", Name: "Knight", HP: 150, MaxHP: 150}
	
	player1.AvailableTroops = []*models.Troop{troop1, troop2, troop3}
	player2.AvailableTroops = []*models.Troop{troop1, troop2, troop3}
	
	gameObj.AddPlayer(player1)
	gameObj.AddPlayer(player2)
	
	err := manager.StartGame(gameObj)
	if err != nil {
		t.Errorf("Failed to start game: %v", err)
	}
	
	if gameObj.State != models.InProgress {
		t.Errorf("Expected game state to be InProgress, got %v", gameObj.State)
	}
}