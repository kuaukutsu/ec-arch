package fiber

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"bookmarks/internal/config"
	"bookmarks/internal/handler"
)

func ErrorResponse(ctx *fiber.Ctx, err string, code int) error {
	errCtx := &handler.ErrorContext{
		RequestID: ctx.Locals(config.RequestID).(string),
	}

	if code >= 500 {
		err = http.StatusText(code)
	}

	return ctx.Status(code).JSON(handler.NewError(err, errCtx))
}
