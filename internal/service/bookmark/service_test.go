package bookmark

import (
	"testing"

	"github.com/stretchr/testify/require"

	"bookmarks/internal/repository/bookmark"
	"bookmarks/internal/storage/memory"
)

func TestAppend_Success(t *testing.T) {
	storage := memory.NewStorage()
	repo := bookmark.NewRepository(storage)
	srv := NewService(repo)

	bookmark, err := srv.Append("test", "value")
	require.NoError(t, err)
	require.Equal(t, "test", bookmark.Title)
	require.Equal(t, "value", bookmark.Value)
}

func TestAppend_NotErrExists(t *testing.T) {
	storage := memory.NewStorage()
	repo := bookmark.NewRepository(storage)
	srv := NewService(repo)

	entity1, err := srv.Append("test", "value")
	require.NoError(t, err)

	entity2, err := srv.Append("test", "value")
	require.NoError(t, err)
	require.Equal(t, entity1.Uuid, entity2.Uuid)
}
