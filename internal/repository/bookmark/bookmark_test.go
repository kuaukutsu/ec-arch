package bookmark

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"bookmarks/internal/model"
	core "bookmarks/internal/repository"
	"bookmarks/internal/storage/memory"
)

func TestCreate_Success(t *testing.T) {
	storage := memory.NewStorage()
	repo := NewRepository(storage)

	bookmark, err := model.NewBookmark("test", "value")
	require.NoError(t, err)

	bookmark, err = repo.Create(bookmark)
	require.NoError(t, err)
	require.Equal(t, "test", bookmark.Title)
	require.Equal(t, "value", bookmark.Value)
	require.GreaterOrEqual(t, time.Now(), bookmark.CreatedAt)
}

func TestCreate_ErrExists(t *testing.T) {
	bookmark, _ := model.NewBookmark("test", "value")
	repo := makePrepareRepository(bookmark)

	_, err := repo.Create(bookmark)
	require.ErrorIs(t, err, core.ErrExists)
}

func TestGetUUID_Success(t *testing.T) {
	bookmark, _ := model.NewBookmark("test", "value")
	repo := makePrepareRepository(bookmark)

	entity, err := repo.GetByUUID(bookmark.Uuid)
	require.NoError(t, err)
	require.Equal(t, "test", entity.Title)
	require.Equal(t, bookmark.Uuid, entity.Uuid)
}

func TestGetUUID_NotFound(t *testing.T) {
	bookmark, _ := model.NewBookmark("test", "value")
	repo := makePrepareRepository(bookmark)

	uuid7, _ := uuid.NewV7()
	_, err := repo.GetByUUID(uuid7)
	require.ErrorIs(t, err, core.ErrNotFound)
}

func TestGetValue_Success(t *testing.T) {
	bookmark, _ := model.NewBookmark("test", "value")
	repo := makePrepareRepository(bookmark)

	entity, err := repo.GetByValue("value")
	require.NoError(t, err)
	require.Equal(t, "test", entity.Title)
	require.Equal(t, bookmark.Uuid, entity.Uuid)
}

func TestGetValue_NotFound(t *testing.T) {
	bookmark, _ := model.NewBookmark("test", "value")
	repo := makePrepareRepository(bookmark)

	_, err := repo.GetByValue("not-found")
	require.ErrorIs(t, err, core.ErrNotFound)
}

func TestDelete_Success(t *testing.T) {
	bookmark, _ := model.NewBookmark("test", "value")
	repo := makePrepareRepository(bookmark)

	err := repo.Delete(bookmark.Uuid)
	require.NoError(t, err)

	_, err = repo.GetByUUID(bookmark.Uuid)
	require.ErrorIs(t, err, core.ErrNotFound)

	_, err = repo.GetByValue("value")
	require.ErrorIs(t, err, core.ErrNotFound)
}

func TestDelete_NotFound(t *testing.T) {
	bookmark, _ := model.NewBookmark("test", "value")
	repo := makePrepareRepository(bookmark)

	uuid7, _ := uuid.NewV7()
	err := repo.Delete(uuid7)
	require.ErrorIs(t, err, core.ErrNotFound)
}

func makePrepareRepository(bookmark model.Bookmark) *repository {
	storage := memory.NewStorage()
	repo := NewRepository(storage)

	_, _ = repo.Create(bookmark)

	return repo
}
