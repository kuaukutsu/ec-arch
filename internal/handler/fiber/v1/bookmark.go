package v1

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"

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

// @Summary     Append bookmark
// @Description Append bookmark to collection
// @ID          append
// @Tags  	    bookmark
// @Accept      json
// @Produce     json
// @Success     200 {object} model.Bookmark
// @Failure     500 {object} handler.ErrorResponse
// @Router      /bookmark/append [post]
func (h *bookmarkHandler) Append(ctx fiber.Ctx) error {
	var input CreateBookmarkRequest

	log := h.logger.With(
		slog.String("op", "handler.v1.bookmark.Append"),
		slog.String("request_id", requestid.FromContext(ctx)),
	)

	if err := ctx.Bind().Body(&input); err != nil {
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

		if errors.Is(err, bookmark.ErrBookmarkExists) {
			error := fmt.Sprintf("[%s] %s", entity.Uuid, bookmark.ErrBookmarkExists)
			return router.ErrorResponse(ctx, error, http.StatusConflict)
		}

		return router.ErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
	}

	return ctx.Status(http.StatusCreated).JSON(entity)
}

// @Summary     Show bookmark
// @Description Show info bookmark
// @ID          view
// @Tags  	    bookmark
// @Accept      json
// @Produce     json
// @Param       uuid   path      string  true  "Bookmark UUID"
// @Success     200 {object} model.Bookmark
// @Failure     500 {object} handler.ErrorResponse
// @Router      /bookmark/{uuid} [get]
func (h *bookmarkHandler) View(ctx fiber.Ctx) error {
	log := h.logger.With(
		slog.String("op", "handler.v1.bookmark.View"),
		slog.String("request_id", requestid.FromContext(ctx)),
	)

	entity, err := h.service.View(ctx.Params("uuid"))
	if err != nil {
		log.Error(err.Error())

		if errors.Is(err, bookmark.ErrBookmarkNotFound) {
			return router.ErrorResponse(ctx, bookmark.ErrBookmarkNotFound.Error(), http.StatusNotFound)
		}

		return router.ErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
	}

	return ctx.Status(http.StatusOK).JSON(entity)
}

func (h *bookmarkHandler) Change(ctx fiber.Ctx) error {
	return ctx.SendString("Change")
}

func (h *bookmarkHandler) Delete(ctx fiber.Ctx) error {
	log := h.logger.With(
		slog.String("op", "handler.v1.bookmark.Delete"),
		slog.String("request_id", requestid.FromContext(ctx)),
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
