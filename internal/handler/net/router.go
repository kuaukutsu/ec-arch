package net

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"bookmarks/internal/config"
	customMiddleware "bookmarks/internal/handler/net/middleware"
)

type BookmarkHandler interface {
	Append(w http.ResponseWriter, r *http.Request)
	View(w http.ResponseWriter, r *http.Request)
	Change(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func Register(
	log *slog.Logger,
	bookmarkHnd BookmarkHandler,
) func(*http.Server) {
	return func(s *http.Server) {
		router := chi.NewRouter()

		router.Use(middleware.RequestID)
		router.Use(middleware.Logger)
		router.Use(customMiddleware.Logger(log))
		router.Use(middleware.Recoverer)
		router.Use(middleware.URLFormat)

		router.Get("/health", healthHandler)

		router.Route("/v1", func(r chi.Router) {
			r.Route("/bookmark", func(r chi.Router) {
				r.Post("/append", bookmarkHnd.Append)

				r.Route("/{uuid}", func(r chi.Router) {
					r.Use(uuidCtx)

					r.Get("/", bookmarkHnd.View)
					r.Post("/", bookmarkHnd.Change)
					r.Delete("/", bookmarkHnd.Delete)
				})
			})

			r.Get("/bookmarks", func(w http.ResponseWriter, r *http.Request) {
				data := map[string]string{"status": "bookmarks list"}
				render.Status(r, http.StatusOK)
				render.JSON(w, r, data)
			})
		})

		s.Handler = router
	}
}

func uuidCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid := chi.URLParam(r, "uuid")
		ctx := context.WithValue(r.Context(), config.FieldUUID, uuid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, data)
}
