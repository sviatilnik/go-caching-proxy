package cache

import (
	"context"
	"sync"
)

type InMemoryStore struct {
	store map[string]Response
	mu    sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		store: make(map[string]Response),
		mu:    sync.RWMutex{},
	}
}

func (s *InMemoryStore) Get(ctx context.Context, key string) *Response {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if v, ok := s.store[key]; ok {
		return &v
	}

	return nil
}

func (s *InMemoryStore) Save(ctx context.Context, key string, value *Response, ttl int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Сохраняем копию значения: если где-то дальше в коде изменят исходный Response,
	// это не затронет данные в нашем хранилище.

	s.store[key] = *value
	return nil
}

func (s *InMemoryStore) Has(ctx context.Context, key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.store[key]
	return ok
}

func (s *InMemoryStore) Remove(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.store, key)
	return nil
}
