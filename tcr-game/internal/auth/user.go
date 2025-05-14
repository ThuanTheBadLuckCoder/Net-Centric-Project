// internal/auth/user.go - User management
package auth

import (
	"tcr-game/internal/models"
	"tcr-game/internal/storage"
)

type UserManager struct {
	storage *storage.JSONStorage
}

func NewUserManager(storage *storage.JSONStorage) *UserManager {
	return &UserManager{
		storage: storage,
	}
}

func (um *UserManager) UpdatePlayerStats(player *models.Player, won bool, drawn bool) error {
	player.Stats.GamesPlayed++
	
	if won {
		player.Stats.GamesWon++
	} else if drawn {
		player.Stats.GamesDrawn++
	} else {
		player.Stats.GamesLost++
	}
	
	return um.storage.SavePlayer(player)
}

func (um *UserManager) GetPlayerProfile(playerID string) (*models.Player, error) {
	return um.storage.LoadPlayer(playerID)
}

func (um *UserManager) UpdatePlayerExperience(player *models.Player, exp int) error {
	oldLevel := player.Level
	player.AddExperience(exp)
	
	// Check if player leveled up
	if player.Level > oldLevel {
		// When player levels up, they can upgrade a troop or tower
		// This is just a placeholder - you might want to implement
		// a system where players choose what to upgrade
	}
	
	return um.storage.SavePlayer(player)
}

func (um *UserManager) LoadAvailableTroops(player *models.Player) error {
	// Load all troop templates
	troops, err := um.storage.LoadTroops()
	if err != nil {
		return err
	}
	
	// Apply player's troop levels
	player.AvailableTroops = make([]*models.Troop, len(troops))
	for i, troop := range troops {
		// Create a copy of the troop
		playerTroop := *troop
		playerTroop.ApplyLevel(player.GetTroopLevel(troop.ID))
		player.AvailableTroops[i] = &playerTroop
	}
	
	return nil
}