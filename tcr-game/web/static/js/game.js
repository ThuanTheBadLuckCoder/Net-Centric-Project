// Complete web/static/js/game.js
class TCRGame {
    constructor() {
        this.wsConnection = null;
        this.token = null;
        this.gameId = null;
        this.playerId = null;
        this.selectedTroop = null;
        this.selectedTower = null;
        this.gameState = null;
        
        this.initEventListeners();
        this.showScreen('login-screen');
    }

    initEventListeners() {
        // Login/Register
        document.getElementById('login-btn').addEventListener('click', () => this.login());
        document.getElementById('register-btn').addEventListener('click', () => this.register());
        
        // Lobby events
        document.getElementById('create-simple-btn').addEventListener('click', () => this.createGame('simple'));
        document.getElementById('create-enhanced-btn').addEventListener('click', () => this.createGame('enhanced'));
        document.getElementById('join-game-btn').addEventListener('click', () => this.joinGame());
        
        // Game actions
        document.getElementById('attack-btn').addEventListener('click', () => this.attack());
        document.getElementById('end-turn-btn').addEventListener('click', () => this.endTurn());
        document.getElementById('new-game-btn').addEventListener('click', () => this.newGame());
        
        // Enter key for login
        document.getElementById('password').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.login();
        });
        
        // Enter key for join game
        document.getElementById('game-id-input').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.joinGame();
        });
    }

    showScreen(screenId) {
        document.querySelectorAll('.screen').forEach(screen => {
            screen.classList.add('hidden');
        });
        document.getElementById(screenId).classList.remove('hidden');
    }

    showError(elementId, message) {
        document.getElementById(elementId).textContent = message;
        setTimeout(() => {
            document.getElementById(elementId).textContent = '';
        }, 5000);
    }

    showStatus(message, type = 'info') {
        const statusEl = document.getElementById('lobby-status');
        statusEl.textContent = message;
        statusEl.className = `status-message ${type}`;
        setTimeout(() => {
            statusEl.textContent = '';
            statusEl.className = 'status-message';
        }, 5000);
    }

    async login() {
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;

        if (!username || !password) {
            this.showError('login-error', 'Please enter username and password');
            return;
        }

        try {
            const response = await fetch('/api/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, password })
            });

            const data = await response.json();
            
            if (data.success) {
                this.token = data.token;
                this.playerId = data.player.id;
                this.showScreen('game-lobby');
                console.log('Login successful!');
            } else {
                this.showError('login-error', data.error || 'Login failed');
            }
        } catch (error) {
            console.error('Login error:', error);
            this.showError('login-error', 'Network error: ' + error.message);
        }
    }

    async register() {
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;

        if (!username || !password) {
            this.showError('login-error', 'Please enter username and password');
            return;
        }

        if (password.length < 6) {
            this.showError('login-error', 'Password must be at least 6 characters');
            return;
        }

        try {
            const response = await fetch('/api/register', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, password })
            });

            const data = await response.json();
            
            if (data.success) {
                this.showStatus('Registration successful! Please login.', 'success');
            } else {
                this.showError('login-error', data.error || 'Registration failed');
            }
        } catch (error) {
            console.error('Register error:', error);
            this.showError('login-error', 'Network error: ' + error.message);
        }
    }

    async createGame(mode) {
        const gameId = `game_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
        
        console.log('Creating game:', gameId, 'mode:', mode);
        
        try {
            const response = await fetch('/api/games', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': this.token
                },
                body: JSON.stringify({ mode, game_id: gameId })
            });

            const data = await response.json();
            console.log('Create game response:', data);
            
            if (data.success) {
                this.gameId = gameId;
                this.connectWebSocket();
                this.showScreen('game-screen');
                this.showStatus(`Game created: ${gameId}`, 'success');
            } else {
                this.showStatus(data.error || 'Failed to create game', 'error');
            }
        } catch (error) {
            console.error('Create game error:', error);
            this.showStatus('Network error: ' + error.message, 'error');
        }
    }

    async joinGame() {
        const gameId = document.getElementById('game-id-input').value.trim();
        
        if (!gameId) {
            this.showStatus('Please enter a game ID', 'error');
            return;
        }

        console.log('Joining game:', gameId);

        try {
            const response = await fetch(`/api/games/${gameId}/join`, {
                method: 'POST',
                headers: { 'Authorization': this.token }
            });

            const data = await response.json();
            console.log('Join game response:', data);
            
            if (data.success) {
                this.gameId = gameId;
                this.connectWebSocket();
                this.showScreen('game-screen');
            } else {
                this.showStatus(data.error || 'Failed to join game', 'error');
            }
        } catch (error) {
            console.error('Join game error:', error);
            this.showStatus('Network error: ' + error.message, 'error');
        }
    }

    connectWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws/${this.gameId}?token=${this.token}`;
        
        console.log('Connecting to WebSocket:', wsUrl);
        
        this.wsConnection = new WebSocket(wsUrl);
        
        this.wsConnection.onopen = () => {
            console.log('WebSocket connected');
            this.requestGameState();
        };
        
        this.wsConnection.onmessage = (event) => {
            const message = JSON.parse(event.data);
            console.log('WebSocket message:', message);
            this.handleWebSocketMessage(message);
        };
        
        this.wsConnection.onclose = () => {
            console.log('WebSocket disconnected');
            setTimeout(() => this.connectWebSocket(), 2000);
        };
        
        this.wsConnection.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    }

    handleWebSocketMessage(message) {
        switch (message.type) {
            case 'game_state':
                this.updateGameState(message.data);
                break;
            case 'turn_result':
            case 'action_result':
                this.handleActionResult(message.result);
                this.requestGameState();
                break;
            case 'player_joined':
                this.addLogEntry(`${message.data.username} joined the game`);
                break;
            case 'game_end':
                this.handleGameEnd(message);
                break;
            case 'pong':
                break;
            default:
                console.log('Unknown message type:', message.type);
        }
    }

    requestGameState() {
        if (this.wsConnection && this.wsConnection.readyState === WebSocket.OPEN) {
            this.wsConnection.send(JSON.stringify({ type: 'get_state' }));
        }
    }

    updateGameState(gameState) {
        this.gameState = gameState;
        console.log('Updating game state:', gameState);
        
        // Update game info
        document.getElementById('game-mode').textContent = `Mode: ${gameState.mode}`;
        
        if (gameState.mode === 'enhanced' && gameState.time_remaining !== undefined) {
            const minutes = Math.floor(gameState.time_remaining / 60);
            const seconds = gameState.time_remaining % 60;
            document.getElementById('game-time').textContent = 
                `Time: ${minutes}:${seconds.toString().padStart(2, '0')}`;
        }
        
        // Update players
        if (gameState.players && gameState.players.length >= 2) {
            const player = gameState.players.find(p => p.id === this.playerId);
            const opponent = gameState.players.find(p => p.id !== this.playerId);
            
            if (player && opponent) {
                document.getElementById('player-name').textContent = `You (${player.username})`;
                document.getElementById('opponent-name').textContent = opponent.username;
                
                this.updateTowers(player.towers, 'player');
                this.updateTowers(opponent.towers, 'opp');
                this.updateTroops(player.troops || player.available_troops);
                
                if (gameState.mode === 'enhanced' && player.mana !== undefined) {
                    document.getElementById('player-mana').textContent = 
                        `Mana: ${player.mana}/${player.max_mana}`;
                }
                
                if (gameState.mode === 'simple' && gameState.current_player) {
                    const isMyTurn = gameState.current_player.id === this.playerId;
                    document.getElementById('current-turn').textContent = 
                        `Turn: ${isMyTurn ? 'Your Turn' : gameState.current_player.username}`;
                    document.getElementById('attack-btn').disabled = !isMyTurn;
                    document.getElementById('end-turn-btn').disabled = !isMyTurn;
                }
            }
        }
    }

    updateTowers(towers, prefix) {
        towers.forEach((tower, index) => {
            const towerEl = document.getElementById(`${prefix}-tower-${index}`);
            if (towerEl) {
                const hpEl = towerEl.querySelector('.tower-hp');
                if (hpEl) {
                    hpEl.textContent = `${tower.hp}/${tower.max_hp}`;
                }
                
                if (!tower.alive) {
                    towerEl.classList.add('destroyed');
                }
                
                // Add click event for opponent towers (targets)
                if (prefix === 'opp') {
                    towerEl.addEventListener('click', () => this.selectTower(index));
                }
            }
        });
    }

    updateTroops(troops) {
        const troopsContainer = document.getElementById('troops-list');
        troopsContainer.innerHTML = '';
        
        if (!troops) return;
        
        troops.forEach(troop => {
            const troopCard = document.createElement('div');
            troopCard.className = 'troop-card';
            troopCard.dataset.troopId = troop.id;
            
            troopCard.innerHTML = `
                <div class="troop-name">${troop.name}</div>
                <div class="troop-stats">
                    HP: ${troop.hp}/${troop.max_hp}<br>
                    ATK: ${troop.attack} DEF: ${troop.defense}
                    ${troop.mana_cost ? `<br><span class="troop-mana">Mana: ${troop.mana_cost}</span>` : ''}
                </div>
            `;
            
            if (!troop.alive && troop.alive !== undefined) {
                troopCard.classList.add('used');
            } else {
                troopCard.addEventListener('click', () => this.selectTroop(troop.id));
            }
            
            troopsContainer.appendChild(troopCard);
        });
    }

    selectTroop(troopId) {
        // Remove previous selection
        document.querySelectorAll('.troop-card').forEach(card => {
            card.classList.remove('selected');
        });
        
        // Select new troop
        const troopCard = document.querySelector(`[data-troop-id="${troopId}"]`);
        if (troopCard) {
            troopCard.classList.add('selected');
            this.selectedTroop = troopId;
        }
        
        this.updateAttackButton();
    }

    selectTower(towerIndex) {
        // Remove previous selection
        document.querySelectorAll('.tower').forEach(tower => {
            tower.classList.remove('selected');
        });
        
        // Select new tower
        const towerEl = document.getElementById(`opp-tower-${towerIndex}`);
        if (towerEl) {
            towerEl.classList.add('selected');
            this.selectedTower = towerIndex;
        }
        
        this.updateAttackButton();
    }

    updateAttackButton() {
        const attackBtn = document.getElementById('attack-btn');
        const canAttack = this.selectedTroop && this.selectedTower !== null;
        
        if (this.gameState && this.gameState.mode === 'simple' && this.gameState.current_player) {
            const isMyTurn = this.gameState.current_player.id === this.playerId;
            attackBtn.disabled = !canAttack || !isMyTurn;
        } else {
            attackBtn.disabled = !canAttack;
        }
        
        if (canAttack) {
            attackBtn.textContent = `Attack with ${this.selectedTroop}`;
        } else {
            attackBtn.textContent = 'Select troop and target';
        }
    }

    async attack() {
        if (!this.selectedTroop || this.selectedTower === null) return;
        
        const actionData = {
            type: 'attack',
            troop_id: this.selectedTroop,
            target_tower: this.selectedTower
        };
        
        if (this.gameState.mode === 'enhanced') {
            actionData.type = 'spawn_troop';
            actionData.timestamp = new Date().toISOString();
        }
        
        console.log('Sending attack:', actionData);
        
        try {
            const response = await fetch(`/api/games/${this.gameId}/action`, {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': this.token
                },
                body: JSON.stringify(actionData)
            });
            
            const result = await response.json();
            console.log('Attack result:', result);
            this.handleActionResult(result);
        } catch (error) {
            console.error('Attack failed:', error);
            this.addLogEntry(`Attack failed: ${error.message}`, 'error');
        }
    }

    handleActionResult(result) {
        if (result.battle_result) {
            const br = result.battle_result;
            let logText = `${br.troop_used} attacks for ${br.damage} damage${br.critical_hit ? ' (CRIT!)' : ''}`;
            
            if (br.tower_destroyed) {
                logText += ' - Tower destroyed!';
            }
            
            this.addLogEntry(logText, br.critical_hit ? 'crit' : 'damage');
            
            if (br.tower_destroyed) {
                this.addLogEntry('Tower destroyed!', 'tower-destroyed');
            }
        }
        
        if (result.error) {
            this.addLogEntry(`Error: ${result.error}`, 'error');
        }
        
        // Clear selections after attack
        this.selectedTroop = null;
        this.selectedTower = null;
        document.querySelectorAll('.selected').forEach(el => {
            el.classList.remove('selected');
        });
        this.updateAttackButton();
    }

    addLogEntry(text, type = '') {
        const logContent = document.getElementById('log-content');
        const entry = document.createElement('div');
        entry.className = `log-entry ${type}`;
        entry.textContent = text;
        logContent.appendChild(entry);
        logContent.scrollTop = logContent.scrollHeight;
    }

    endTurn() {
        // In simple mode, ending turn without action
        // This could be implemented if needed
        console.log('End turn clicked');
    }

    handleGameEnd(message) {
        const resultEl = document.getElementById('game-result');
        const statsEl = document.getElementById('game-stats');
        
        if (message.winner === this.playerId) {
            resultEl.textContent = 'Victory!';
            resultEl.className = 'winner';
        } else if (message.winner) {
            resultEl.textContent = 'Defeat!';
            resultEl.className = 'loser';
        } else {
            resultEl.textContent = 'Draw!';
            resultEl.className = 'draw';
        }
        
        statsEl.innerHTML = `Reason: ${message.reason}`;
        this.showScreen('game-over');
    }

    newGame() {
        this.gameId = null;
        this.selectedTroop = null;
        this.selectedTower = null;
        this.gameState = null;
        
        if (this.wsConnection) {
            this.wsConnection.close();
        }
        
        this.showScreen('game-lobby');
    }
}

// Initialize game when page loads
document.addEventListener('DOMContentLoaded', () => {
    window.game = new TCRGame();
    console.log('Game initialized');
});

// Handle page visibility for WebSocket reconnection
document.addEventListener('visibilitychange', () => {
    if (!document.hidden && window.game && window.game.gameId) {
        // Reconnect WebSocket if needed
        if (!window.game.wsConnection || window.game.wsConnection.readyState !== WebSocket.OPEN) {
            window.game.connectWebSocket();
        }
    }
});