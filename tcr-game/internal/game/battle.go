// internal/game/battle.go - Battle mechanics
package game

import (
	"fmt"
	"math/rand"
	
	"tcr-game/internal/models"
)

type BattleResult struct {
	AttackerID    string `json:"attacker_id"`
	DefenderID    string `json:"defender_id"`
	TroopUsed     string `json:"troop_used"`
	TargetTower   int    `json:"target_tower"`
	Damage        int    `json:"damage"`
	CriticalHit   bool   `json:"critical_hit"`
	TowerDestroyed bool   `json:"tower_destroyed"`
	CanContinue   bool   `json:"can_continue"`
	GameEnded     bool   `json:"game_ended"`
	Winner        string `json:"winner,omitempty"`
}

type BattleEngine struct {
	critMultiplier float64
}

func NewBattleEngine(critMultiplier float64) *BattleEngine {
	return &BattleEngine{
		critMultiplier: critMultiplier,
	}
}

// CalculateDamage implements the damage formula: DMG = ATK_A - DEF_B (if ≥ 0)
// With critical hit chance for enhanced mode
func (be *BattleEngine) CalculateDamage(attacker *models.Troop, defender models.Attackable, enhancedMode bool) (int, bool) {
	baseDamage := attacker.Attack - defender.GetDefense()
	if baseDamage < 0 {
		baseDamage = 0
	}
	
	// Check for critical hit in enhanced mode
	criticalHit := false
	damage := baseDamage
	
	if enhancedMode && rand.Float64() < attacker.CritChance {
		criticalHit = true
		damage = int(float64(baseDamage) * be.critMultiplier)
	}
	
	return damage, criticalHit
}

// ExecuteAttack performs a single attack from troop to target tower
func (be *BattleEngine) ExecuteAttack(game *models.Game, playerID string, troopID string, targetTowerIndex int) (*BattleResult, error) {
	// Find attacker and defender
	var attacker, defender *models.Player
	for _, player := range game.Players {
		if player.ID == playerID {
			attacker = player
		} else {
			defender = player
		}
	}
	
	if attacker == nil || defender == nil {
		return nil, fmt.Errorf("invalid players")
	}
	
	// Find the troop being used
	var troopUsed *models.Troop
	for _, troop := range attacker.AvailableTroops {
		if troop.ID == troopID {
			troopUsed = troop
			break
		}
	}
	
	if troopUsed == nil {
		return nil, fmt.Errorf("troop not found: %s", troopID)
	}
	
	// Check if target tower index is valid
	if targetTowerIndex < 0 || targetTowerIndex >= len(defender.Towers) {
		return nil, fmt.Errorf("invalid tower index: %d", targetTowerIndex)
	}
	
	// Validate attack rules (Simple mode)
	if game.Mode == models.SimpleMode {
		if err := be.validateSimpleAttackRules(defender, targetTowerIndex); err != nil {
			return nil, err
		}
	}
	
	// Check mana cost (Enhanced mode)
	if game.Mode == models.EnhancedMode {
		if !attacker.CanSpendMana(troopUsed.ManaCost) {
			return nil, fmt.Errorf("insufficient mana: need %d, have %d", troopUsed.ManaCost, attacker.Mana)
		}
		attacker.SpendMana(troopUsed.ManaCost)
	}
	
	// Get target tower
	targetTower := defender.Towers[targetTowerIndex]
	if !targetTower.IsAlive() {
		return nil, fmt.Errorf("target tower is already destroyed")
	}
	
	// Calculate damage
	enhancedMode := game.Mode == models.EnhancedMode
	damage, criticalHit := be.CalculateDamage(troopUsed, targetTower, enhancedMode)
	
	// Apply damage
	targetTower.TakeDamage(damage)
	
	// Create battle result
	result := &BattleResult{
		AttackerID:     attacker.ID,
		DefenderID:     defender.ID,
		TroopUsed:     troopID,
		TargetTower:   targetTowerIndex,
		Damage:        damage,
		CriticalHit:   criticalHit,
		TowerDestroyed: !targetTower.IsAlive(),
		CanContinue:   false,
		GameEnded:     false,
	}
	
	// Check if tower was destroyed
	if result.TowerDestroyed {
		// In simple mode, if a tower is destroyed, the troop can continue attacking
		if game.Mode == models.SimpleMode {
			result.CanContinue = true
		}
		
		// Check for game end conditions
		result.GameEnded, result.Winner = be.checkGameEndConditions(game, defender)
	}
	
	// Add battle event to game
	game.AddEvent("attack", attacker.ID, result)
	
	return result, nil
}

// validateSimpleAttackRules validates attack rules for simple mode
func (be *BattleEngine) validateSimpleAttackRules(defender *models.Player, targetTowerIndex int) error {
	// Rule: Must destroy guard towers before attacking king tower
	if targetTowerIndex == 2 { // King tower
		// Check if both guard towers are destroyed
		if defender.Towers[0].IsAlive() || defender.Towers[1].IsAlive() {
			return fmt.Errorf("must destroy guard towers before attacking king tower")
		}
	}
	
	// Rule: Must destroy first guard tower before second
	if targetTowerIndex == 1 && defender.Towers[0].IsAlive() {
		return fmt.Errorf("must destroy left guard tower before right guard tower")
	}
	
	return nil
}

// checkGameEndConditions checks if the game has ended
func (be *BattleEngine) checkGameEndConditions(game *models.Game, defender *models.Player) (bool, string) {
	// Check if king tower is destroyed
	if !defender.Towers[2].IsAlive() {
		// Find the winner (the other player)
		for _, player := range game.Players {
			if player.ID != defender.ID {
				return true, player.ID
			}
		}
	}
	
	// For enhanced mode, also check tower count at time limit
	if game.Mode == models.EnhancedMode {
		// This will be handled by the game engine when time runs out
		return false, ""
	}
	
	return false, ""
}

// GetValidTargets returns valid tower targets for a player's attack
func (be *BattleEngine) GetValidTargets(game *models.Game, defenderID string) []int {
	var defender *models.Player
	for _, player := range game.Players {
		if player.ID == defenderID {
			defender = player
			break
		}
	}
	
	if defender == nil {
		return []int{}
	}
	
	validTargets := []int{}
	
	if game.Mode == models.SimpleMode {
		// Simple mode rules: attack in order
		if defender.Towers[0].IsAlive() {
			validTargets = append(validTargets, 0)
		} else if defender.Towers[1].IsAlive() {
			validTargets = append(validTargets, 1)
		} else if defender.Towers[2].IsAlive() {
			validTargets = append(validTargets, 2)
		}
	} else {
		// Enhanced mode: can attack any tower
		for i, tower := range defender.Towers {
			if tower.IsAlive() {
				validTargets = append(validTargets, i)
			}
		}
	}
	
	return validTargets
}

// CountDestroyedTowers counts how many towers a player has destroyed
func (be *BattleEngine) CountDestroyedTowers(player *models.Player) int {
	destroyed := 0
	for _, tower := range player.Towers {
		if !tower.IsAlive() {
			destroyed++
		}
	}
	return destroyed
}

// GetGameWinner determines the winner based on current game state
func (be *BattleEngine) GetGameWinner(game *models.Game) string {
	if len(game.Players) != 2 {
		return ""
	}
	
	player1 := game.Players[0]
	player2 := game.Players[1]
	
	// Check for king tower destruction
	if !player1.Towers[2].IsAlive() {
		return player2.ID
	}
	if !player2.Towers[2].IsAlive() {
		return player1.ID
	}
	
	// For enhanced mode, compare tower counts
	if game.Mode == models.EnhancedMode {
		destroyed1 := be.CountDestroyedTowers(player1)
		destroyed2 := be.CountDestroyedTowers(player2)
		
		if destroyed2 > destroyed1 {
			return player2.ID
		} else if destroyed1 > destroyed2 {
			return player1.ID
		}
		// If equal, it's a draw (return empty string)
	}
	
	return ""
}