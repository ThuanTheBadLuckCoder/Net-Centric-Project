// Fixed cmd/client/main.go - Handle the auth.LoginResponse correctly
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	
	"github.com/gorilla/websocket"
	"tcr-game/internal/auth"
)

type TestClient struct {
	serverURL string
	token     string
	gameID    string
	playerID  string
	wsConn    *websocket.Conn
}

func main() {
	serverURL := flag.String("server", "http://localhost:8080", "Server URL")
	flag.Parse()
	
	client := &TestClient{
		serverURL: *serverURL,
	}
	
	fmt.Println("TCR Game Test Client")
	fmt.Println("Commands: register, login, create, join, attack, state, ws, quit")
	fmt.Println("Example: login username password")
	
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		
		command := strings.TrimSpace(scanner.Text())
		if command == "quit" {
			break
		}
		
		parts := strings.Fields(command)
		if len(parts) == 0 {
			continue
		}
		
		switch parts[0] {
		case "register":
			if len(parts) != 3 {
				fmt.Println("Usage: register <username> <password>")
				continue
			}
			client.register(parts[1], parts[2])
		case "login":
			if len(parts) != 3 {
				fmt.Println("Usage: login <username> <password>")
				continue
			}
			client.login(parts[1], parts[2])
		case "create":
			if len(parts) != 3 {
				fmt.Println("Usage: create <mode> <gameID>")
				continue
			}
			client.createGame(parts[1], parts[2])
		case "join":
			if len(parts) != 2 {
				fmt.Println("Usage: join <gameID>")
				continue
			}
			client.joinGame(parts[1])
		case "attack":
			if len(parts) != 3 {
				fmt.Println("Usage: attack <troopID> <towerIndex>")
				continue
			}
			client.attack(parts[1], parts[2])
		case "state":
			client.getGameState()
		case "ws":
			client.connectWebSocket()
		default:
			fmt.Printf("Unknown command: %s\n", parts[0])
		}
	}
	
	if client.wsConn != nil {
		client.wsConn.Close()
	}
}

func (c *TestClient) register(username, password string) {
	creds := auth.Credentials{
		Username: username,
		Password: password,
	}
	
	data, _ := json.Marshal(creds)
	resp, err := http.Post(c.serverURL+"/api/register", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Register error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	
	if success, ok := result["success"].(bool); ok && success {
		fmt.Println("Registration successful")
	} else {
		fmt.Printf("Registration failed: %v\n", result["error"])
	}
}

func (c *TestClient) login(username, password string) {
	creds := auth.Credentials{
		Username: username,
		Password: password,
	}
	
	data, _ := json.Marshal(creds)
	resp, err := http.Post(c.serverURL+"/api/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Login error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	// Use map[string]interface{} instead of auth.LoginResponse to handle type assertion better
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	
	if success, ok := result["success"].(bool); ok && success {
		if token, ok := result["token"].(string); ok {
			c.token = token
		}
		
		// Handle player data as map[string]interface{}
		if playerData, ok := result["player"].(map[string]interface{}); ok {
			if id, ok := playerData["id"].(string); ok {
				c.playerID = id
			}
		}
		
		fmt.Printf("Login successful. Token: %s\n", c.token[:10]+"...")
	} else {
		if errorMsg, ok := result["error"].(string); ok {
			fmt.Printf("Login failed: %s\n", errorMsg)
		} else {
			fmt.Println("Login failed: unknown error")
		}
	}
}

func (c *TestClient) createGame(mode, gameID string) {
	if c.token == "" {
		fmt.Println("Please login first")
		return
	}
	
	gameData := map[string]string{
		"mode":    mode,
		"game_id": gameID,
	}
	
	data, _ := json.Marshal(gameData)
	req, _ := http.NewRequest("POST", c.serverURL+"/api/games", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.token)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Create game error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	
	if success, ok := result["success"].(bool); ok && success {
		c.gameID = gameID
		fmt.Printf("Game created: %s\n", gameID)
	} else {
		fmt.Printf("Create game failed: %v\n", result["error"])
	}
}

func (c *TestClient) joinGame(gameID string) {
	if c.token == "" {
		fmt.Println("Please login first")
		return
	}
	
	req, _ := http.NewRequest("POST", c.serverURL+"/api/games/"+gameID+"/join", nil)
	req.Header.Set("Authorization", c.token)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Join game error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	
	if success, ok := result["success"].(bool); ok && success {
		c.gameID = gameID
		fmt.Printf("Joined game: %s\n", gameID)
	} else {
		fmt.Printf("Join game failed: %v\n", result["error"])
	}
}

func (c *TestClient) attack(troopID, towerIndex string) {
	if c.token == "" || c.gameID == "" {
		fmt.Println("Please login and join a game first")
		return
	}
	
	// Convert towerIndex to int
	towerIdx := 0
	if towerIndex == "1" {
		towerIdx = 1
	} else if towerIndex == "2" {
		towerIdx = 2
	}
	
	actionData := map[string]interface{}{
		"type":        "attack",
		"troop_id":    troopID,
		"target_tower": towerIdx,
	}
	
	data, _ := json.Marshal(actionData)
	req, _ := http.NewRequest("POST", c.serverURL+"/api/games/"+c.gameID+"/action", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.token)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Attack error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	
	if success, ok := result["success"].(bool); ok && success {
		fmt.Println("Attack successful")
		if br, ok := result["battle_result"].(map[string]interface{}); ok {
			fmt.Printf("Damage: %v, Crit: %v, Tower destroyed: %v\n",
				br["damage"], br["critical_hit"], br["tower_destroyed"])
		}
	} else {
		fmt.Printf("Attack failed: %v\n", result["error"])
	}
}

func (c *TestClient) getGameState() {
	if c.token == "" || c.gameID == "" {
		fmt.Println("Please login and join a game first")
		return
	}
	
	req, _ := http.NewRequest("GET", c.serverURL+"/api/games/"+c.gameID+"/state", nil)
	req.Header.Set("Authorization", c.token)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Get state error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	var state map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&state)
	
	fmt.Println("Game State:")
	stateJSON, _ := json.MarshalIndent(state, "", "  ")
	fmt.Println(string(stateJSON))
}

func (c *TestClient) connectWebSocket() {
	if c.token == "" || c.gameID == "" {
		fmt.Println("Please login and join a game first")
		return
	}
	
	wsURL := "ws://localhost:8080/ws/" + c.gameID + "?token=" + c.token
	
	var err error
	c.wsConn, _, err = websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		fmt.Printf("WebSocket error: %v\n", err)
		return
	}
	
	fmt.Println("WebSocket connected")
	
	// Listen for messages
	go func() {
		for {
			_, message, err := c.wsConn.ReadMessage()
			if err != nil {
				fmt.Printf("WebSocket read error: %v\n", err)
				break
			}
			
			var msg map[string]interface{}
			json.Unmarshal(message, &msg)
			fmt.Printf("WebSocket: %s\n", string(message))
		}
	}()
	
	// Send ping periodically
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		
		for range ticker.C {
			if c.wsConn == nil {
				break
			}
			
			ping := map[string]string{"type": "ping"}
			c.wsConn.WriteJSON(ping)
		}
	}()
}