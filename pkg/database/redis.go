package database

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	Client *redis.Client
}

// NewRedisDB creates a new Redis client with secure connection handling
func NewRedisDB(connection string, password string, db int) (*RedisDB, error) {
	var client *redis.Client
	var safeConn string

	// Parse and sanitize connection details for logging
	if strings.HasPrefix(connection, "redis://") {
		parsedURL, err := url.Parse(connection)
		if err != nil {
			return nil, fmt.Errorf("invalid Redis URL: %w", err)
		}

		// Extract password from URL if not explicitly provided
		if password == "" {
			password, _ = parsedURL.User.Password()
		}

		// Create sanitized version for logging
		safeURL := *parsedURL
		if safeURL.User != nil {
			if _, hasPassword := safeURL.User.Password(); hasPassword {
				safeURL.User = url.User(safeURL.User.Username())
			}
		}
		safeConn = safeURL.String()

		client = redis.NewClient(&redis.Options{
			Addr:     parsedURL.Host,
			Password: password,
			DB:       db,
		})
	} else {
		// For direct host:port connections
		safeConn = connection
		client = redis.NewClient(&redis.Options{
			Addr:     connection,
			Password: password,
			DB:       db,
		})
	}

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// Secure logging
	log.Printf("✅ Redis connection established | Addr: %s | DB: %d | Pool: %d",
		safeConn,
		db,
		client.Options().PoolSize,
	)

	return &RedisDB{Client: client}, nil
}

func (r *RedisDB) Close() error {
	if r.Client != nil {
		if err := r.Client.Close(); err != nil {
			return fmt.Errorf("error closing Redis connection: %w", err)
		}
		log.Print("Redis connection closed")
	}
	return nil
}
