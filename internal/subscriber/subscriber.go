package subscriber

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"supmap-navigation/internal/config"
)

type Subscriber struct {
	config *config.Config
	logger *slog.Logger
	client *redis.Client
	topic  string
}

func NewSubscriber(config *config.Config, logger *slog.Logger, client *redis.Client, topic string) *Subscriber {
	return &Subscriber{
		config,
		logger,
		client,
		topic,
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

	for {
		select {
		case msg, ok := <-msgCh:
			if !ok {
				s.logger.Warn("pubsub channel closed by Redis")
				return nil
			}
			if err := s.handleMessage(ctx, msg); err != nil {
				s.logger.Error("error handling message", "error", err)
			}
		case <-ctx.Done():
			s.logger.Info("shutting down Redis subscriber")
			return nil
		}
	}
}

func (s *Subscriber) handleMessage(ctx context.Context, msg *redis.Message) error {
	s.logger.Info("received message", "payload", msg.Payload)
	return nil
}
