// Fixed internal/server/handlers.go - Move types to game package
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	
	"github.com/gorilla/mux"
	"tcr-game/internal/auth"
	"tcr-game/internal/game"
	"tcr-game/internal/models"
)

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var creds auth.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	player, err := s.authService.Register(creds.Username, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	response := map[string]interface{}{
		"success": true,
		"player":  player,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var creds auth.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	response, err := s.authService.Login(creds.Username, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Missing authorization token", http.StatusUnauthorized)
		return
	}
	
	s.authService.Logout(token)
	
	response := map[string]interface{}{
		"success": true,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleCreateGame(w http.ResponseWriter, r *http.Request) {
	player, err := s.validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	var request struct {
		Mode   string `json:"mode"`
		GameID string `json:"game_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	var gameMode models.GameMode
	switch request.Mode {
	case "simple":
		gameMode = models.SimpleMode
	case "enhanced":
		gameMode = models.EnhancedMode
	default:
		http.Error(w, "Invalid game mode", http.StatusBadRequest)
		return
	}
	
	game, err := s.gameEngine.CreateGame(request.GameID, gameMode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Join the game immediately
	if err := s.gameEngine.JoinGame(request.GameID, player); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	response := map[string]interface{}{
		"success": true,
		"game":    game,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleJoinGame(w http.ResponseWriter, r *http.Request) {
	player, err := s.validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	
	if err := s.gameEngine.JoinGame(gameID, player); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	gameObj, _ := s.gameEngine.GetGame(gameID)
	
	response := map[string]interface{}{
		"success": true,
		"game":    gameObj,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleGetGameState(w http.ResponseWriter, r *http.Request) {
	_, err := s.validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	
	state, err := s.gameEngine.GetGameState(gameID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

func (s *Server) handleGameAction(w http.ResponseWriter, r *http.Request) {
	player, err := s.validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	
	gameObj, err := s.gameEngine.GetGame(gameID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	var response interface{}
	
	if gameObj.Mode == models.SimpleMode {
		var action game.TurnAction
		if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		
		result, err := s.gameEngine.ProcessSimpleAction(gameID, player.ID, action)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response = result
		
		// Broadcast to WebSocket clients
		s.wsManager.BroadcastToGame(gameID, map[string]interface{}{
			"type":   "turn_result",
			"result": result,
		})
	} else {
		var action game.EnhancedAction
		if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		action.Timestamp = time.Now()
		
		result, err := s.gameEngine.ProcessEnhancedAction(gameID, player.ID, action)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response = result
		
		// Broadcast to WebSocket clients
		s.wsManager.BroadcastToGame(gameID, map[string]interface{}{
			"type":   "action_result",
			"result": result,
		})
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) validateToken(r *http.Request) (*models.Player, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return nil, fmt.Errorf("missing authorization token")
	}
	
	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	
	return s.authService.ValidateToken(token)
}