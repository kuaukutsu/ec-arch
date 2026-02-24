package bookmark

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"bookmarks/internal/model"
	"bookmarks/internal/repository"
)

var (
	ErrBookmarkExists   = errors.New("bookmark already exists")
	ErrBookmarkNotFound = errors.New("bookmark not found")
)

type Repository interface {
	Create(bookmark model.Bookmark) (model.Bookmark, error)
	Update(bookmark model.Bookmark) (model.Bookmark, error)
	GetByUUID(uuid uuid.UUID) (model.Bookmark, error)
	GetByValue(val string) (model.Bookmark, error)
	Delete(uuid uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) Append(title, val string) (model.Bookmark, error) {
	const op = "service.bookmark.Append"

	// exists
	bookmark, err := s.repo.GetByValue(val)
	if err == nil {
		return model.Bookmark{}, fmt.Errorf("%s: %w", op, ErrBookmarkExists)
	}

	if errors.Is(err, repository.ErrNotFound) {
		bookmark, err = model.NewBookmark(title, val)
		if err != nil {
			return model.Bookmark{}, fmt.Errorf("%s: %w", op, err)
		}

		bookmark, err = s.repo.Create(bookmark)
	}

	if err != nil {
		return model.Bookmark{}, fmt.Errorf("%s: %w", op, err)
	}

	return bookmark, nil
}

func (s *service) View(u string) (model.Bookmark, error) {
	const op = "service.bookmark.View"

	uuid, err := uuid.Parse(u)
	if err != nil {
		return model.Bookmark{}, fmt.Errorf("%s: %w", op, err)
	}

	bookmark, err := s.repo.GetByUUID(uuid)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.Bookmark{}, ErrBookmarkNotFound
		}

		return model.Bookmark{}, fmt.Errorf("%s: %w", op, err)
	}

	return bookmark, nil
}

func (s *service) Change() {
}

func (s *service) Delete(u string) error {
	const op = "service.bookmark.Delete"

	uuid, err := uuid.Parse(u)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.repo.Delete(uuid); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrBookmarkNotFound
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
