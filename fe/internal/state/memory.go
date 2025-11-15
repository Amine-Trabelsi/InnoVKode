package state

import (
	"sync"
	"time"

	"github.com/escalopa/inno-vkode/internal/domain"
)

type Store interface {
	Get(chatID int64) (*domain.Session, bool)
	Save(session *domain.Session)
	Delete(chatID int64)
	All() []*domain.Session
}

type MemoryStore struct {
	now func() time.Time
	mu  sync.RWMutex
	db  map[int64]*domain.Session
}

func NewMemoryStore(now func() time.Time) *MemoryStore {
	return &MemoryStore{
		now: now,
		db:  make(map[int64]*domain.Session),
	}
}

func (s *MemoryStore) Get(chatID int64) (*domain.Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if sess, ok := s.db[chatID]; ok {
		return sess, true
	}
	return nil, false
}

func (s *MemoryStore) Save(session *domain.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session.LastActivity = s.now()
	s.db[session.ChatID] = session
}

func (s *MemoryStore) Delete(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.db, chatID)
}

func (s *MemoryStore) All() []*domain.Session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]*domain.Session, 0, len(s.db))
	for _, sess := range s.db {
		items = append(items, sess)
	}
	return items
}
