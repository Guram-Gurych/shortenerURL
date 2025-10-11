package repository

import (
	"fmt"
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

func (rep *MemoryRepository) Save(id, url string) error {
	rep.mu.Lock()
	defer rep.mu.Unlock()
	rep.urls[id] = url
	return nil
}

func (rep *MemoryRepository) Get(id string) (string, error) {
	rep.mu.Lock()
	defer rep.mu.Unlock()
	value, ok := rep.urls[id]
	if !ok {
		return "", fmt.Errorf("URL с таким id не найден")
	}

	return value, nil
}
