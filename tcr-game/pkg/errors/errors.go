// pkg/errors/errors.go - Custom error types
package errors

import "fmt"

type GameError struct {
	Code    string
	Message string
}

func (e *GameError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Error codes
const (
	ErrCodeGameNotFound    = "GAME_NOT_FOUND"
	ErrCodeGameFull        = "GAME_FULL"
	ErrCodeInvalidAction   = "INVALID_ACTION"
	ErrCodeInsufficientMana = "INSUFFICIENT_MANA"
	ErrCodeNotYourTurn     = "NOT_YOUR_TURN"
	ErrCodeGameEnded       = "GAME_ENDED"
	ErrCodeInvalidTarget   = "INVALID_TARGET"
	ErrCodeTroopNotFound   = "TROOP_NOT_FOUND"
	ErrCodeUnauthorized    = "UNAUTHORIZED"
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
)

func NewGameError(code, message string) *GameError {
	return &GameError{Code: code, Message: message}
}

func GameNotFound(gameID string) *GameError {
	return NewGameError(ErrCodeGameNotFound, fmt.Sprintf("game %s not found", gameID))
}

func GameFull() *GameError {
	return NewGameError(ErrCodeGameFull, "game is full")
}

func InvalidAction(msg string) *GameError {
	return NewGameError(ErrCodeInvalidAction, msg)
}

func InsufficientMana(need, have int) *GameError {
	return NewGameError(ErrCodeInsufficientMana, fmt.Sprintf("need %d mana, have %d", need, have))
}

func NotYourTurn() *GameError {
	return NewGameError(ErrCodeNotYourTurn, "not your turn")
}

func GameEnded() *GameError {
	return NewGameError(ErrCodeGameEnded, "game has ended")
}

func InvalidTarget(msg string) *GameError {
	return NewGameError(ErrCodeInvalidTarget, msg)
}

func TroopNotFound(troopID string) *GameError {
	return NewGameError(ErrCodeTroopNotFound, fmt.Sprintf("troop %s not found", troopID))
}

func Unauthorized() *GameError {
	return NewGameError(ErrCodeUnauthorized, "unauthorized")
}

func InvalidCredentials() *GameError {
	return NewGameError(ErrCodeInvalidCredentials, "invalid credentials")
}