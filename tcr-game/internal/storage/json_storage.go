// internal/storage/json_storage.go - JSON storage implementation
package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	
	"tcr-game/internal/models"
)

type JSONStorage struct {
	playersDir string
	troopsFile string
	towersFile string
}

func NewJSONStorage(playersDir, troopsFile, towersFile string) *JSONStorage {
	return &JSONStorage{
		playersDir: playersDir,
		troopsFile: troopsFile,
		towersFile: towersFile,
	}
}

func (js *JSONStorage) LoadPlayer(id string) (*models.Player, error) {
	filename := filepath.Join(js.playersDir, fmt.Sprintf("%s.json", id))
	
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("player not found: %s", id)
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

func (js *JSONStorage) SavePlayer(player *models.Player) error {
	// Ensure directory exists
	if err := os.MkdirAll(js.playersDir, 0755); err != nil {
		return err
	}
	
	filename := filepath.Join(js.playersDir, fmt.Sprintf("%s.json", player.ID))
	
	data, err := json.MarshalIndent(player, "", "  ")
	if err != nil {
		return err
	}
	
	return ioutil.WriteFile(filename, data, 0644)
}

func (js *JSONStorage) LoadTroops() ([]*models.Troop, error) {
	data, err := ioutil.ReadFile(js.troopsFile)
	if err != nil {
		return nil, err
	}
	
	var troopSpecs []map[string]interface{}
	if err := json.Unmarshal(data, &troopSpecs); err != nil {
		return nil, err
	}
	
	troops := make([]*models.Troop, 0, len(troopSpecs))
	for _, spec := range troopSpecs {
		troop := &models.Troop{
			ID:          spec["id"].(string),
			Name:        spec["name"].(string),
			HP:          int(spec["hp"].(float64)),
			MaxHP:       int(spec["hp"].(float64)),
			Attack:      int(spec["attack"].(float64)),
			Defense:     int(spec["defense"].(float64)),
			CritChance:  spec["crit_chance"].(float64),
			ManaCost:    int(spec["mana_cost"].(float64)),
			Description: spec["description"].(string),
			Level:       1,
		}
		troops = append(troops, troop)
	}
	
	return troops, nil
}

func (js *JSONStorage) LoadTowers() (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(js.towersFile)
	if err != nil {
		return nil, err
	}
	
	var towers map[string]interface{}
	if err := json.Unmarshal(data, &towers); err != nil {
		return nil, err
	}
	
	return towers, nil
}