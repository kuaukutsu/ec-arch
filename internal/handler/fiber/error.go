package fiber

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"bookmarks/internal/handler"
)

func ErrorResponse(ctx fiber.Ctx, err string, code int) error {
	errCtx := &handler.ErrorContext{
		RequestID: requestid.FromContext(ctx),
	}

	if code >= 500 {
		err = http.StatusText(code)
	}

	return ctx.Status(code).JSON(handler.NewError(err, errCtx))
}
