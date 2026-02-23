package bookmark

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"bookmarks/internal/model"
	"bookmarks/internal/storage"
)

type Storage interface {
	Create(uuid uuid.UUID, title, val string, time time.Time) (storage.Bookmark, error)
	Update(uuid uuid.UUID, title string) (storage.Bookmark, error)
	GetByUUID(uuid uuid.UUID) (storage.Bookmark, error)
	GetByValue(val string) (storage.Bookmark, error)
	Delete(uuid uuid.UUID) error
}

type repository struct {
	storage Storage
}

func NewRepository(s Storage) *repository {
	return &repository{storage: s}
}

func (r *repository) Create(bookmark model.Bookmark) (model.Bookmark, error) {
	const op = "repository.bookmark.Create"

	_, err := r.storage.Create(
		bookmark.Uuid,
		bookmark.Title,
		bookmark.Value,
		bookmark.CreatedAt,
	)
	if err != nil {
		return model.Bookmark{}, fmt.Errorf("%s: %w", op, err)
	}

	return bookmark, nil
}

func (r *repository) Update(bookmark model.Bookmark) (model.Bookmark, error) {
	return bookmark, nil
}

func (r *repository) GetByUUID(uuid uuid.UUID) (model.Bookmark, error) {
	const op = "repository.bookmark.GetByUUID"

	record, err := r.storage.GetByUUID(uuid)
	if err != nil {
		return model.Bookmark{}, fmt.Errorf("%s: %w", op, err)
	}

	return castToModel(record)
}

func (r *repository) GetByValue(val string) (model.Bookmark, error) {
	const op = "repository.bookmark.GetByValue"

	record, err := r.storage.GetByValue(val)
	if err != nil {
		return model.Bookmark{}, fmt.Errorf("%s: %w", op, err)
	}

	return castToModel(record)
}

func (r *repository) Delete(uuid uuid.UUID) error {
	const op = "repository.bookmark.Delete"

	if err := r.storage.Delete(uuid); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func castToModel(r storage.Bookmark) (model.Bookmark, error) {
	const op = "repository.bookmark.castModel"

	uuid, err := uuid.Parse(r.Uuid)
	if err != nil {
		return model.Bookmark{}, fmt.Errorf("%s: %w", op, err)
	}

	return model.Bookmark{
		Uuid:      uuid,
		Title:     r.Title,
		Value:     r.Value,
		CreatedAt: r.CreatedAt,
	}, nil
}
