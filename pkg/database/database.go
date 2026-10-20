package database

import (
	"context"
	"currencies/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateDbPool(config *config.DatabaseConfig) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), config.Url)
	/*
		if err != nil {
			panic(err)
		}
	*/
	return dbpool, err
}
