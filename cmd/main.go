package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"supmap-navigation/internal/api"
	"supmap-navigation/internal/config"
	"supmap-navigation/internal/subscriber"
	"supmap-navigation/internal/ws"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conf, err := config.New()
	if err != nil {
		return err
	}

	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonHandler)

	redisClient := redis.NewClient(&redis.Options{Addr: net.JoinHostPort(conf.RedisHost, conf.RedisPort)})
	sub := subscriber.NewSubscriber(conf, logger, redisClient, "incidents", 10)

	go func() {
		if err := sub.Start(ctx); err != nil {
			logger.Error("subscriber stopped with error: ", err)
		}
	}()

	wsManager := ws.NewManager(ctx)
	go wsManager.Start()

	server := api.NewServer(conf, wsManager, logger)
	if err := server.Start(ctx); err != nil {
		return err
	}

	return nil
}
