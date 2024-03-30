package dbpool

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ffigari/stored-strings/internal/config"
)

func NewFromConfig(ctx context.Context, dbName string) (*pgxpool.Pool, error) {
	config, err := config.Get()
	if err != nil {
		return nil, err
	}

	dbPool, err := pgxpool.New(
		ctx,
		config.PostgresServerConnectionString+"/"+dbName,
	)
	if err != nil {
		return nil, err
	}

	if err = dbPool.Ping(ctx); err != nil {
		return nil, err
	}

	return dbPool, nil
}
