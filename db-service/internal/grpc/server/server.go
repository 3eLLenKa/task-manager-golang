package server

import (
	"fmt"
	"log/slog"
	"net"
	"todo/db/internal/grpc/handlers"
	"todo/db/internal/lib/sl"
	"todo/db/internal/service"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	log        *slog.Logger
	gRPCserver *grpc.Server
	port       int
}

func New(log *slog.Logger, taskService *service.TaskService, port int) *Server {
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(),
	))

	handlers.Register(server, taskService)

	reflection.Register(server)

	return &Server{
		log:        log,
		gRPCserver: server,
		port:       port,
	}
}

func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		panic("failed to run gRPC server")
	}
}

func (s *Server) Run() error {
	const op = "server.Run"

	log := s.log.With(
		slog.String("op", op),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))

	if err != nil {
		log.Error("lisening error", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.gRPCserver.Serve(l); err != nil {
		log.Error("failed to serve", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Server) Stop() {
	const op = "server.Stop"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("Gracefully stop")

	s.gRPCserver.GracefulStop()
}
