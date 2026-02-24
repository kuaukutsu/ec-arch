package bookmark

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"bookmarks/internal/model"
	core "bookmarks/internal/repository"
	"bookmarks/internal/storage/memory"
)

func TestCreate_Success(t *testing.T) {
	storage := memory.NewBookmarkStorage()
	repo := NewRepository(storage)

	bookmark := makeBookmark()

	entity, err := repo.Create(bookmark)
	require.NoError(t, err)
	require.Equal(t, bookmark.Title, entity.Title)
	require.Equal(t, bookmark.Value, entity.Value)
	require.GreaterOrEqual(t, time.Now(), entity.CreatedAt)
}

func TestCreate_ErrExists(t *testing.T) {
	bookmark := makeBookmark()
	repo := makePrepareRepository(bookmark)

	_, err := repo.Create(bookmark)
	require.ErrorIs(t, err, core.ErrExists)
}

func TestGetUUID_Success(t *testing.T) {
	bookmark := makeBookmark()
	repo := makePrepareRepository(bookmark)

	entity, err := repo.GetByUUID(bookmark.Uuid)
	require.NoError(t, err)
	require.Equal(t, bookmark.Title, entity.Title)
	require.Equal(t, bookmark.Uuid, entity.Uuid)
}

func TestGetUUID_NotFound(t *testing.T) {
	bookmark := makeBookmark()
	repo := makePrepareRepository(bookmark)

	uuid7, _ := uuid.NewV7()
	_, err := repo.GetByUUID(uuid7)
	require.ErrorIs(t, err, core.ErrNotFound)
}

func TestGetValue_Success(t *testing.T) {
	bookmark := makeBookmark()
	repo := makePrepareRepository(bookmark)

	entity, err := repo.GetByValue(bookmark.Value)
	require.NoError(t, err)
	require.Equal(t, bookmark.Title, entity.Title)
	require.Equal(t, bookmark.Uuid, entity.Uuid)
}

func TestGetValue_NotFound(t *testing.T) {
	bookmark := makeBookmark()
	repo := makePrepareRepository(bookmark)

	_, err := repo.GetByValue("not-found")
	require.ErrorIs(t, err, core.ErrNotFound)
}

func TestDelete_Success(t *testing.T) {
	bookmark := makeBookmark()
	repo := makePrepareRepository(bookmark)

	err := repo.Delete(bookmark.Uuid)
	require.NoError(t, err)

	_, err = repo.GetByUUID(bookmark.Uuid)
	require.ErrorIs(t, err, core.ErrNotFound)

	_, err = repo.GetByValue(bookmark.Value)
	require.ErrorIs(t, err, core.ErrNotFound)
}

func TestDelete_NotFound(t *testing.T) {
	bookmark := makeBookmark()
	repo := makePrepareRepository(bookmark)

	uuid7, _ := uuid.NewV7()
	err := repo.Delete(uuid7)
	require.ErrorIs(t, err, core.ErrNotFound)
}

func makeBookmark() model.Bookmark {
	bookmark, _ := model.NewBookmark(gofakeit.Word(), gofakeit.Animal())

	return bookmark
}

func makePrepareRepository(bookmark model.Bookmark) *repository {
	storage := memory.NewBookmarkStorage()
	repo := NewRepository(storage)

	_, _ = repo.Create(bookmark)

	return repo
}
