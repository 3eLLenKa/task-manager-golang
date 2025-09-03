package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"todo/db/internal/app"
	"todo/db/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := InitLogger(cfg.Env)

	application := app.New(
		log,
		cfg.GRPC.Port,
		cfg.Postgres.DSN,
		cfg.Redis.DSN,
		cfg.Redis.TTL,
	)

	go func() {
		application.Server.MustRun()
	}()

	log.Info("zaebis rabotaet")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.Server.Stop()

	log.Info("Gracefully stopped")

}

func InitLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		panic("failed to init logger")
	}

	return log
}
