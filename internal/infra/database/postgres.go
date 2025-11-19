package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnection(connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return pgxpool.NewWithConfig(ctx, config)
}
