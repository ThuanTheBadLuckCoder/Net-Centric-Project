// internal/models/battle.go - Battle result model
package models

import "time"

type BattleAction struct {
	ID         string    `json:"id"`
	GameID     string    `json:"game_id"`
	PlayerID   string    `json:"player_id"`
	TroopID    string    `json:"troop_id"`
	TargetID   string    `json:"target_id"`
	Damage     int       `json:"damage"`
	CritHit    bool      `json:"crit_hit"`
	Timestamp  time.Time `json:"timestamp"`
}

type BattleState struct {
	Attacker   *Player `json:"attacker"`
	Defender   *Player `json:"defender"`
	TroopUsed  *Troop  `json:"troop_used"`
	Target     *Tower  `json:"target"`
	Damage     int     `json:"damage"`
	CritHit    bool    `json:"crit_hit"`
	Destroyed  bool    `json:"destroyed"`
}