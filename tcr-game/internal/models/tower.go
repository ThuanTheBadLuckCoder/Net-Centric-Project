// internal/models/tower.go - Tower model
package models

type TowerType string

const (
	KingTower  TowerType = "king_tower"
	GuardTower TowerType = "guard_tower"
)

type Tower struct {
	Type        TowerType `json:"type"`
	Name        string    `json:"name"`
	HP          int       `json:"hp"`
	MaxHP       int       `json:"max_hp"`
	Attack      int       `json:"attack"`
	Defense     int       `json:"defense"`
	CritChance  float64   `json:"crit_chance"`
	Description string    `json:"description"`
	Level       int       `json:"level"`
	Position    int       `json:"position"` // 0=left guard, 1=right guard, 2=king
}

func NewTower(towerType TowerType, name string, hp, attack, defense int, critChance float64, description string, position int) *Tower {
	return &Tower{
		Type:        towerType,
		Name:        name,
		HP:          hp,
		MaxHP:       hp,
		Attack:      attack,
		Defense:     defense,
		CritChance:  critChance,
		Description: description,
		Level:       1,
		Position:    position,
	}
}

func (t *Tower) ApplyLevel(level int) {
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

func (t *Tower) IsAlive() bool {
	return t.HP > 0
}

func (t *Tower) TakeDamage(damage int) {
	t.HP -= damage
	if t.HP < 0 {
		t.HP = 0
	}
}

func (t *Tower) GetDefense() int {
	return t.Defense
}

// Attackable interface for both troops and towers
type Attackable interface {
	TakeDamage(damage int)
	GetDefense() int
	IsAlive() bool
}