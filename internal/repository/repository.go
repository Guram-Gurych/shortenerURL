package repository

import (
	"context"
	"sync"
)

type MemoryRepository struct {
	urls map[string]string
	mu   sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		urls: make(map[string]string),
	}
}

func (rep *MemoryRepository) Save(_ context.Context, id, url string) error {
	rep.mu.Lock()
	defer rep.mu.Unlock()

	_, ok := rep.urls[id]
	if ok {
		return ErrorAlreadyExists
	}

	rep.urls[id] = url
	return nil
}

func (rep *MemoryRepository) Get(_ context.Context, id string) (string, error) {
	rep.mu.RLock()
	defer rep.mu.RUnlock()
	value, ok := rep.urls[id]
	if !ok {
		return "", ErrorNotFound
	}

	return value, nil
}
