package repository

import "context"

type URLRepository interface {
	Save(ctx context.Context, id, url string) error
	Get(ctx context.Context, id string) (string, error)
}
