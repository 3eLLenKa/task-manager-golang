package app

import (
	"log/slog"
	"time"
	"todo/db/internal/grpc/server"
	"todo/db/internal/lib/sl"
	"todo/db/internal/service"
	"todo/db/internal/storage/postgres"
	"todo/db/internal/storage/redis"
)

type App struct {
	Server *server.Server
}

func New(
	log *slog.Logger,
	grpcPort int,
	postgresDsn string,
	redisDsn string,
	cacheTTL time.Duration,
) *App {
	pgStorage, err := postgres.New(postgresDsn)

	if err != nil {
		log.Error("failed to init postgres", sl.Err(err))
		panic(err)
	}

	redisCache, err := redis.New(redisDsn, cacheTTL)

	if err != nil {
		log.Error("failed to init redis", sl.Err(err))
		panic(err)
	}

	taskService := service.New(log, pgStorage, redisCache)

	grpcServer := server.New(log, taskService, grpcPort)

	return &App{
		Server: grpcServer,
	}
}
