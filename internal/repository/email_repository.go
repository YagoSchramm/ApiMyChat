package repository

import (
	"errors"
	"sync"
	"time"
)

type item struct {
	code      string
	expiresAt time.Time
}

type MemoryCache struct {
	mu    sync.RWMutex
	store map[string]item
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		store: make(map[string]item),
	}
}

func (m *MemoryCache) SaveOTP(email, code string, duration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store[email] = item{
		code:      code,
		expiresAt: time.Now().Add(duration),
	}
	return nil
}

func (m *MemoryCache) GetOTP(email string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	val, ok := m.store[email]
	if !ok || time.Now().After(val.expiresAt) {
		return "", errors.New("código expirado ou inexistente")
	}
	return val.code, nil
}
