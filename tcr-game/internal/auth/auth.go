// internal/auth/auth.go - Authentication logic
package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
	
	"tcr-game/internal/models"
	"tcr-game/internal/storage"
)

type AuthService struct {
	storage  *storage.JSONStorage
	sessions map[string]*Session
}

type Session struct {
	PlayerID  string
	Token     string
	ExpiresAt time.Time
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool           `json:"success"`
	Token   string         `json:"token,omitempty"`
	Player  *models.Player `json:"player,omitempty"`
	Error   string         `json:"error,omitempty"`
}

func NewAuthService(storage *storage.JSONStorage) *AuthService {
	return &AuthService{
		storage:  storage,
		sessions: make(map[string]*Session),
	}
}

func (as *AuthService) hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func (as *AuthService) generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (as *AuthService) Register(username, password string) (*models.Player, error) {
	// Check if player already exists
	playerID := as.generatePlayerID(username)
	if _, err := as.storage.LoadPlayer(playerID); err == nil {
		return nil, errors.New("player already exists")
	}
	
	// Create new player
	hashedPassword := as.hashPassword(password)
	player := models.NewPlayer(playerID, username, hashedPassword)
	
	// Initialize default troop and tower levels
	player.TroopLevels = map[string]int{
		"goblin": 1,
		"archer": 1,
		"knight": 1,
		"wizard": 1,
		"dragon": 1,
	}
	player.TowerLevels = map[models.TowerType]int{
		models.KingTower:  1,
		models.GuardTower: 1,
	}
	
	// Save player to storage
	if err := as.storage.SavePlayer(player); err != nil {
		return nil, fmt.Errorf("failed to save player: %v", err)
	}
	
	return player, nil
}

func (as *AuthService) Login(username, password string) (*LoginResponse, error) {
	playerID := as.generatePlayerID(username)
	
	// Load player from storage
	player, err := as.storage.LoadPlayer(playerID)
	if err != nil {
		return &LoginResponse{
			Success: false,
			Error:   "Invalid credentials",
		}, nil
	}
	
	// Verify password
	hashedPassword := as.hashPassword(password)
	if player.Password != hashedPassword {
		return &LoginResponse{
			Success: false,
			Error:   "Invalid credentials",
		}, nil
	}
	
	// Generate session token
	token, err := as.generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}
	
	// Create session
	session := &Session{
		PlayerID:  player.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hour session
	}
	
	as.sessions[token] = session
	
	return &LoginResponse{
		Success: true,
		Token:   token,
		Player:  player,
	}, nil
}

func (as *AuthService) ValidateToken(token string) (*models.Player, error) {
	session, exists := as.sessions[token]
	if !exists {
		return nil, errors.New("invalid token")
	}
	
	// Check if session expired
	if time.Now().After(session.ExpiresAt) {
		delete(as.sessions, token)
		return nil, errors.New("session expired")
	}
	
	// Load player
	player, err := as.storage.LoadPlayer(session.PlayerID)
	if err != nil {
		return nil, fmt.Errorf("failed to load player: %v", err)
	}
	
	return player, nil
}

func (as *AuthService) Logout(token string) {
	delete(as.sessions, token)
}

func (as *AuthService) generatePlayerID(username string) string {
	// Simple player ID generation (username-based)
	// In production, you might want more sophisticated ID generation
	return fmt.Sprintf("player_%s", username)
}

// Middleware for token validation
func (as *AuthService) AuthMiddleware(next func(*models.Player)) func(string) error {
	return func(token string) error {
		player, err := as.ValidateToken(token)
		if err != nil {
			return err
		}
		next(player)
		return nil
	}
}

