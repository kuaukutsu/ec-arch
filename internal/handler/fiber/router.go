package fiber

import (
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"bookmarks/internal/config"
	"bookmarks/internal/handler/fiber/middleware"
)

type BookmarkHandler interface {
	Append(ctx *fiber.Ctx) error
	View(ctx *fiber.Ctx) error
	Change(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}

func Register(
	log *slog.Logger,
	bookmarkHnd BookmarkHandler,
) func(s *fiber.App) {
	return func(s *fiber.App) {
		s.Use(requestid.New(requestid.Config{
			ContextKey: config.RequestID,
		}))
		s.Use(middleware.Logger(log))

		s.Get("/health", healthHandler)

		v1 := s.Group("/v1")

		bookmark := v1.Group("/bookmark")
		bookmark.Post("/append", bookmarkHnd.Append)
		bookmark.Get("/:uuid", bookmarkHnd.View)
		bookmark.Post("/:uuid", bookmarkHnd.Change)
		bookmark.Delete("/:uuid", bookmarkHnd.Delete)
	}
}

func healthHandler(ctx *fiber.Ctx) error {
	data := map[string]string{
		"status": "ok",
	}

	return ctx.Status(http.StatusOK).JSON(data)
}
