// pkg/protocol/constants.go - Game constants
package protocol

const (
	// Message types
	MsgTypeLogin          = "login"
	MsgTypeLoginResponse  = "login_response"
	MsgTypeCreateGame     = "create_game"
	MsgTypeJoinGame       = "join_game"
	MsgTypeSimpleAction   = "simple_action"
	MsgTypeEnhancedAction = "enhanced_action"
	MsgTypeGameState      = "game_state"
	MsgTypeActionResult   = "action_result"
	MsgTypePlayerJoined   = "player_joined"
	MsgTypeGameEnd        = "game_end"
	MsgTypePing           = "ping"
	MsgTypePong           = "pong"
	
	// Game modes
	GameModeSimple   = "simple"
	GameModeEnhanced = "enhanced"
	
	// Tower types
	TowerTypeKing  = "king_tower"
	TowerTypeGuard = "guard_tower"
	
	// Action types
	ActionTypeAttack     = "attack"
	ActionTypeSpawnTroop = "spawn_troop"
	
	// Game states
	GameStateWaiting    = "waiting"
	GameStateInProgress = "in_progress"
	GameStateFinished   = "finished"
)

// Default game values
const (
	DefaultManaStart     = 5
	DefaultManaMax       = 10
	DefaultManaRegen     = 1.0
	DefaultCritMultiplier = 1.2
	DefaultExpWin = 30
	DefaultExpDraw = 10
	DefaultGameDuration = 180 // 3 minutes
	DefaultTurnTime = 30     // 30 seconds per turn
)