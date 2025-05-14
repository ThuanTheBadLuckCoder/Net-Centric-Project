// internal/models/gamestate.go - Game state model
package models

import "time"

type GameMode string

const (
	SimpleMode   GameMode = "simple"
	EnhancedMode GameMode = "enhanced"
)

type GameState string

const (
	Waiting    GameState = "waiting"
	InProgress GameState = "in_progress"
	Finished   GameState = "finished"
)

type Game struct {
	ID          string           `json:"id"`
	Mode        GameMode         `json:"mode"`
	State       GameState        `json:"state"`
	Players     []*Player        `json:"players"`
	CurrentTurn int              `json:"current_turn"`
	StartTime   time.Time        `json:"start_time"`
	EndTime     *time.Time       `json:"end_time,omitempty"`
	Duration    int              `json:"duration"` // seconds for enhanced mode
	Winner      *Player          `json:"winner,omitempty"`
	Events      []GameEvent      `json:"events"`
}

type GameEvent struct {
	Type      string      `json:"type"`
	PlayerID  string      `json:"player_id"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

func NewGame(id string, mode GameMode) *Game {
	return &Game{
		ID:          id,
		Mode:        mode,
		State:       Waiting,
		Players:     make([]*Player, 0, 2),
		CurrentTurn: 0,
		Events:      make([]GameEvent, 0),
	}
}

func (g *Game) AddPlayer(player *Player) bool {
	if len(g.Players) >= 2 {
		return false
	}
	g.Players = append(g.Players, player)
	return true
}

func (g *Game) Start() {
	g.State = InProgress
	g.StartTime = time.Now()
	
	// Initialize towers for both players
	for _, player := range g.Players {
		player.InitializeTowers(nil) // TODO: Pass actual tower specs
	}
}

func (g *Game) AddEvent(eventType, playerID string, data interface{}) {
	event := GameEvent{
		Type:      eventType,
		PlayerID:  playerID,
		Timestamp: time.Now(),
		Data:      data,
	}
	g.Events = append(g.Events, event)
}

func (g *Game) IsFinished() bool {
	return g.State == Finished
}

func (g *Game) GetOpponent(playerID string) *Player {
	for _, player := range g.Players {
		if player.ID != playerID {
			return player
		}
	}
	return nil
}