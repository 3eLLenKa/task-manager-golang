package app

import "todo/api/internal/http/server"

type App struct {
	app *server.Server
}

func New(addr string) *App {
	// grpc := client.New()
	// handlers := handlers.New()
	// router := router.New().InitRouter()
	// server := server.New()

	// return &App{
	// 	app: app,
	// }

	return &App{}
}
