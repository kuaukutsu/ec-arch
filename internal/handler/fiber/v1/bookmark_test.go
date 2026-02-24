package v1

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/render"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/stretchr/testify/require"

	"bookmarks/internal/config"
	"bookmarks/internal/handler"
	"bookmarks/internal/model"
	repo "bookmarks/internal/repository/bookmark"
	srv "bookmarks/internal/service/bookmark"
	"bookmarks/internal/storage/memory"
)

func TestAppend_Success(t *testing.T) {
	target := "/v1/bookmark/append"

	hdl := makeHandler()
	app := makeFiber(target, hdl.Append)

	// 2. Создаем тестовый запрос
	body := strings.NewReader(`{"title": "test", "value": "value"}`)
	req := httptest.NewRequest(http.MethodPost, target, body)
	req.Header.Set("Content-Type", "application/json")

	// Выполняем тест (второй аргумент - таймаут в мс, -1 — без таймаута)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	defer func() {
		_ = resp.Body.Close()
	}()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var response model.Bookmark
	err = render.DecodeJSON(resp.Body, &response)
	require.NoError(t, err)
	require.Equal(t, "value", response.Value)
}

func TestAppend_BodyEmptyError(t *testing.T) {
	target := "/v1/bookmark/append"

	hdl := makeHandler()
	app := makeFiber(target, hdl.Append)

	// 2. Создаем тестовый запрос
	body := strings.NewReader("")
	req := httptest.NewRequest(http.MethodPost, target, body)
	req.Header.Set("Content-Type", "application/json")

	// Выполняем тест (второй аргумент - таймаут в мс, -1 — без таймаута)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	defer func() {
		_ = resp.Body.Close()
	}()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response handler.ErrorResponse
	err = render.DecodeJSON(resp.Body, &response)
	require.NoError(t, err)
	require.Equal(t, "unexpected end of JSON input", response.Error)
}

func TestAppend_BodyBadRequest(t *testing.T) {
	target := "/v1/bookmark/append"

	hdl := makeHandler()
	app := makeFiber(target, hdl.Append)

	// 2. Создаем тестовый запрос
	body := strings.NewReader(`{"title": "test"}`)
	req := httptest.NewRequest(http.MethodPost, target, body)
	req.Header.Set("Content-Type", "application/json")

	// Выполняем тест (второй аргумент - таймаут в мс, -1 — без таймаута)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	defer func() {
		_ = resp.Body.Close()
	}()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response handler.ErrorResponse
	err = render.DecodeJSON(resp.Body, &response)
	require.NoError(t, err)
	require.Contains(t, response.Error, "Field validation for 'Value' failed")
}

func makeFiber(target string, handler fiber.Handler) *fiber.App {
	app := fiber.New()
	app.Use(requestid.New(requestid.Config{
		ContextKey: config.RequestID,
	}))

	app.Post(target, handler)

	return app
}

func makeHandler() *bookmarkHandler {
	storage := memory.NewBookmarkStorage()
	repository := repo.NewRepository(storage)
	service := srv.NewService(repository)
	logger := slog.New(slog.DiscardHandler)

	return NewHandler(logger, service)
}
