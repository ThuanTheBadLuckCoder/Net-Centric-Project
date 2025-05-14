// internal/storage/player_storage.go - Player data persistence
package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	
	"tcr-game/internal/models"
)

type PlayerStorage struct {
	playersDir string
}

func NewPlayerStorage(playersDir string) *PlayerStorage {
	return &PlayerStorage{
		playersDir: playersDir,
	}
}

func (ps *PlayerStorage) SavePlayer(player *models.Player) error {
	// Ensure directory exists
	if err := os.MkdirAll(ps.playersDir, 0755); err != nil {
		return err
	}
	
	filename := filepath.Join(ps.playersDir, fmt.Sprintf("%s.json", player.ID))
	
	data, err := json.MarshalIndent(player, "", "  ")
	if err != nil {
		return err
	}
	
	return ioutil.WriteFile(filename, data, 0644)
}

func (ps *PlayerStorage) LoadPlayer(playerID string) (*models.Player, error) {
	filename := filepath.Join(ps.playersDir, fmt.Sprintf("%s.json", playerID))
	
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("player not found: %s", playerID)
	}
	
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	
	var player models.Player
	if err := json.Unmarshal(data, &player); err != nil {
		return nil, err
	}
	
	return &player, nil
}

func (ps *PlayerStorage) PlayerExists(playerID string) bool {
	filename := filepath.Join(ps.playersDir, fmt.Sprintf("%s.json", playerID))
	_, err := os.Stat(filename)
	return err == nil
}

func (ps *PlayerStorage) DeletePlayer(playerID string) error {
	filename := filepath.Join(ps.playersDir, fmt.Sprintf("%s.json", playerID))
	return os.Remove(filename)
}

func (ps *PlayerStorage) ListPlayers() ([]string, error) {
	files, err := ioutil.ReadDir(ps.playersDir)
	if err != nil {
		return nil, err
	}
	
	players := make([]string, 0)
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			playerID := file.Name()[:len(file.Name())-5] // Remove .json extension
			players = append(players, playerID)
		}
	}
	
	return players, nil
}
