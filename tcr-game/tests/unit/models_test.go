// tests/unit/models_test.go - Model tests
package unit

import (
	"testing"
	
	"tcr-game/internal/models"
)

func TestPlayer_ApplyExperience(t *testing.T) {
	player := models.NewPlayer("test", "testuser", "pass")
	
	player.AddExperience(150)
	
	if player.Experience != 150 {
		t.Errorf("Expected experience 150, got %d", player.Experience)
	}
	
	if player.Level != 2 {
		t.Errorf("Expected level 2, got %d", player.Level)
	}
}

func TestTroop_ApplyLevel(t *testing.T) {
	troop := models.NewTroop("goblin", "Goblin", 100, 30, 5, 0.1, 2, "Test troop")
	
	troop.ApplyLevel(3)
	
	expectedHP := int(float64(100) * 1.2) // 10% increase per level above 1
	if troop.HP != expectedHP {
		t.Errorf("Expected HP %d, got %d", expectedHP, troop.HP)
	}
	
	expectedAttack := int(float64(30) * 1.2)
	if troop.Attack != expectedAttack {
		t.Errorf("Expected attack %d, got %d", expectedAttack, troop.Attack)
	}
}

func TestTower_TakeDamage(t *testing.T) {
	tower := models.NewTower(models.GuardTower, "Guard", 300, 20, 10, 0.05, "Test tower", 0)
	
	tower.TakeDamage(100)
	
	if tower.HP != 200 {
		t.Errorf("Expected HP 200, got %d", tower.HP)
	}
	
	if !tower.IsAlive() {
		t.Errorf("Expected tower to be alive")
	}
	
	tower.TakeDamage(300)
	
	if tower.HP != 0 {
		t.Errorf("Expected HP 0, got %d", tower.HP)
	}
	
	if tower.IsAlive() {
		t.Errorf("Expected tower to be dead")
	}
}