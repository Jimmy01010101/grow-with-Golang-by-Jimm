//DB Latihan Nect Kita pake DB Real
package store

import (
	"errors"
	"sync"

	"go-login/models"
)

var (
	ErrUserNotFound = errors.New("user tidak ditemukan")
	ErrUserExists   = errors.New("username sudah terdaftar")
)

type MemStore struct {
	mu     sync.RWMutex
	users  map[string]*models.User
	nextID int
}

func NewMemStore() *MemStore {
	return &MemStore{
		users:  make(map[string]*models.User),
		nextID: 1,
	}
}

func (s *MemStore) Create(u *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ada := s.users[u.Username]; ada {
		return ErrUserExists
	}
	u.ID = s.nextID
	s.nextID++
	s.users[u.Username] = u
	return nil
}

func (s *MemStore) GetByUsername(username string) (*models.User, error) {
	s.mu.RLock() 
	defer s.mu.RUnlock()

	u, ada := s.users[username]
	if !ada {
		return nil, ErrUserNotFound
	}
	return u, nil
}