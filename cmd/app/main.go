package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"bookmarks/internal/config"
	"bookmarks/internal/handler/fiber"
	fiberv1 "bookmarks/internal/handler/fiber/v1"
	"bookmarks/internal/handler/net"
	netv1 "bookmarks/internal/handler/net/v1"
	bookmarkRepo "bookmarks/internal/repository/bookmark"
	bookmarkServ "bookmarks/internal/service/bookmark"
	"bookmarks/internal/storage/memory"
	"bookmarks/pkg/http"
	"bookmarks/pkg/http/fiberserver"
	"bookmarks/pkg/http/netserver"
)

const (
	envLocal  = "local"
	envProd   = "production"
	servFiber = "fiber"
	servCore  = "net/http"
)

func main() {
	var err error

	cfg := config.NewConfig()
	log := setupLogger(cfg.Env)

	log.Debug("app main", slog.Any("config", cfg))

	server := makeServer(log, cfg.HTTPServer)
	server.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("application.signal", slog.String("signal", s.String()))
	case err = <-server.Notify():
		log.Error("application.Notify", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
	}

	// Shutdown
	err = server.Shutdown()
	if err != nil {
		log.Error("application.Shutdown", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func makeServer(log *slog.Logger, cfg config.HTTPServer) http.Server {
	switch cfg.Type {
	case servFiber:
		// handler(service)
		return fiberserver.New(
			log,
			fiber.Register(
				log,
				fiberv1.NewHandler(log, makeBookmarkService()),
			),
			fiberserver.Address(cfg.Address),
			fiberserver.ReadTimeout(cfg.Timeout),
			fiberserver.WriteTimeout(cfg.Timeout),
			fiberserver.ShutdownTimeout(cfg.Timeout),
			fiberserver.IdleTimeout(cfg.IdleTimeout),
		)
	default:
		return netserver.New(
			log,
			net.Register(
				log,
				netv1.NewHandler(log, makeBookmarkService()),
			),
			netserver.Address(cfg.Address),
			netserver.ReadTimeout(cfg.Timeout),
			netserver.WriteTimeout(cfg.Timeout),
			netserver.ShutdownTimeout(cfg.Timeout),
			netserver.IdleTimeout(cfg.IdleTimeout),
		)
	}
}

func makeBookmarkService() netv1.Service {
	// storage [selected]
	storage := memory.NewStorage()
	// repository(storage)
	repo := bookmarkRepo.NewRepository(storage)
	// service(repository)
	return bookmarkServ.NewService(repo)
}
