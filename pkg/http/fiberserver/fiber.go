package fiberserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/sync/errgroup"

	"bookmarks/pkg/http"
)

type server struct {
	ctx context.Context //nolint:containedctx
	eg  *errgroup.Group
	log *slog.Logger

	app    *fiber.App
	notify chan error

	prefork         bool
	address         string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	idleTimeout     time.Duration
	shutdownTimeout time.Duration
}

func New(
	logger *slog.Logger,
	handler func(s *fiber.App),
	options ...Option,
) http.Server {
	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(1) // Run only one goroutine

	s := &server{
		app:     nil,
		ctx:     ctx,
		eg:      group,
		log:     logger,
		notify:  make(chan error, 1),
		prefork: false,
	}

	for _, opt := range options {
		opt(s)
	}

	app := fiber.New(fiber.Config{
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		IdleTimeout:  s.idleTimeout,
		JSONDecoder:  json.Unmarshal,
		JSONEncoder:  json.Marshal,
	})

	handler(app)

	s.app = app

	return s
}

func (s *server) Start() {
	const op = "http.fiber.Start"

	s.eg.Go(func() error {
		err := s.app.Listen(s.address, fiber.ListenConfig{
			EnablePrefork: s.prefork,
		})
		if err != nil {
			s.notify <- err

			close(s.notify)

			return err
		}

		return nil
	})

	s.log.Info("Start", slog.String("op", op))
}

func (s *server) Notify() <-chan error {
	return s.notify
}

func (s *server) Shutdown() error {
	const op = "http.fiber.Shutdown"
	var shutdownErrors []error

	log := s.log.With(
		slog.String("op", op),
	)

	err := s.app.ShutdownWithTimeout(s.shutdownTimeout)
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Error(err.Error())

		shutdownErrors = append(shutdownErrors, fmt.Errorf("%s: %w", op, err))
	}

	// Wait for all goroutines to finish and get any error
	err = s.eg.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Error(err.Error())

		shutdownErrors = append(shutdownErrors, fmt.Errorf("%s: %w", op, err))
	}

	log.Info("Shutdown")

	return errors.Join(shutdownErrors...)
}
