// internal/utils/calculator.go - Damage calculation utilities
package utils

import (
	"math/rand"
	
	"tcr-game/internal/models"
)

func CalculateDamage(attacker *models.Troop, defender models.Attackable, critMultiplier float64) (int, bool) {
	baseDamage := attacker.Attack - defender.GetDefense()
	if baseDamage < 0 {
		baseDamage = 0
	}
	
	criticalHit := rand.Float64() < attacker.CritChance
	damage := baseDamage
	
	if criticalHit {
		damage = int(float64(baseDamage) * critMultiplier)
	}
	
	return damage, criticalHit
}

func CalculateRequiredExp(level int) int {
	// Each level requires 10% more EXP than the previous
	baseExp := 100
	required := baseExp
	
	for i := 1; i < level; i++ {
		required = int(float64(required) * 1.1)
	}
	
	return required
}