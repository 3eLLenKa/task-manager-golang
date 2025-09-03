package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo/api/internal/app"
	"todo/api/internal/config"
)

func main() {
	cfg := config.MustLoad()

	application := app.New(cfg.Serv.HTTP.Host, cfg.GRPC.DBService.Address)

	go func() {
		application.Server.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Serv.HTTP.ShutdownTimeout*time.Second)
	defer cancel()

	application.Server.Stop(ctx)
}
