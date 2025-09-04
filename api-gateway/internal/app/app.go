package app

import (
	"todo/api/internal/grpc/client"
	"todo/api/internal/http/handlers"
	"todo/api/internal/http/router"
	"todo/api/internal/http/server"
	producer "todo/api/internal/kafka"
)

type App struct {
	Server *server.Server
}

func New(serverAddr string, serviceAddr string) *App {
	grpclient, err := client.New(serviceAddr)
	if err != nil {
		panic("grpc server not connected")
	}

	writer := producer.New(
		[]string{"kafka:9092"},
		"tasks",
	)

	handlers := handlers.New(grpclient, writer)
	router := router.New(handlers).InitRouter()
	app := server.New(serverAddr, router)

	return &App{
		Server: app,
	}
}
