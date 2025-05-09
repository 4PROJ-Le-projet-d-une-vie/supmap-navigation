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
	key := formatKey(session.ID)
	return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r RedisSessionCache) GetSession(ctx context.Context, sessionID string) (*navigation.Session, error) {
	key := formatKey(sessionID)
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

func (r RedisSessionCache) DeleteSession(ctx context.Context, sessionID string) error {
	key := formatKey(sessionID)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("deleting session: %w", err)
	}
	return nil
}

func formatKey(sessionID string) string {
	return fmt.Sprintf("navigation:session:%s", sessionID)
}
