package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(connString string) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), connString)
}
