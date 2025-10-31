package service

import (
	"context"
	"fmt"
	"github.com/Guram-Gurych/shortenerURL.git/internal/repository"
	"github.com/google/uuid"
)

type URLShortener interface {
	CreateShortURL(ctx context.Context, originalURL string) (string, error)
	GetOriginalURL(ctx context.Context, id string) (string, error)
}

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

func (ss *ShortenerService) CreateShortURL(ctx context.Context, originalURL string) (string, error) {
	id := generateID()

	if err := ss.repo.Save(ctx, id, originalURL); err != nil {
		return "", fmt.Errorf("не удалось сохранить URL в сервисе: %w", err)
	}

	return id, nil
}

func (ss *ShortenerService) GetOriginalURL(ctx context.Context, id string) (string, error) {
	value, err := ss.repo.Get(ctx, id)

	if err != nil {
		return "", err
	}

	return value, nil
}
