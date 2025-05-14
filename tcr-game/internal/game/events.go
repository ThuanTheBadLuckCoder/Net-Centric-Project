// Fixed internal/game/events.go - Remove unused encoding/json import
package game

import (
	"sync"
	"time"
	
	"tcr-game/internal/models"
)

type EventType string

const (
	EventPlayerJoined    EventType = "player_joined"
	EventGameStarted     EventType = "game_started"
	EventTurnChanged     EventType = "turn_changed"
	EventAttackMade      EventType = "attack_made"
	EventTowerDestroyed  EventType = "tower_destroyed"
	EventGameEnded       EventType = "game_ended"
	EventManaUpdated     EventType = "mana_updated"
)

type GameEventData struct {
	Type      EventType   `json:"type"`
	GameID    string      `json:"game_id"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type EventManager struct {
	subscribers map[string][]chan GameEventData
	mutex       sync.RWMutex
}

func NewEventManager() *EventManager {
	return &EventManager{
		subscribers: make(map[string][]chan GameEventData),
	}
}

func (em *EventManager) Subscribe(gameID string, ch chan GameEventData) {
	em.mutex.Lock()
	defer em.mutex.Unlock()
	
	em.subscribers[gameID] = append(em.subscribers[gameID], ch)
}

func (em *EventManager) Unsubscribe(gameID string, ch chan GameEventData) {
	em.mutex.Lock()
	defer em.mutex.Unlock()
	
	if subscribers, exists := em.subscribers[gameID]; exists {
		for i, subscriber := range subscribers {
			if subscriber == ch {
				em.subscribers[gameID] = append(subscribers[:i], subscribers[i+1:]...)
				break
			}
		}
		
		// Clean up empty game channels
		if len(em.subscribers[gameID]) == 0 {
			delete(em.subscribers, gameID)
		}
	}
}

func (em *EventManager) Publish(event GameEventData) {
	em.mutex.RLock()
	subscribers := em.subscribers[event.GameID]
	em.mutex.RUnlock()
	
	if len(subscribers) == 0 {
		return
	}
	
	// Create a copy to avoid concurrent access issues
	subscribersCopy := make([]chan GameEventData, len(subscribers))
	copy(subscribersCopy, subscribers)
	
	// Send to all subscribers in goroutines to avoid blocking
	for _, subscriber := range subscribersCopy {
		go func(ch chan GameEventData) {
			select {
			case ch <- event:
				// Event sent successfully
			default:
				// Channel is full or closed, skip
			}
		}(subscriber)
	}
}

func (em *EventManager) PublishPlayerJoined(gameID string, player *models.Player) {
	em.Publish(GameEventData{
		Type:      EventPlayerJoined,
		GameID:    gameID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"player_id": player.ID,
			"username":  player.Username,
		},
	})
}

func (em *EventManager) PublishGameStarted(gameID string, game *models.Game) {
	em.Publish(GameEventData{
		Type:      EventGameStarted,
		GameID:    gameID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"mode":    game.Mode,
			"players": len(game.Players),
		},
	})
}

func (em *EventManager) PublishTurnChanged(gameID string, currentPlayerID string) {
	em.Publish(GameEventData{
		Type:      EventTurnChanged,
		GameID:    gameID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"current_player": currentPlayerID,
		},
	})
}

func (em *EventManager) PublishAttackMade(gameID string, result *BattleResult) {
	em.Publish(GameEventData{
		Type:      EventAttackMade,
		GameID:    gameID,
		Timestamp: time.Now(),
		Data:      result,
	})
}

func (em *EventManager) PublishTowerDestroyed(gameID string, towerIndex int, playerID string) {
	em.Publish(GameEventData{
		Type:      EventTowerDestroyed,
		GameID:    gameID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"tower_index": towerIndex,
			"player_id":   playerID,
		},
	})
}

func (em *EventManager) PublishGameEnded(gameID string, winner string, reason string) {
	em.Publish(GameEventData{
		Type:      EventGameEnded,
		GameID:    gameID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"winner": winner,
			"reason": reason,
		},
	})
}

func (em *EventManager) PublishManaUpdated(gameID string, playerID string, mana int) {
	em.Publish(GameEventData{
		Type:      EventManaUpdated,
		GameID:    gameID,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"player_id": playerID,
			"mana":      mana,
		},
	})
}

// CleanupGame removes all subscribers for a specific game
func (em *EventManager) CleanupGame(gameID string) {
	em.mutex.Lock()
	defer em.mutex.Unlock()
	
	// Close all channels for this game
	if subscribers, exists := em.subscribers[gameID]; exists {
		for _, ch := range subscribers {
			close(ch)
		}
		delete(em.subscribers, gameID)
	}
}

// GetSubscriberCount returns the number of subscribers for a game
func (em *EventManager) GetSubscriberCount(gameID string) int {
	em.mutex.RLock()
	defer em.mutex.RUnlock()
	
	return len(em.subscribers[gameID])
}