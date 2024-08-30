package sessions

import (
	"github.com/Zach51920/connect-four/internal/connectfour"
	"log/slog"
	"sync"
	"time"
)

const (
	MaxIdleTimeout = 1 * time.Minute
	PruneInterval  = 1 * time.Minute
)

type Store interface {
	New(id string, game *connectfour.Game) *Session
	Get(id string) (*Session, bool)
	Close()
}

type MemorySessionStore struct {
	sessionMu sync.RWMutex
	sessions  map[string]*Session

	shutdownOnce sync.Once
	shutdownCh   chan struct{}
}

func NewMemorySessionStore() *MemorySessionStore {
	store := &MemorySessionStore{
		sessions:   make(map[string]*Session),
		shutdownCh: make(chan struct{}),
	}
	go store.start()
	return store
}

func (s *MemorySessionStore) New(id string, game *connectfour.Game) *Session {
	session := &Session{ID: id, Game: game, LastUsed: time.Now()}
	s.sessionMu.Lock()
	defer s.sessionMu.Unlock()
	s.sessions[id] = session
	return session
}

func (s *MemorySessionStore) Get(id string) (*Session, bool) {
	s.sessionMu.RLock()
	defer s.sessionMu.RUnlock()
	sess, ok := s.sessions[id]
	return sess, ok
}

func (s *MemorySessionStore) Close() {
	s.shutdownOnce.Do(func() {
		close(s.shutdownCh)
	})
}

func (s *MemorySessionStore) start() {
	ticker := time.NewTicker(PruneInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.shutdownCh:
			return
		case <-ticker.C:
			s.prune()
		}
	}
}

func (s *MemorySessionStore) prune() {
	s.sessionMu.Lock()
	defer s.sessionMu.Unlock()

	for id, sess := range s.sessions {
		if time.Since(sess.LastUsed) > MaxIdleTimeout {
			slog.Debug("Removing stale session", "session_id", id)
			sess.CloseStream()
			delete(s.sessions, id)
		}
	}
}
