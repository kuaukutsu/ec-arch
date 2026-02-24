package memory

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"bookmarks/internal/repository"
	"bookmarks/internal/storage"
)

type db struct {
	mu    sync.RWMutex
	table map[uuid.UUID]*storage.Bookmark
	uiVal map[string]*storage.Bookmark // unique index by value
}

func NewBookmarkStorage() *db {
	return &db{
		table: make(map[uuid.UUID]*storage.Bookmark),
		uiVal: make(map[string]*storage.Bookmark),
	}
}

func (db *db) Create(
	uuid uuid.UUID,
	title string,
	val string,
	time time.Time,
) (storage.Bookmark, error) {
	const op = "storage.bookmark.Create"

	record := storage.Bookmark{
		Uuid:      uuid.String(),
		Title:     title,
		Value:     val,
		CreatedAt: time,
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.table[uuid]; exists {
		return storage.Bookmark{}, fmt.Errorf("%s: %w", op, repository.ErrExists)
	}

	if _, exists := db.uiVal[val]; exists {
		return storage.Bookmark{}, fmt.Errorf("%s: %w", op, repository.ErrExists)
	}

	db.table[uuid] = &record
	db.uiVal[val] = &record

	return record, nil
}

func (db *db) Update(uuid uuid.UUID, title string) (storage.Bookmark, error) {
	return storage.Bookmark{}, nil
}

func (db *db) GetByUUID(uuid uuid.UUID) (storage.Bookmark, error) {
	const op = "storage.bookmark.GetByUUID"

	db.mu.RLock()
	defer db.mu.RUnlock()

	record, exists := db.table[uuid]
	if !exists {
		return storage.Bookmark{}, fmt.Errorf("%s: %w", op, repository.ErrNotFound)
	}

	return *record, nil
}

func (db *db) GetByValue(val string) (storage.Bookmark, error) {
	const op = "storage.storage.GetByValue"

	db.mu.RLock()
	defer db.mu.RUnlock()

	record, exists := db.uiVal[val]
	if !exists {
		return storage.Bookmark{}, fmt.Errorf("%s: %w", op, repository.ErrNotFound)
	}

	return *record, nil
}

func (db *db) Delete(uuid uuid.UUID) error {
	const op = "storage.bookmark.Delete"

	db.mu.Lock()
	defer db.mu.Unlock()

	record, exists := db.table[uuid]
	if !exists {
		return fmt.Errorf("%s: %w", op, repository.ErrNotFound)
	}

	delete(db.table, uuid)
	delete(db.uiVal, record.Value)

	return nil
}
