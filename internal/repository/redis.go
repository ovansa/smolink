package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

func (r *RedisRepository) Client() *redis.Client {
	return r.client
}

func (r *RedisRepository) GetURL(ctx context.Context, shortCode string) (string, error) {
	return r.client.Get(ctx, "url:"+shortCode).Result()
}

func (r *RedisRepository) SetURL(ctx context.Context, shortCode, originalURL string, expiry time.Duration) error {
	return r.client.Set(ctx, "url:"+shortCode, originalURL, expiry).Err()
}
