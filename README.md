# Tower Castle Rush (TCR) - Net-Centric Documentation

## Table of Contents
1. [Project Overview](#project-overview)
2. [Network Architecture](#network-architecture)
3. [Communication Protocols](#communication-protocols)
4. [Project Structure](#project-structure)
5. [Game Mechanics](#game-mechanics)
6. [Installation & Setup](#installation--setup)
7. [Network Implementation](#network-implementation)

## Project Overview

Tower Castle Rush is a multiplayer tower defense game demonstrating client-server networking concepts. The project implements both TCP and WebSocket communication protocols to support desktop and web clients.

### Key Network Features
- **TCP Socket Server** (Go) for desktop clients
- **WebSocket Support** for web browsers
- **Concurrent Client Handling** using Go goroutines
- **Real-time Communication** for game state synchronization
- **Turn-based and Real-time Protocols**

### Technology Stack
- **Server**: Go with TCP/WebSocket listeners
- **Desktop Client**: Go TCP client
- **Web Client**: HTML/JavaScript with WebSocket
- **Data Format**: JSON for all network messages

## Network Architecture

### Client-Server Model
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Client    │    │  Desktop Client │    │                 │
│ (WebSocket)     │◄──►│   (TCP)         │◄──►│   Go Server     │
│ Port: 8080/ws   │    │  Port: 8080     │    │   Port: 8080    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                       │
                               ┌───────────────────────▼────────────────────────┐
                               │  Concurrent Game Engine                        │
                               │  ┌────────────┐ ┌────────────┐ ┌─────────┐     │
                               │  │ Game 1     │ │ Game 2     │ │ Game N  │     │
                               │  │ P1 ◄──► P2 │ │ P3 ◄──► P4 │ │ ... ... │     │
                               │  └────────────┘ └────────────┘ └─────────┘     │
                               └────────────────────────────────────────────────┘
```

### Concurrency Flow Diagram
```
Main Server Thread
       │
       ▼
┌─────────────────┐
│ Accept          │
│ Connection      │
└─────┬───────────┘
      │
      ▼
┌─────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│ Spawn Goroutine │────►│ Authentication   │────►│ Game Matchmaking │
│ for Client      │     │ Thread           │     │ Thread           │
└─────────────────┘     └──────────────────┘     └──────────────────┘
                                                          │
                                                          ▼
                                                 ┌──────────────────┐
                                                 │ Game Instance    │
                                                 │ Goroutine        │
                                                 │ ┌──────────────┐ │
                                                 │ │ Player 1     │ │
                                                 │ │ Handler      │ │
                                                 │ └──────────────┘ │
                                                 │ ┌──────────────┐ │
                                                 │ │ Player 2     │ │
                                                 │ │ Handler      │ │
                                                 │ └──────────────┘ │
                                                 └──────────────────┘
```

## Communication Protocols

### TCP Protocol Flow (Desktop Client)

#### Connection Sequence Diagram
```
Client                      Server
  │                          │
  │── TCP Connect :8080 ────►│
  │◄── Welcome Message ──────│
  │                          │
  │── "alice mypassword" ───►│
  │◄── "SUCCESS: Welcome" ───│
  │                          │
  │── "SIMPLE" ─────────────►│
  │◄── "MODE_SELECTED" ──────│
  │                          │
  │◄── "WAITING..." ─────────│
  │◄── "GAME_START" ─────────│
  │                          │
  │── "DEPLOY Pawn" ────────►│
  │◄── "Deployed Pawn_L1" ───│
  │◄── "YOUR_TURN|..." ──────│
```

#### Message Flow State Machine
```
[CONNECTED] ──authentication──► [AUTHENTICATED] ──mode_select──► [MODE_SELECTED]
                                                                      │
                                                                      ▼
[GAME_OVER] ◄──game_end──── [IN_GAME] ◄──opponent_found──── [WAITING_MATCH]
     │                        │    ▲
     └──new_game──────────────┘    │
                                   │
                           ┌───────┴────────┐
                           │ Handle Commands│
                           │ - DEPLOY       │
                           │ - SURRENDER    │
                           │ - help         │
                           └────────────────┘
```

### WebSocket Protocol Flow (Web Client)

#### Message Exchange Simulation
```
Browser                                Server
   │                                     │
   │──── WebSocket Handshake ───────────►│
   │◄─── 101 Switching Protocols ────────│
   │                                     │
   │──── {"type":"createGame"} ─────────►│
   │◄─── {"type":"gameCreated"} ─────────│
   │                                     │
   │◄─── {"type":"playerJoined"} ────────│ ◄── Another player joins
   │                                     │
   │──── {"type":"deployTroop"} ────────►│
   │◄─── {"type":"gameStateUpdate"} ─────│
   │                                     │
   └──── Real-time State Sync ──────────►│
   ◄──── Every 2 seconds ────────────────┘
```

#### JSON Message Structure Visualization
```
┌─────────────────────────────────────────┐
│ Client → Server Messages                │
├─────────────────────────────────────────┤
│ ┌─────────────────────────────────────┐ │
│ │ CREATE GAME                         │ │
│ │ {                                   │ │
│ │   "type": "createGame",             │ │
│ │   "mode": "SIMPLE|ENHANCED",        │ │
│ │   "playerName": "string"            │ │
│ │ }                                   │ │
│ └─────────────────────────────────────┘ │
│                                         │
│ ┌─────────────────────────────────────┐ │
│ │ DEPLOY TROOP                        │ │
│ │ {                                   │ │
│ │   "type": "deployTroop",            │ │
│ │   "gameId": "abc123",               │ │
│ │   "troopType": "Pawn|Knight|..."    │ │
│ │ }                                   │ │
│ └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│ Server → Client Messages                │
├─────────────────────────────────────────┤
│ ┌─────────────────────────────────────┐ │
│ │ GAME STATE UPDATE                   │ │
│ │ {                                   │ │
│ │   "type": "gameStateUpdate",        │ │
│ │   "gameState": {                    │ │
│ │     "players": [...],               │ │
│ │     "turn": 1,                      │ │
│ │     "battleLog": [...]              │ │
│ │   }                                 │ │
│ │ }                                   │ │
│ └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

## Project Structure

### File Organization Diagram
```
tcr-game/
├── cmd/
│   ├── server/
│   │   └── main.go ──────────► TCP & WebSocket Server
│   └── client/
│       └── main.go ──────────► Desktop TCP Client
│
├── game/
│   ├── game.go ─────────────► Game Logic & State Management
│   ├── types.go ────────────► Data Structures
│   └── util.go ─────────────► File I/O Operations
│
├── web/templates/
│   ├── index.html ──────────► Web Interface Entry
│   └── game.html ───────────► Game Board UI
│
└── Data Files
    ├── players.json ────────► Player Database
    ├── towers.json ─────────► Tower Configuration
    └── troops.json ─────────► Unit Definitions
```

### Component Interaction Flow
```
┌─────────────┐    reads     ┌─────────────┐    loads     ┌─────────────┐
│ Server      │─────────────►│ JSON Files  │◄─────────────│ Game Engine │
│ main.go     │              │ (.json)     │              │ game.go     │
└─────────────┘              └─────────────┘              └─────────────┘
       │                                                          │
       │ creates                                                  │ manages
       ▼                                                          ▼
┌─────────────┐                                          ┌─────────────┐
│ Client      │                                          │ Game State  │
│ Connections │                                          │ & Logic     │
└─────────────┘                                          └─────────────┘
       │                                                          │
       │ communicates                                             │ updates
       ▼                                                          ▼
┌─────────────┐    ┌─────────────┐                      ┌─────────────┐
│ Desktop     │    │ Web         │                      │ Player Data │
│ Client      │    │ Client      │                      │ & Progress  │
└─────────────┘    └─────────────┘                      └─────────────┘
```

## Game Mechanics

### Game Mode Flow Comparison

#### Simple Mode (Turn-based) Flow
```
Game Start
    │
    ▼
┌─────────────┐    wait for    ┌─────────────┐    validate    ┌─────────────┐
│ Player 1    │   command      │ Receive     │   & execute    │ Process     │
│ Turn        │───────────────►│ DEPLOY      │───────────────►│ Combat      │
└─────────────┘                └─────────────┘                └─────────────┘
    ▲                                                                │
    │                                                                │
    │ switch turns                                                   │
    │                                                                ▼
┌─────────────┐    wait for    ┌─────────────┐    validate    ┌─────────────┐
│ Player 2    │   command      │ Receive     │   & execute    │ Check Win   │
│ Turn        │◄───────────────│ DEPLOY      │◄───────────────│ Condition   │
└─────────────┘                └─────────────┘                └─────────────┘
                                                                      │
                                                                      ▼
                                                               [GAME_END or CONTINUE]
```

#### Enhanced Mode (Real-time) Flow
```
Game Start (3 minute timer)
         │
         ▼
┌─────────────────────────────────────────────────────────────┐
│                Real-time Event Loop                         │
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │ Mana Regen  │  │ Auto Deploy │  │ Combat      │          │
│  │ Every 1s    │  │ Every 15s   │  │ Every 2s    │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
│         │                │                │                 │
│         ▼                ▼                ▼                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │ +1 Mana     │  │ Smart Troop │  │ Troops ──►  │          │
│  │ (max 10)    │  │ Selection   │  │ Towers      │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
│                                           │                 │
│         ┌─────────────────────────────────┘                 │
│         │ async                                             │
│         ▼                                                   │
│  ┌─────────────┐                                            │
│  │ Player      │                                            │
│  │ Commands    │                                            │
│  │ (anytime)   │                                            │
│  └─────────────┘                                            │
└─────────────────────────────────────────────────────────────┘
         │
         ▼
[Time Expires or King Tower Destroyed]
```

### Combat Resolution Simulation
```
Turn Combat Phase:

Troops ────attack───► Towers
  │                     │
  │ (damage calc)       │ (damage calc)
  │                     │
  ▼                     ▼
┌─────────────┐       ┌─────────────┐
│ ATK - DEF   │       │ ATK - DEF   │
│ = Damage    │       │ = Damage    │
└─────────────┘       └─────────────┘
  │                     │
  │ apply               │ apply
  ▼                     ▼
┌─────────────┐       ┌─────────────┐
│ Update      │       │ Update      │
│ Tower HP    │       │ Troop HP    │
└─────────────┘       └─────────────┘
  │                     │
  │ check               │ check
  ▼                     ▼
[Destroyed?]          [Destroyed?]
  │                     │
  ▼                     ▼
[Remove/Continue]     [Remove/Continue]
```

## Installation & Setup

### Setup Flow Diagram
```
Prerequisites Check
         │
         ▼
┌─────────────────┐
│ Go 1.16+        │
│ installed?      │
└─────┬───────────┘
      │ Yes
      ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Download/Clone  │───►│ Navigate to     │───►│ Start Server    │
│ Project         │    │ Directory       │    │ go run cmd/...  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                       │
                                                       ▼
                                              ┌─────────────────┐
                                              │ Server Running  │
                                              │ Port 8080       │
                                              └─────────────────┘
                                                       │
                        ┌──────────────────────────────┼──────────────────────────────┐
                        │                              │                              │
                        ▼                              ▼                              ▼
                ┌─────────────────┐            ┌─────────────────┐          ┌─────────────────┐
                │ Start Desktop   │            │ Open Web        │          │ Multiple        │
                │ Client          │            │ Browser         │          │ Clients...      │
                │ go run cmd/...  │            │ index.html      │          │                 │
                └─────────────────┘            └─────────────────┘          └─────────────────┘
```

### Configuration Structure
```
Server Configuration
         │
         ├── Network Settings
         │   ├── Port: 8080
         │   ├── Protocol: TCP + WebSocket
         │   └── Max Connections: Unlimited
         │
         ├── Data Files
         │   ├── players.json    (User accounts)
         │   ├── towers.json     (Game balance)
         │   └── troops.json     (Unit stats)
         │
         └── Game Settings
             ├── Enhanced Mode: 3min duration
             ├── Mana: Start 5, Max 10, +1/sec
             └── Combat: Every 2 seconds
```

## Network Implementation

### Server Architecture Simulation

#### Multi-Client Connection Handling
```
Main Server Process
         │
         ▼
┌─────────────────┐
│ Listen on       │
│ Port 8080       │
└─────┬───────────┘
      │
      ▼ (for each connection)
┌─────────────────┐
│ Accept()        │
│ New Connection  │
└─────┬───────────┘
      │
      ▼
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│ Goroutine 1     │     │ Goroutine 2     │     │ Goroutine N     │
│ ┌─────────────┐ │     │ ┌─────────────┐ │     │ ┌─────────────┐ │
│ │ Client A    │ │     │ │ Client B    │ │     │ │ Client Z    │ │
│ │ Handler     │ │     │ │ Handler     │ │     │ │ Handler     │ │
│ └─────────────┘ │     │ └─────────────┘ │     │ └─────────────┘ │
└─────────────────┘     └─────────────────┘     └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                                 ▼
                        ┌─────────────────┐
                        │ Shared Game     │
                        │ State           │
                        │ (Mutex Locked)  │
                        └─────────────────┘
```

#### Game State Synchronization Pattern
```
Player Action Trigger
         │
         ▼
┌─────────────────┐
│ Acquire Mutex   │
│ Lock            │
└─────┬───────────┘
      │
      ▼
┌─────────────────┐
│ Validate        │
│ Action          │
└─────┬───────────┘
      │ Valid
      ▼
┌─────────────────┐
│ Update Game     │
│ State           │
└─────┬───────────┘
      │
      ▼
┌─────────────────┐
│ Release Mutex   │
│ Lock            │
└─────┬───────────┘
      │
      ▼
┌─────────────────┐     ┌─────────────────┐
│ Broadcast to    │────►│ All Connected   │
│ All Players     │     │ Clients         │
└─────────────────┘     └─────────────────┘
```

### Client Connection Patterns

#### TCP Client Communication Loop
```
Connect to Server
         │
         ▼
┌─────────────────┐
│ Authentication  │
│ Phase           │
└─────┬───────────┘
      │ Success
      ▼
┌─────────────────┐     ┌─────────────────┐
│ Background      │     │ Main Input      │
│ Message Reader  │     │ Loop            │
│ Goroutine       │     │ Thread          │
└─────┬───────────┘     └─────┬───────────┘
      │                       │
      │ reads                 │ writes
      ▼                       ▼
┌─────────────────┐     ┌─────────────────┐
│ Server          │◄────│ User            │
│ Messages        │────►│ Commands        │
└─────────────────┘     └─────────────────┘
      │                       │
      ▼                       ▼
┌─────────────────┐     ┌─────────────────┐
│ Display to      │     │ Send to         │
│ Console         │     │ Server          │
└─────────────────┘     └─────────────────┘
```

#### WebSocket Client State Management
```
WebSocket Connection
         │
         ▼
┌─────────────────┐
│ Connection      │
│ Event Handlers  │
└─────┬───────────┘
      │
      ├── onopen ──────────► Connection Established
      │                       │
      ├── onmessage ────────► ┌─────────────────┐
      │                       │ Parse JSON      │
      │                       │ Update UI       │
      │                       │ Handle Events   │
      │                       └─────────────────┘
      │
      ├── onclose ─────────► ┌─────────────────┐
      │                      │ Reconnection    │
      │                      │ Logic           │
      │                      └─────────────────┘
      │
      └── onerror ─────────► ┌─────────────────┐
                             │ Error Display   │
                             │ & Recovery      │
                             └─────────────────┘
```

### Error Handling and Recovery

#### Connection Loss Recovery Simulation
```
Normal Operation
         │
         ▼
┌─────────────────┐
│ Client Active   │
│ Communication   │
└─────┬───────────┘
      │
      ▼ (connection lost)
┌─────────────────┐
│ Detect          │
│ Disconnection   │
└─────┬───────────┘
      │
      ▼
┌─────────────────┐    Retry 1    ┌─────────────────┐
│ Attempt         │──────────────►│ Wait 3 seconds  │
│ Reconnection    │◄──────────────│ Exponential     │
└─────┬───────────┘    Failed     │ Backoff         │
      │                           └─────────────────┘
      ▼ Success
┌─────────────────┐
│ Resume          │
│ Communication   │
└─────────────────┘
```

