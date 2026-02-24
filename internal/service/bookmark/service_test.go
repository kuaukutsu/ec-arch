package bookmark

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	"bookmarks/internal/repository/bookmark"
	"bookmarks/internal/storage/memory"
)

func TestAppend_Success(t *testing.T) {
	storage := memory.NewBookmarkStorage()
	repo := bookmark.NewRepository(storage)
	srv := NewService(repo)

	title := gofakeit.Word()
	value := gofakeit.CarModel()
	
	bookmark, err := srv.Append(title, value)
	require.NoError(t, err)
	require.Equal(t, title, bookmark.Title)
	require.Equal(t, value, bookmark.Value)
}

func TestAppend_ErrExists(t *testing.T) {
	storage := memory.NewBookmarkStorage()
	repo := bookmark.NewRepository(storage)
	srv := NewService(repo)

	value := gofakeit.CarModel()
	
	_, err := srv.Append(gofakeit.Word(), value)
	require.NoError(t, err)

	_, err = srv.Append(gofakeit.Word(), value)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrBookmarkExists)
}
