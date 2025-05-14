// pkg/protocol/messages.go - Game protocol messages
package protocol

import "time"

// Client to Server messages
type LoginMessage struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateGameMessage struct {
	Type   string `json:"type"`
	GameID string `json:"game_id"`
	Mode   string `json:"mode"`
}

type JoinGameMessage struct {
	Type   string `json:"type"`
	GameID string `json:"game_id"`
}

type SimpleActionMessage struct {
	Type        string `json:"type"`
	TroopID     string `json:"troop_id"`
	TargetTower int    `json:"target_tower"`
}

type EnhancedActionMessage struct {
	Type        string    `json:"type"`
	TroopID     string    `json:"troop_id"`
	TargetTower int       `json:"target_tower"`
	Timestamp   time.Time `json:"timestamp"`
}

// Server to Client messages
type LoginResponse struct {
	Type    string      `json:"type"`
	Success bool        `json:"success"`
	Token   string      `json:"token,omitempty"`
	Player  interface{} `json:"player,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type GameStateMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type ActionResultMessage struct {
	Type   string      `json:"type"`
	Result interface{} `json:"result"`
}

type PlayerJoinedMessage struct {
	Type     string `json:"type"`
	PlayerID string `json:"player_id"`
	Username string `json:"username"`
}

type GameEndMessage struct {
	Type   string `json:"type"`
	Winner string `json:"winner,omitempty"`
	Reason string `json:"reason"`
}