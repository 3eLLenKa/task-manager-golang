package app

import (
	"todo/api/internal/grpc/client"
	"todo/api/internal/http/handlers"
	"todo/api/internal/http/router"
	"todo/api/internal/http/server"
)

type App struct {
	Server *server.Server
}

func New(
	serverAddr string,
	serviceAddr string,
) *App {
	grpclient, err := client.New(serviceAddr)

	if err != nil {
		panic("grpc server not connected")
	}

	handlers := handlers.New(grpclient)
	router := router.New(handlers).InitRouter()
	app := server.New(serverAddr, router)

	return &App{
		Server: app,
	}
}
