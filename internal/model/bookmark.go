package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Bookmark struct {
	Uuid      uuid.UUID
	Title     string
	Value     string // основное значение, которое нужно запомнить
	CreatedAt time.Time
}

func NewBookmark(title, value string) (Bookmark, error) {
	const op = "model.bookmark.New"

	title = strings.TrimSpace(title)
	if title == "" {
		return Bookmark{}, fmt.Errorf("%s: %w", op, ErrInvalidTitle)
	}

	value = strings.TrimSpace(value)
	if value == "" {
		return Bookmark{}, fmt.Errorf("%s: %w", op, ErrInvalidValue)
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		return Bookmark{}, fmt.Errorf("%s: %w", op, err)
	}

	return Bookmark{
		Uuid:      uuid,
		Title:     title,
		Value:     value,
		CreatedAt: time.Now(),
	}, nil
}
