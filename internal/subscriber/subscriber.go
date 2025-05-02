package subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"supmap-navigation/internal/config"
	"supmap-navigation/internal/incidents"
)

type Subscriber struct {
	config        *config.Config
	logger        *slog.Logger
	client        *redis.Client
	topic         string
	maxConcurrent int
}

func NewSubscriber(config *config.Config, logger *slog.Logger, client *redis.Client, topic string, maxConcurrent int) *Subscriber {
	return &Subscriber{
		config,
		logger,
		client,
		topic,
		maxConcurrent,
	}
}

func (s *Subscriber) Start(ctx context.Context) error {
	s.logger.Info("Redis subscriber is running", "topic", s.topic)
	pubsub := s.client.Subscribe(ctx, s.topic)
	defer func() {
		if err := pubsub.Close(); err != nil {
			s.logger.Warn("failed to close pubsub", "error", err)
		}
	}()

	msgCh := pubsub.Channel()

	// Improvised semaphore to avoid too many go routines (handleMessage).
	// A better but more complex solution would be a worker pool. But franchement OSEF.
	// https://medium.com/@deckarep/gos-extended-concurrency-semaphores-part-1-5eeabfa351ce
	// https://gobyexample.com/worker-pools
	sem := make(chan struct{}, s.maxConcurrent)
	for {
		select {
		case msg, ok := <-msgCh:
			if !ok {
				s.logger.Warn("pubsub channel closed by Redis")
				return nil
			}
			select {
			case sem <- struct{}{}:
				go func(m *redis.Message) {
					defer func() { <-sem }()
					if err := s.handleMessage(ctx, m); err != nil {
						s.logger.Error("error handling message", "error", err)
					}
				}(msg)
			case <-ctx.Done():
				return nil
			}
		case <-ctx.Done():
			s.logger.Info("shutting down Redis subscriber")
			return nil
		}
	}
}

func (s *Subscriber) handleMessage(ctx context.Context, msg *redis.Message) error {
	var incident incidents.Incident

	if err := json.Unmarshal([]byte(msg.Payload), &incident); err != nil {
		return fmt.Errorf("failed to decode incidents: %w", err)
	}

	if err := incident.Validate(); err != nil {
		return fmt.Errorf("invalid incidents: %w", err)
	}

	s.logger.Info("incidents received", "incidents", incident)
	return nil
}
