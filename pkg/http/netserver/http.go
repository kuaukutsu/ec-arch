package netserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	nethttp "net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"bookmarks/pkg/http"
)

type server struct {
	ctx context.Context //nolint:containedctx
	eg  *errgroup.Group
	log *slog.Logger

	app    *nethttp.Server
	notify chan error

	address         string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	idleTimeout     time.Duration
	shutdownTimeout time.Duration
}

func New(
	logger *slog.Logger,
	handler func(s *nethttp.Server),
	options ...Option,
) http.Server {
	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(1) // Run only one goroutine

	s := &server{
		app:    nil,
		ctx:    ctx,
		eg:     group,
		log:    logger,
		notify: make(chan error, 1),
	}

	for _, opt := range options {
		opt(s)
	}

	app := &nethttp.Server{
		Addr:         s.address,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		IdleTimeout:  s.idleTimeout,
	}

	handler(app)

	s.app = app

	return s
}

func (s *server) Start() {
	const op = "http.net.Start"

	s.eg.Go(func() error {
		err := s.app.ListenAndServe()
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
	const op = "http.net.Shutdown"
	var shutdownErrors []error

	log := s.log.With(
		slog.String("op", op),
	)

	ctx, cancel := context.WithTimeout(s.ctx, s.shutdownTimeout)
	defer cancel()

	err := s.app.Shutdown(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Error(err.Error())

		shutdownErrors = append(shutdownErrors, fmt.Errorf("%s: %w", op, err))
	}

	err = s.eg.Wait()
	if err != nil && !errors.Is(err, nethttp.ErrServerClosed) {
		log.Error(err.Error())

		shutdownErrors = append(shutdownErrors, fmt.Errorf("%s: %w", op, err))
	}

	log.Info("Shutdown")

	return errors.Join(shutdownErrors...)
}
