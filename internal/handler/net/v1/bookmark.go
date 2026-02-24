package v1

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"bookmarks/internal/config"
	"bookmarks/internal/handler/net"
	"bookmarks/internal/model"
	"bookmarks/internal/service/bookmark"
)

var (
	ErrRequestBodyIsEmpty = errors.New("request body is empty")
	ErrUUIDIsEmpty        = errors.New("uuid is empty")
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

func (h *bookmarkHandler) Append(w http.ResponseWriter, r *http.Request) {
	var input CreateBookmarkRequest

	log := h.logger.With(
		slog.String("op", "handler.v1.bookmark.Append"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	err := render.DecodeJSON(r.Body, &input)
	if errors.Is(err, io.EOF) {
		log.Error(ErrRequestBodyIsEmpty.Error())
		net.ErrorResponse(w, r, ErrRequestBodyIsEmpty.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Error(err.Error())
		net.ErrorResponse(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(input); err != nil {
		log.Error(err.Error())
		net.ErrorResponse(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	entity, err := h.service.Append(input.Title, input.Value)
	if err != nil {
		log.Error(err.Error())

		if errors.Is(err, bookmark.ErrBookmarkExists) {
			error := fmt.Sprintf("[%s] %s", entity.Uuid, bookmark.ErrBookmarkExists)
			net.ErrorResponse(w, r, error, http.StatusConflict)
			return
		}

		net.ErrorResponse(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, entity)
}

func (h *bookmarkHandler) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.logger.With(
		slog.String("op", "handler.v1.bookmark.View"),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	uuid, err := prepareUuid(ctx)
	if err != nil {
		log.Error(err.Error())
		net.ErrorResponse(w, r, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	entity, err := h.service.View(uuid)
	if err != nil {
		log.Error(err.Error())

		if errors.Is(err, bookmark.ErrBookmarkNotFound) {
			net.ErrorResponse(w, r, bookmark.ErrBookmarkNotFound.Error(), http.StatusNotFound)
			return
		}

		net.ErrorResponse(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, entity)
}

func (h *bookmarkHandler) Change(w http.ResponseWriter, r *http.Request) {
}

func (h *bookmarkHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.logger.With(
		slog.String("op", "handler.v1.bookmark.Delete"),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	uuid, err := prepareUuid(ctx)
	if err != nil {
		log.Error(err.Error())
		net.ErrorResponse(w, r, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := h.service.Delete(uuid); err != nil {
		log.Error(err.Error())

		if errors.Is(err, bookmark.ErrBookmarkNotFound) {
			net.ErrorResponse(w, r, err.Error(), http.StatusNotFound)
			return
		}

		net.ErrorResponse(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

func prepareUuid(ctx context.Context) (string, error) {
	uuid, ok := ctx.Value(config.FieldUUID).(string)
	if !ok {
		return "", ErrUUIDIsEmpty
	}

	return uuid, nil
}
