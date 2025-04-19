package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	Pool *pgxpool.Pool
}

func NewPostgresDB(dsn string) (*PostgresDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	log.Println("✅ Successfully connected to PostgreSQL")
	return &PostgresDB{Pool: pool}, nil
}

func (db *PostgresDB) Close() {
	db.Pool.Close()
}
