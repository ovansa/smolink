package database

import (
	"context"
	"fmt"
	"log"
	"net/url"
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

	log.Printf("✅ Successfully connected to PostgreSQL at %s", sanitizeDSN(dsn))
	return &PostgresDB{Pool: pool}, nil
}

func (db *PostgresDB) Close() {
	db.Pool.Close()
}

func sanitizeDSN(dsn string) string {
	parsed, err := url.Parse(dsn)
	if err != nil {
		return "unknown DSN (failed to parse)"
	}

	if parsed.User != nil {
		username := parsed.User.Username()
		parsed.User = url.UserPassword(username, "****")
	}

	return parsed.Redacted()
}
