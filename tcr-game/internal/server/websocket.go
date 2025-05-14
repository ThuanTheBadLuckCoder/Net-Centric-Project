// Fixed internal/server/websocket.go - Completely clean version
package server

import (
	"log"
	"net/http"
	"sync"
	
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WebSocketManager struct {
	connections map[string]map[*websocket.Conn]bool // gameID -> connections
	mutex       sync.RWMutex
	upgrader    websocket.Upgrader
}

type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		connections: make(map[string]map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in development
			},
		},
	}
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	
	// Validate token
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}
	
	player, err := s.authService.ValidateToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	
	// Upgrade connection
	conn, err := s.wsManager.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()
	
	// Add connection to game
	s.wsManager.AddConnection(gameID, conn)
	defer s.wsManager.RemoveConnection(gameID, conn)
	
	// Send initial game state
	state, err := s.gameEngine.GetGameState(gameID)
	if err == nil {
		s.wsManager.SendToConnection(conn, WSMessage{
			Type: "game_state",
			Data: state,
		})
	}
	
	// Send player joined event
	s.wsManager.BroadcastToGame(gameID, WSMessage{
		Type: "player_joined",
		Data: map[string]interface{}{
			"player_id": player.ID,
			"username":  player.Username,
		},
	})
	
	// Handle incoming messages
	for {
		var msg WSMessage
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}
		
		// Handle different message types
		switch msg.Type {
		case "ping":
			s.wsManager.SendToConnection(conn, WSMessage{Type: "pong", Data: nil})
		case "get_state":
			state, err := s.gameEngine.GetGameState(gameID)
			if err == nil {
				s.wsManager.SendToConnection(conn, WSMessage{
					Type: "game_state",
					Data: state,
				})
			}
		}
	}
}

func (wsm *WebSocketManager) AddConnection(gameID string, conn *websocket.Conn) {
	wsm.mutex.Lock()
	defer wsm.mutex.Unlock()
	
	if wsm.connections[gameID] == nil {
		wsm.connections[gameID] = make(map[*websocket.Conn]bool)
	}
	wsm.connections[gameID][conn] = true
}

func (wsm *WebSocketManager) RemoveConnection(gameID string, conn *websocket.Conn) {
	wsm.mutex.Lock()
	defer wsm.mutex.Unlock()
	
	if connections, exists := wsm.connections[gameID]; exists {
		delete(connections, conn)
		if len(connections) == 0 {
			delete(wsm.connections, gameID)
		}
	}
}

func (wsm *WebSocketManager) BroadcastToGame(gameID string, message interface{}) {
	wsm.mutex.RLock()
	connections, exists := wsm.connections[gameID]
	wsm.mutex.RUnlock()
	
	if !exists {
		return
	}
	
	wsm.mutex.RLock()
	for conn := range connections {
		go wsm.SendToConnection(conn, message)
	}
	wsm.mutex.RUnlock()
}

func (wsm *WebSocketManager) SendToConnection(conn *websocket.Conn, message interface{}) {
	if err := conn.WriteJSON(message); err != nil {
		log.Printf("WebSocket write error: %v", err)
		conn.Close()
	}
}