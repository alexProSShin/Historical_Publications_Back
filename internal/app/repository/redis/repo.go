package redisrepo

import (
	"backend/internal/app/repository"
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisRepository(client *redis.Client) repository.RedisRepo {
	return &RedisRepository{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisRepository) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *RedisRepository) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisRepository) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisRepository) Exists(key string) (bool, error) {
	count, err := r.client.Exists(r.ctx, key).Result()
	return count > 0, err
}
