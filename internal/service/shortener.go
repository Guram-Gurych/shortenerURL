package service

import (
	"fmt"
	"github.com/Guram-Gurych/shortenerURL.git/internal/repository"
	"github.com/google/uuid"
)

type ShortenerService struct {
	repo repository.URLRepository
}

func NewShortenerService(repo repository.URLRepository) *ShortenerService {
	return &ShortenerService{
		repo: repo,
	}
}

func generateID() string {
	return uuid.New().String()[:8]
}

func (ss *ShortenerService) CreateShortURL(originalURL string) (string, error) {
	id := generateID()

	if err := ss.repo.Save(id, originalURL); err != nil {
		return "", fmt.Errorf("не удалось сохранить URL в сервисе: %w", err)
	}

	return id, nil
}

func (ss *ShortenerService) GetOriginalURL(id string) (string, error) {
	value, err := ss.repo.Get(id)

	if err != nil {
		return "", err
	}

	return value, nil
}
