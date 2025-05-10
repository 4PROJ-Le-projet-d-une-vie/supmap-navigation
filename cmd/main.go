package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"supmap-navigation/internal/api"
	"supmap-navigation/internal/cache"
	"supmap-navigation/internal/config"
	routing "supmap-navigation/internal/gis/routing"
	"supmap-navigation/internal/incidents"
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
	sessionCache := cache.NewRedisSessionCache(redisClient, 30*time.Minute)

	wsManager := ws.NewManager(ctx, logger, sessionCache)

	supmapGISURL := fmt.Sprintf("http://%s:%s", conf.SupmapGISHost, conf.SupmapGISPort)
	routingClient := routing.NewClient(supmapGISURL)
	logger.Info("supmap-gis client initialized", "url", supmapGISURL)

	multicaster := incidents.NewMulticaster(wsManager, sessionCache, routingClient)
	sub := subscriber.NewSubscriber(conf, logger, redisClient, conf.RedisIncidentsChannel, 10, multicaster)

	go wsManager.Start()

	go func() {
		if err := sub.Start(ctx); err != nil {
			logger.Error("subscriber stopped with error: ", err)
		}
	}()

	server := api.NewServer(conf, wsManager, logger)
	if err := server.Start(ctx); err != nil {
		return err
	}

	return nil
}
