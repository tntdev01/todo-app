package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(connString string) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), connString)
}
