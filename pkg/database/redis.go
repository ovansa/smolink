package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisDB struct {
	Client *redis.Client
}

func NewRedisDB(addr, password string, db int) (*RedisDB, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("✅ Successfully connected to Redis")
	return &RedisDB{Client: client}, nil
}

func (r *RedisDB) Close() error {
	return r.Client.Close()
}
