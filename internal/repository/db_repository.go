package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgconn"
)

type DBRepository struct {
	db *sql.DB
}

func NewDBRepository(db *sql.DB) *DBRepository {
	return &DBRepository{db: db}
}

func (db *DBRepository) Save(ctx context.Context, id, url string) error {
	query := "INSERT INTO urls (short_id, original_url) VALUES ($1, $2)"

	_, err := db.db.ExecContext(ctx, query, id, url)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrorAlreadyExists
			}
		}
		return err
	}

	return nil
}

func (db *DBRepository) Get(ctx context.Context, id string) (string, error) {
	var originalURL string
	query := "SELECT original_url FROM urls WHERE short_id = $1"

	err := db.db.QueryRowContext(ctx, query, id).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrorNotFound
		}
		return "", err
	}

	return originalURL, nil
}
