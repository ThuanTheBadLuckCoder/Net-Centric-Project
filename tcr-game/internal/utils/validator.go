// internal/utils/validator.go - Input validation utilities
package utils

import (
	"errors"
	"regexp"
	"strings"
)

func ValidateUsername(username string) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if len(username) > 20 {
		return errors.New("username must be at most 20 characters")
	}
	
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if !matched {
		return errors.New("username can only contain letters, numbers, and underscores")
	}
	
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if len(password) > 100 {
		return errors.New("password is too long")
	}
	
	return nil
}

func ValidateGameID(gameID string) error {
	gameID = strings.TrimSpace(gameID)
	if len(gameID) < 1 {
		return errors.New("game ID cannot be empty")
	}
	if len(gameID) > 50 {
		return errors.New("game ID is too long")
	}
	
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", gameID)
	if !matched {
		return errors.New("game ID can only contain letters, numbers, hyphens, and underscores")
	}
	
	return nil
}