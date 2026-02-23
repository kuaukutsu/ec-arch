package v1

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"bookmarks/internal/config"
	router "bookmarks/internal/handler/fiber"
	"bookmarks/internal/model"
	"bookmarks/internal/service/bookmark"
)

type Service interface {
	Append(title, val string) (model.Bookmark, error)
	View(uuid string) (model.Bookmark, error)
	Change()
	Delete(uuid string) error
}

type bookmarkHandler struct {
	service   Service
	validator *validator.Validate
	logger    *slog.Logger
}

func NewHandler(l *slog.Logger, s Service) *bookmarkHandler {
	return &bookmarkHandler{
		service:   s,
		validator: validator.New(validator.WithRequiredStructEnabled()),
		logger:    l,
	}
}

func (h *bookmarkHandler) Append(ctx *fiber.Ctx) error {
	var input CreateBookmarkRequest

	log := h.logger.With(
		slog.String("op", "handler.v1.bookmark.Append"),
		slog.String("request_id", ctx.Locals(config.RequestID).(string)),
	)

	if err := ctx.BodyParser(&input); err != nil {
		log.Error(err.Error())
		return router.ErrorResponse(ctx, err.Error(), http.StatusBadRequest)
	}

	if err := h.validator.Struct(input); err != nil {
		log.Error(err.Error())
		return router.ErrorResponse(ctx, err.Error(), http.StatusBadRequest)
	}

	entity, err := h.service.Append(input.Title, input.Value)
	if err != nil {
		log.Error(err.Error())
		return router.ErrorResponse(ctx, err.Error(), http.StatusBadRequest)
	}

	return ctx.Status(http.StatusCreated).JSON(entity)
}

func (h *bookmarkHandler) View(ctx *fiber.Ctx) error {
	log := h.logger.With(
		slog.String("op", "handler.v1.bookmark.View"),
		slog.String("request_id", ctx.Locals(config.RequestID).(string)),
	)

	entity, err := h.service.View(ctx.Params("uuid"))
	if err != nil {
		log.Error(err.Error())

		if errors.Is(err, bookmark.ErrBookmarkNotFound) {
			return router.ErrorResponse(ctx, err.Error(), http.StatusNotFound)
		}

		return router.ErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
	}

	return ctx.Status(http.StatusOK).JSON(entity)
}

func (h *bookmarkHandler) Change(ctx *fiber.Ctx) error {
	return ctx.SendString("Change")
}

func (h *bookmarkHandler) Delete(ctx *fiber.Ctx) error {
	log := h.logger.With(
		slog.String("op", "handler.v1.bookmark.Delete"),
		slog.String("request_id", ctx.Locals(config.RequestID).(string)),
	)

	if err := h.service.Delete(ctx.Params("uuid")); err != nil {
		log.Error(err.Error())

		if errors.Is(err, bookmark.ErrBookmarkNotFound) {
			return router.ErrorResponse(ctx, err.Error(), http.StatusNotFound)
		}

		return router.ErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
	}

	return ctx.SendStatus(http.StatusNoContent)
}
