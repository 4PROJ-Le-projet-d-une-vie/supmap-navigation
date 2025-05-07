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
	"supmap-navigation/internal/cache"
	"supmap-navigation/internal/config"
	"supmap-navigation/internal/subscriber"
	"supmap-navigation/internal/ws"
	"syscall"
	"time"
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

	var loggerOpts slog.HandlerOptions
	if conf.Env == config.EnvDev {
		loggerOpts = slog.HandlerOptions{Level: slog.LevelDebug}
	}

	jsonHandler := slog.NewJSONHandler(os.Stdout, &loggerOpts)
	logger := slog.New(jsonHandler)

	redisClient := redis.NewClient(&redis.Options{Addr: net.JoinHostPort(conf.RedisHost, conf.RedisPort)})
	sub := subscriber.NewSubscriber(conf, logger, redisClient, conf.RedisIncidentsChannel, 10)
	sessionCache := cache.NewRedisSessionCache(redisClient, 30*time.Minute)

	go func() {
		if err := sub.Start(ctx); err != nil {
			logger.Error("subscriber stopped with error: ", err)
		}
	}()

	wsManager := ws.NewManager(ctx, logger, sessionCache)
	go wsManager.Start()

	server := api.NewServer(conf, wsManager, logger)
	if err := server.Start(ctx); err != nil {
		return err
	}

	return nil
}
