package storage

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	Client *redis.Client
}

var ctx = context.Background()

func NewRedisStorage(addr string) *RedisStorage {
	return &RedisStorage{
		Client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (r *RedisStorage) Increment(key string) (int, error) {
	count, err := r.Client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	r.Client.Expire(ctx, key, time.Second)
	return int(count), nil
}

func (r *RedisStorage) Get(key string) (int, error) {
	val, err := r.Client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}
