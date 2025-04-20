package database

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisDB struct {
	Client *redis.Client
}

// NewRedisDB creates a new Redis client. Accepts either:
// 1. Direct connection parameters (addr, password, db)
// 2. Redis URL format (redis://user:password@host:port/db)
func NewRedisDB(connection string, password string, db int) (*RedisDB, error) {
	var client *redis.Client

	// If connection string is a Redis URL
	if strings.HasPrefix(connection, "redis://") {
		parsedURL, err := url.Parse(connection)
		if err != nil {
			return nil, fmt.Errorf("invalid Redis URL: %w", err)
		}

		// Extract password from URL if not explicitly provided
		if password == "" {
			var hasPassword bool
			password, hasPassword = parsedURL.User.Password()
			if !hasPassword {
				password = "" // If no password in URL
			}
		}

		client = redis.NewClient(&redis.Options{
			Addr:     parsedURL.Host,
			Password: password,
			DB:       db, // Default to 0 if not specified in URL
		})
	} else {
		// Standard connection parameters
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

	log.Println("✅ Successfully connected to Redis")
	return &RedisDB{Client: client}, nil
}

func (r *RedisDB) Close() error {
	if r.Client != nil {
		return r.Client.Close()
	}
	return nil
}
