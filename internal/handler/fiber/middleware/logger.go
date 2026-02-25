package middleware

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func Logger(log *slog.Logger) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		start := time.Now()

		err := ctx.Next()

		duration := time.Since(start)

		entry := log.With(
			slog.String("method", ctx.Method()),
			slog.String("remote_addr", ctx.IP()),
			slog.String("user_agent", ctx.Get(fiber.HeaderUserAgent)),
			slog.String("request_id", requestid.FromContext(ctx)),
			slog.Duration("latency", duration),
		)

		entry.Info(
			ctx.OriginalURL(),
			slog.String("status", strconv.Itoa(ctx.Response().StatusCode())),
			slog.String("bytes", strconv.Itoa(len(ctx.Response().Body()))),
		)

		return err
	}
}
