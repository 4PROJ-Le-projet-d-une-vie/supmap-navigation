package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"supmap-navigation/internal/navigation"
	"time"
)

type RedisSessionCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisSessionCache(client *redis.Client, ttl time.Duration) *RedisSessionCache {
	return &RedisSessionCache{client: client, ttl: ttl}
}

func (r RedisSessionCache) SetSession(ctx context.Context, session *navigation.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("marshalling session: %w", err)
	}
	key := formatKey(session.UserID)
	return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r RedisSessionCache) GetSession(ctx context.Context, userID string) (*navigation.Session, error) {
	key := formatKey(userID)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("getting session: %w", err)
	}
	var session navigation.Session
	if err := json.Unmarshal([]byte(val), &session); err != nil {
		return nil, fmt.Errorf("unmarshalling session: %w", err)
	}
	return &session, nil
}

func (r RedisSessionCache) DeleteSession(ctx context.Context, userID string) error {
	key := formatKey(userID)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("deleting session: %w", err)
	}
	return nil
}

func formatKey(userID string) string {
	return fmt.Sprintf("navigation:session:%s", userID)
}
