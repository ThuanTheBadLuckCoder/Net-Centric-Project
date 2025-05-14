// Updated internal/server/server.go - Serve index.html at root
package server

import (
	"context"
	"net/http"
	"time"
	
	"github.com/gorilla/mux"
	"tcr-game/config"
	"tcr-game/internal/auth"
	"tcr-game/internal/game"
	"tcr-game/internal/storage"
)

type Server struct {
	config      *config.Config
	router      *mux.Router
	httpServer  *http.Server
	gameEngine  *game.GameEngine
	authService *auth.AuthService
	wsManager   *WebSocketManager
}

func New(cfg *config.Config) *Server {
	// Initialize storage
	storage := storage.NewJSONStorage(
		cfg.Database.PlayersDir,
		cfg.Database.TroopsFile,
		cfg.Database.TowersFile,
	)
	
	// Initialize services
	authService := auth.NewAuthService(storage)
	gameEngine := game.NewGameEngine(cfg, storage)
	wsManager := NewWebSocketManager()
	
	s := &Server{
		config:      cfg,
		gameEngine:  gameEngine,
		authService: authService,
		wsManager:   wsManager,
	}
	
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router = mux.NewRouter()
	
	// Apply middleware
	s.router.Use(s.corsMiddleware)
	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.rateLimitMiddleware)
	
	// Auth routes
	s.router.HandleFunc("/api/register", s.handleRegister).Methods("POST")
	s.router.HandleFunc("/api/login", s.handleLogin).Methods("POST")
	s.router.HandleFunc("/api/logout", s.handleLogout).Methods("POST")
	
	// Game routes
	s.router.HandleFunc("/api/games", s.handleCreateGame).Methods("POST")
	s.router.HandleFunc("/api/games/{gameID}/join", s.handleJoinGame).Methods("POST")
	s.router.HandleFunc("/api/games/{gameID}/state", s.handleGetGameState).Methods("GET")
	s.router.HandleFunc("/api/games/{gameID}/action", s.handleGameAction).Methods("POST")
	
	// WebSocket route
	s.router.HandleFunc("/ws/{gameID}", s.handleWebSocket)
	
	// Serve index.html at root
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/index.html")
	}).Methods("GET")
	
	// Static files
	s.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("web/static/"))))
}

func (s *Server) Start(port string) error {
	s.httpServer = &http.Server{
		Addr:         ":" + port,
		Handler:      s.router,
		ReadTimeout:  time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Server.WriteTimeout) * time.Second,
	}
	
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}