package v1

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/render"
	"github.com/stretchr/testify/require"

	"bookmarks/internal/handler"
	"bookmarks/internal/model"
	repo "bookmarks/internal/repository/bookmark"
	srv "bookmarks/internal/service/bookmark"
	"bookmarks/internal/storage/memory"
)

func TestAppend_Success(t *testing.T) {
	hdl := makeHandler()

	// 2. Создаем тестовый запрос
	body := strings.NewReader(`{"title": "test", "value": "value"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/bookmark/append", body)
	req.Header.Set("Content-Type", "application/json")

	// 3. Создаем "записыватель" ответа
	rr := httptest.NewRecorder()

	hdl.Append(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)

	var response model.Bookmark
	err := render.DecodeJSON(rr.Body, &response)

	require.NoError(t, err)
	require.Equal(t, "value", response.Value)
}

func TestAppend_BodyBadRequest(t *testing.T) {
	hdl := makeHandler()

	// 2. Создаем тестовый запрос
	body := strings.NewReader(`{"title": "test"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/bookmark/append", body)
	req.Header.Set("Content-Type", "application/json")

	// 3. Создаем "записыватель" ответа
	rr := httptest.NewRecorder()

	hdl.Append(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)

	var response handler.ErrorResponse
	err := render.DecodeJSON(rr.Body, &response)
	require.NoError(t, err)
	require.Contains(t, response.Error, "Field validation for 'Value' failed")
}

func makeHandler() *bookmarkHandler {
	storage := memory.NewBookmarkStorage()
	repository := repo.NewRepository(storage)
	service := srv.NewService(repository)
	logger := slog.New(slog.DiscardHandler)

	return NewHandler(logger, service)
}
