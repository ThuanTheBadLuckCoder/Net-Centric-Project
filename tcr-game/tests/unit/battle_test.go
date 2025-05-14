// tests/unit/battle_test.go - Battle mechanics tests
package unit

import (
	"testing"
	
	"tcr-game/internal/game"
	"tcr-game/internal/models"
)

func TestBattleEngine_ExecuteAttack(t *testing.T) {
	engine := game.NewBattleEngine(1.2)
	game := models.NewGame("test_battle", models.SimpleMode)
	
	// Setup players
	player1 := models.NewPlayer("p1", "player1", "pass1")
	player2 := models.NewPlayer("p2", "player2", "pass2")
	
	// Initialize towers for player2
	player2.InitializeTowers(nil)
	
	// Add troop to player1
	troop := &models.Troop{
		ID:      "goblin",
		Name:    "Goblin",
		HP:      100,
		MaxHP:   100,
		Attack:  30,
		Defense: 5,
	}
	player1.AvailableTroops = []*models.Troop{troop}
	
	game.AddPlayer(player1)
	game.AddPlayer(player2)
	game.State = models.InProgress
	
	// Execute attack
	result, err := engine.ExecuteAttack(game, player1.ID, "goblin", 0)
	
	if err != nil {
		t.Errorf("Attack failed: %v", err)
	}
	
	if result.AttackerID != player1.ID {
		t.Errorf("Expected attacker ID %s, got %s", player1.ID, result.AttackerID)
	}
	
	if result.DefenderID != player2.ID {
		t.Errorf("Expected defender ID %s, got %s", player2.ID, result.DefenderID)
	}
}

func TestBattleEngine_GetValidTargets(t *testing.T) {
	engine := game.NewBattleEngine(1.2)
	game := models.NewGame("test_targets", models.SimpleMode)
	
	player := models.NewPlayer("p1", "player1", "pass1")
	player.InitializeTowers(nil)
	
	// Destroy first guard tower
	player.Towers[0].HP = 0
	
	targets := engine.GetValidTargets(game, player.ID)
	
	// In simple mode, only second guard tower should be valid
	if len(targets) != 1 || targets[0] != 1 {
		t.Errorf("Expected valid target [1], got %v", targets)
	}
}