$folders = @(
    "tcr-game/configs",
    "tcr-game/data/players",
    "tcr-game/internal/auth",
    "tcr-game/internal/game",
    "tcr-game/internal/models",
    "tcr-game/internal/server",
    "tcr-game/internal/storage",
    "tcr-game/internal/utils",
    "tcr-game/cmd/server",
    "tcr-game/cmd/client",
    "tcr-game/pkg/protocol",
    "tcr-game/pkg/errors",
    "tcr-game/web/static/css",
    "tcr-game/web/static/js",
    "tcr-game/web/static/images",
    "tcr-game/web/templates",
    "tcr-game/tests/unit",
    "tcr-game/tests/integration",
    "tcr-game/tests/testdata",
    "tcr-game/scripts"
)

$files = @(
    "tcr-game/main.go",
    "tcr-game/go.mod",
    "tcr-game/go.sum",
    "tcr-game/README.md",
    "tcr-game/.gitignore",
    "tcr-game/configs/config.go",
    "tcr-game/configs/game_config.json",
    "tcr-game/data/troops.json",
    "tcr-game/data/towers.json",
    "tcr-game/data/players/player1.json",
    "tcr-game/data/players/player2.json",
    "tcr-game/internal/auth/auth.go",
    "tcr-game/internal/auth/user.go",
    "tcr-game/internal/game/game.go",
    "tcr-game/internal/game/engine.go",
    "tcr-game/internal/game/simple_rules.go",
    "tcr-game/internal/game/enhanced_rules.go",
    "tcr-game/internal/game/battle.go",
    "tcr-game/internal/game/events.go",
    "tcr-game/internal/models/player.go",
    "tcr-game/internal/models/tower.go",
    "tcr-game/internal/models/troop.go",
    "tcr-game/internal/models/battle.go",
    "tcr-game/internal/models/gamestate.go",
    "tcr-game/internal/server/server.go",
    "tcr-game/internal/server/handlers.go",
    "tcr-game/internal/server/websocket.go",
    "tcr-game/internal/server/middleware.go",
    "tcr-game/internal/storage/json_storage.go",
    "tcr-game/internal/storage/player_storage.go",
    "tcr-game/internal/storage/game_storage.go",
    "tcr-game/internal/utils/calculator.go",
    "tcr-game/internal/utils/validator.go",
    "tcr-game/internal/utils/logger.go",
    "tcr-game/cmd/server/main.go",
    "tcr-game/cmd/client/main.go",
    "tcr-game/pkg/protocol/messages.go",
    "tcr-game/pkg/protocol/constants.go",
    "tcr-game/pkg/errors/errors.go",
    "tcr-game/web/templates/index.html",
    "tcr-game/web/templates/game.html",
    "tcr-game/tests/unit/game_test.go",
    "tcr-game/tests/unit/battle_test.go",
    "tcr-game/tests/unit/models_test.go",
    "tcr-game/tests/integration/server_test.go",
    "tcr-game/tests/integration/game_flow_test.go",
    "tcr-game/tests/testdata/test_troops.json",
    "tcr-game/tests/testdata/test_players.json",
    "tcr-game/scripts/build.sh",
    "tcr-game/scripts/test.sh",
    "tcr-game/scripts/setup.sh"
)

foreach ($folder in $folders) {
    New-Item -Path $folder -ItemType Directory -Force | Out-Null
}

foreach ($file in $files) {
    New-Item -Path $file -ItemType File -Force | Out-Null
}

Write-Host "Project structure created successfully!"
