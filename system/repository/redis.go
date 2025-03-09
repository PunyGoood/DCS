package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisRepository[T any] struct {
	client    *redis.Client
	keyPrefix string
}

func NewRedisRepository[T any](client *redis.Client, keyPrefix string) Repository[T] {
	return &RedisRepository[T]{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

func (r *RedisRepository[T]) GetByID(id int64) (*T, error) {
	key := fmt.Sprintf("%s:%d", r.keyPrefix, id)

	val, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("key %s not found", key)
		}
		return nil, err
	}

	var result T
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %v", err)
	}

	return &result, nil
}

// 可选：添加Set方法用于存储数据
func (r *RedisRepository[T]) Set(entity *T, id int64) error {
	key := fmt.Sprintf("%s:%d", r.keyPrefix, id)

	data, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	return r.client.Set(context.Background(), key, data, 0).Err()
}
