package core_tg_dialog

import (
	"sync"
	"time"
)

type GetStatus int

const (
	NotFound GetStatus = iota
	Expired
	Found
)

type Store struct {
	mu   sync.Mutex
	ttl  time.Duration
	data map[int64]Session
}

func NewStore(ttl time.Duration) *Store {
	return &Store{
		ttl:  ttl,
		data: make(map[int64]Session),
	}
}

func (s *Store) Get(chatID int64) (Session, GetStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.data[chatID]
	if !ok {
		return Session{}, NotFound
	}

	if time.Since(session.UpdatedAt) > s.ttl {
		delete(s.data, chatID)
		return Session{}, Expired
	}

	return session, Found
}

func (s *Store) Set(chatID int64, session Session) {
	session.UpdatedAt = time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[chatID] = session
}

func (s *Store) Delete(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, chatID)
}
