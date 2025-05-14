// internal/storage/game_storage.go - Game data persistence
package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	
	"tcr-game/internal/models"
)

type GameStorage struct {
	gamesDir string
}

func NewGameStorage(gamesDir string) *GameStorage {
	return &GameStorage{
		gamesDir: gamesDir,
	}
}

func (gs *GameStorage) SaveGame(game *models.Game) error {
	// Ensure directory exists
	if err := os.MkdirAll(gs.gamesDir, 0755); err != nil {
		return err
	}
	
	filename := filepath.Join(gs.gamesDir, fmt.Sprintf("%s.json", game.ID))
	
	data, err := json.MarshalIndent(game, "", "  ")
	if err != nil {
		return err
	}
	
	return ioutil.WriteFile(filename, data, 0644)
}

func (gs *GameStorage) LoadGame(gameID string) (*models.Game, error) {
	filename := filepath.Join(gs.gamesDir, fmt.Sprintf("%s.json", gameID))
	
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("game not found: %s", gameID)
	}
	
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	
	var game models.Game
	if err := json.Unmarshal(data, &game); err != nil {
		return nil, err
	}
	
	return &game, nil
}

func (gs *GameStorage) DeleteGame(gameID string) error {
	filename := filepath.Join(gs.gamesDir, fmt.Sprintf("%s.json", gameID))
	return os.Remove(filename)
}
