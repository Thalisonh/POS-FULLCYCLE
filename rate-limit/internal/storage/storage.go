package storage

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type StorageInterface interface {
	Increment(ctx context.Context, key string) (int, error)
	Get(ctx context.Context, key string) (int, error)
	SetNX(ctx context.Context, key string, value int, expiration time.Duration) (bool, error)
	FlushAll() error
}

type RedisStorage struct {
	Client *redis.Client
}

func NewRedisStorage(addr string) StorageInterface {
	return &RedisStorage{
		Client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (r *RedisStorage) Increment(ctx context.Context, key string) (int, error) {
	count, err := r.Client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	r.Client.Expire(ctx, key, time.Second)
	return int(count), nil
}

func (r *RedisStorage) Get(ctx context.Context, key string) (int, error) {
	val, err := r.Client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}

	return val, err
}

func (r *RedisStorage) SetNX(ctx context.Context, key string, value int, expiration time.Duration) (bool, error) {
	set, err := r.Client.SetNX(ctx, key, value, expiration).Result()
	if err != nil {
		return false, err
	}
	return set, nil
}

func (r *RedisStorage) FlushAll() error {
	return r.Client.FlushAll(context.Background()).Err()
}
