package auth

import (
	"sync"
	"time"

	"github.com/tarro-dev/discord-oauth/pkg/option"
)

const (
	tickerDuration = 5 * time.Minute
	ttlDuration    = 30 * time.Minute
)

type store struct {
	sessions map[string]entry
	mu       sync.Mutex
	ticker   *time.Ticker
}

type entry struct {
	session Session
	ttl     time.Time
}

func newStore() *store {
	return &store{
		sessions: make(map[string]entry),
		ticker:   time.NewTicker(5 * time.Minute),
	}
}

func (s *store) get(id string) option.Option[Session] {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.sessions[id]
	if !ok {
		return option.None[Session]()
	}

	s.sessions[id] = entry{
		session: e.session,
		ttl:     time.Now().Add(ttlDuration),
	}

	return option.Some(e.session)
}

func (s *store) store(session Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[session.ID] = entry{
		session: session,
		ttl:     time.Now().Add(ttlDuration),
	}
}

func (s *store) delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, id)
}

func (s *store) garbageCollect() {
	for range s.ticker.C {
		s.cleanup()
	}
}

func (s *store) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for id, e := range s.sessions {
		if e.ttl.Before(now) {
			delete(s.sessions, id)
		}
	}
}
