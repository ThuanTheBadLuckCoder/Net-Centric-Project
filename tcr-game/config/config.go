// config/config.go - Configuration management
package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Game     GameConfig     `json:"game"`
	Database DatabaseConfig `json:"database"`
}

type ServerConfig struct {
	Port            string `json:"port"`
	ReadTimeout     int    `json:"read_timeout"`
	WriteTimeout    int    `json:"write_timeout"`
	MaxConnections  int    `json:"max_connections"`
}

type GameConfig struct {
	Simple   SimpleGameConfig   `json:"simple"`
	Enhanced EnhancedGameConfig `json:"enhanced"`
}

type SimpleGameConfig struct {
	MaxPlayers int `json:"max_players"`
	TurnTime   int `json:"turn_time_seconds"`
}

type EnhancedGameConfig struct {
	GameDuration  int     `json:"game_duration_seconds"`
	ManaStart     int     `json:"mana_start"`
	ManaMax       int     `json:"mana_max"`
	ManaRegen     float64 `json:"mana_regen_per_second"`
	CritMultiplier float64 `json:"crit_multiplier"`
	ExpWin        int     `json:"exp_win"`
	ExpDraw       int     `json:"exp_draw"`
}

type DatabaseConfig struct {
	TroopsFile  string `json:"troops_file"`
	TowersFile  string `json:"towers_file"`
	PlayersDir  string `json:"players_directory"`
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}