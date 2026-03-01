package fiber

import (
	"log/slog"
	"net/http"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/swaggo/swag"

	"bookmarks/docs"
	"bookmarks/internal/handler/fiber/middleware"
)

type BookmarkHandler interface {
	Append(ctx fiber.Ctx) error
	View(ctx fiber.Ctx) error
	Change(ctx fiber.Ctx) error
	Delete(ctx fiber.Ctx) error
}

// Swagger spec:
// @title       Go Example REST API
// @version     1.0
// @host        localhost:8082
// @BasePath    /v1
func Register(
	log *slog.Logger,
	bookmarkHnd BookmarkHandler,
) func(s *fiber.App) {
	swag.Register(swag.Name, docs.SwaggerInfo)

	return func(s *fiber.App) {
		s.Use(requestid.New())
		s.Use(middleware.Logger(log))

		s.Get("/health", healthHandler)

		v1 := s.Group("/v1")

		v1.Get("/swagger/*", swaggo.HandlerDefault)

		bookmark := v1.Group("/bookmark")
		bookmark.Post("/append", bookmarkHnd.Append)
		bookmark.Get("/:uuid<guid>", bookmarkHnd.View)
		bookmark.Post("/:uuid<guid>", bookmarkHnd.Change)
		bookmark.Delete("/:uuid<guid>", bookmarkHnd.Delete)
	}
}

func healthHandler(ctx fiber.Ctx) error {
	data := map[string]string{
		"status": "ok",
	}

	return ctx.Status(http.StatusOK).JSON(data)
}
