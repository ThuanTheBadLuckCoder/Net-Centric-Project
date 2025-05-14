// internal/models/player.go - Player model
package models

import "time"

type Player struct {
	ID          string            `json:"id"`
	Username    string            `json:"username"`
	Password    string            `json:"password"` // In production, use hashed passwords
	Experience  int               `json:"experience"`
	Level       int               `json:"level"`
	TroopLevels map[string]int    `json:"troop_levels"`
	TowerLevels map[TowerType]int `json:"tower_levels"`
	Stats       PlayerStats       `json:"stats"`
	Towers      []*Tower          `json:"towers"`
	AvailableTroops []*Troop      `json:"available_troops"`
	Mana        int               `json:"mana"`
	MaxMana     int               `json:"max_mana"`
	LastManaUpdate time.Time      `json:"last_mana_update"`
}

type PlayerStats struct {
	GamesPlayed int `json:"games_played"`
	GamesWon    int `json:"games_won"`
	GamesLost   int `json:"games_lost"`
	GamesDrawn  int `json:"games_drawn"`
}

func NewPlayer(id, username, password string) *Player {
	return &Player{
		ID:          id,
		Username:    username,
		Password:    password,
		Experience:  0,
		Level:       1,
		TroopLevels: make(map[string]int),
		TowerLevels: make(map[TowerType]int),
		Stats:       PlayerStats{},
		Towers:      make([]*Tower, 3),
		AvailableTroops: make([]*Troop, 0),
		Mana:        5,
		MaxMana:     10,
		LastManaUpdate: time.Now(),
	}
}

func (p *Player) AddExperience(exp int) {
	p.Experience += exp
	// Simple leveling: every 100 EXP = 1 level
	newLevel := (p.Experience / 100) + 1
	if newLevel > p.Level {
		p.Level = newLevel
	}
}

func (p *Player) CanSpendMana(cost int) bool {
	return p.Mana >= cost
}

func (p *Player) SpendMana(cost int) bool {
	if p.CanSpendMana(cost) {
		p.Mana -= cost
		return true
	}
	return false
}

func (p *Player) UpdateMana(manaRegenRate float64) {
	now := time.Now()
	elapsed := now.Sub(p.LastManaUpdate).Seconds()
	
	manaToAdd := int(elapsed * manaRegenRate)
	if manaToAdd > 0 {
		p.Mana += manaToAdd
		if p.Mana > p.MaxMana {
			p.Mana = p.MaxMana
		}
		p.LastManaUpdate = now
	}
}

func (p *Player) GetTowerLevel(towerType TowerType) int {
	if level, exists := p.TowerLevels[towerType]; exists {
		return level
	}
	return 1
}

func (p *Player) GetTroopLevel(troopID string) int {
	if level, exists := p.TroopLevels[troopID]; exists {
		return level
	}
	return 1
}

func (p *Player) InitializeTowers(towerSpecs map[string]interface{}) {
	// Left Guard Tower
	p.Towers[0] = NewTower(GuardTower, "Left Guard Tower", 300, 20, 10, 0.05, "Left defensive tower", 0)
	p.Towers[0].ApplyLevel(p.GetTowerLevel(GuardTower))
	
	// Right Guard Tower
	p.Towers[1] = NewTower(GuardTower, "Right Guard Tower", 300, 20, 10, 0.05, "Right defensive tower", 1)
	p.Towers[1].ApplyLevel(p.GetTowerLevel(GuardTower))
	
	// King Tower
	p.Towers[2] = NewTower(KingTower, "King Tower", 500, 25, 15, 0.1, "Main tower", 2)
	p.Towers[2].ApplyLevel(p.GetTowerLevel(KingTower))
}