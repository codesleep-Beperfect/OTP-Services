package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo(c *redis.Client) *RedisRepo {
	return &RedisRepo{client: c}
}

func (r *RedisRepo) Set(key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisRepo) Get(key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisRepo) Exists(key string) (bool, error) {
	val, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val == 1, nil
}

func (r *RedisRepo) Delete(key string) error {
	return r.client.Del(ctx, key).Err()
}