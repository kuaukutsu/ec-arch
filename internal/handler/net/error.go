package net

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"bookmarks/internal/handler"
)

func ErrorResponse(w http.ResponseWriter, r *http.Request, err string, status int) {
	ctx := &handler.ErrorContext{
		RequestID: middleware.GetReqID(r.Context()),
	}

	if status >= 500 {
		err = http.StatusText(status)
	}

	render.Status(r, status)
	render.JSON(w, r, handler.NewError(err, ctx))
}
