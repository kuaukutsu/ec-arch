package middleware

import (
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"bookmarks/internal/config"
)

func Logger(log *slog.Logger) func(c *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		err := ctx.Next()

		entry := log.With(
			slog.String("method", ctx.Method()),
			slog.String("remote_addr", ctx.Context().RemoteAddr().String()),
			slog.String("user_agent", string(ctx.Context().UserAgent())),
			slog.String("request_id", ctx.Locals(config.RequestID).(string)),
		)

		entry.Info(
			ctx.OriginalURL(),
			slog.String("status", strconv.Itoa(ctx.Response().StatusCode())),
			slog.String("bytes", strconv.Itoa(len(ctx.Response().Body()))),
		)

		return err
	}
}
