package handlers

import (
	"github.com/Zach51920/connect-four/internal/connectfour"
	"sync"
	"time"
)

type Sessions struct {
	gamesMu sync.RWMutex
	games   map[string]*connectfour.Game

	shutdownOnce sync.Once
	shutdownCh   chan struct{}
}

func NewSessions() *Sessions {
	s := &Sessions{
		games:      make(map[string]*connectfour.Game),
		shutdownCh: make(chan struct{}),
	}
	go s.start()
	return s
}

func (s *Sessions) Get(sessionID string) *connectfour.Game {
	s.gamesMu.RLock()
	defer s.gamesMu.RUnlock()
	game := s.games[sessionID]
	return game
}

func (s *Sessions) Set(sessionID string, game *connectfour.Game) {
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()
	s.games[sessionID] = game
}

func (s *Sessions) Shutdown() {
	s.shutdownOnce.Do(func() {
		close(s.shutdownCh)
	})
}

func (s *Sessions) start() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.shutdownCh:
			return
		case <-ticker.C:
			s.pruneStaleGames()
		}
	}
}

func (s *Sessions) pruneStaleGames() {
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()
	for sessionID, game := range s.games {
		if time.Since(game.Meta.LastMove) > 10*time.Minute || time.Since(game.Meta.StartTime) > time.Hour {
			delete(s.games, sessionID)
		}
	}
}
