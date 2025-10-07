package repository

import (
	"fmt"
)

type MemoryRepository struct {
	urls map[string]string
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		urls: make(map[string]string),
	}
}

func (rep *MemoryRepository) Save(id, url string) error {
	rep.urls[id] = url
	return nil
}

func (rep *MemoryRepository) Get(id string) (string, error) {
	value, ok := rep.urls[id]
	if !ok {
		return "", fmt.Errorf("URL с таким id не найден")
	}

	return value, nil
}
