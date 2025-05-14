// internal/models/troop.go - Troop model
package models

import "math/rand"

type Troop struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	HP          int     `json:"hp"`
	MaxHP       int     `json:"max_hp"`
	Attack      int     `json:"attack"`
	Defense     int     `json:"defense"`
	CritChance  float64 `json:"crit_chance"`
	ManaCost    int     `json:"mana_cost"`
	Description string  `json:"description"`
	Level       int     `json:"level"`
}

func NewTroop(id, name string, hp, attack, defense int, critChance float64, manaCost int, description string) *Troop {
	return &Troop{
		ID:          id,
		Name:        name,
		HP:          hp,
		MaxHP:       hp,
		Attack:      attack,
		Defense:     defense,
		CritChance:  critChance,
		ManaCost:    manaCost,
		Description: description,
		Level:       1,
	}
}

func (t *Troop) ApplyLevel(level int) {
	if level <= 0 {
		level = 1
	}
	t.Level = level
	
	// Each level increases stats by 10%
	multiplier := 1.0 + (0.1 * float64(level-1))
	
	t.HP = int(float64(t.MaxHP) * multiplier)
	t.MaxHP = t.HP
	t.Attack = int(float64(t.Attack) * multiplier)
	t.Defense = int(float64(t.Defense) * multiplier)
}

func (t *Troop) IsAlive() bool {
	return t.HP > 0
}

func (t *Troop) TakeDamage(damage int) {
	t.HP -= damage
	if t.HP < 0 {
		t.HP = 0
	}
}

func (t *Troop) AttackTarget(target Attackable) int {
	// Calculate base damage
	baseDamage := t.Attack - target.GetDefense()
	if baseDamage < 0 {
		baseDamage = 0
	}
	
	// Check for critical hit
	damage := baseDamage
	if rand.Float64() < t.CritChance {
		damage = int(float64(baseDamage) * 1.2) // 20% crit bonus
	}
	
	// Apply damage to target
	target.TakeDamage(damage)
	
	return damage
}