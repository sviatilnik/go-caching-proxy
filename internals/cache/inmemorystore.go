package cache

import (
	"context"
	"sync"
	"time"
)

type cachedResponse struct {
	Response
	until int64
}

type InMemoryStore struct {
	store map[string]cachedResponse
	mu    sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		store: make(map[string]cachedResponse),
		mu:    sync.RWMutex{},
	}
}

func (s *InMemoryStore) Get(ctx context.Context, key string) *Response {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if v, ok := s.store[key]; ok {
		if v.until > time.Now().Unix() {
			return &v.Response
		}
	}

	return nil
}

func (s *InMemoryStore) Save(ctx context.Context, key string, value *Response, ttl int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Сохраняем копию значения: если где-то дальше в коде изменят исходный Response,
	// это не затронет данные в нашем хранилище.

	valueCopy := *value

	headers := make([]Header, len(value.Headers))
	copy(headers, value.Headers)

	valueCopy.Headers = headers

	s.store[key] = cachedResponse{
		valueCopy,
		time.Now().Unix() + int64(ttl),
	}
	return nil
}

func (s *InMemoryStore) Has(ctx context.Context, key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if v, ok := s.store[key]; ok {
		return v.until > time.Now().Unix()
	}

	return false
}

func (s *InMemoryStore) Remove(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.store, key)
	return nil
}
