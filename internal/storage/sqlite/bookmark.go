package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"

	"bookmarks/internal/repository"
	"bookmarks/internal/storage"
	"bookmarks/pkg/sqlite"
)

type sqli struct {
	db *sql.DB
}

func NewStorage(sqlite *sqlite.Sqlite) (*sqli, error) {
	const op = "storage.sqlite.New"

	err := sqlite.Migrate()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &sqli{db: sqlite.Instance}, nil
}

func (s *sqli) Create(
	uuid uuid.UUID,
	title, val string,
	time time.Time,
) (storage.Bookmark, error) {
	const op = "storage.bookmark.Create"

	stmt, err := s.db.Prepare(`
		INSERT INTO bookmark(uuid, title, value, created_at)
		VALUES(?, ?, ?, ?)
		`)
	if err != nil {
		return storage.Bookmark{}, fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(uuid.String(), title, val, time)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			return storage.Bookmark{}, fmt.Errorf("%s: %w", op, repository.ErrExists)
		}

		return storage.Bookmark{}, fmt.Errorf("%s: %w", op, err)
	}

	return storage.Bookmark{
		Uuid:      uuid.String(),
		Title:     title,
		Value:     val,
		CreatedAt: time,
	}, nil
}

func (s *sqli) Update(
	uuid uuid.UUID,
	title string,
) (storage.Bookmark, error) {
	return storage.Bookmark{}, nil
}

func (s *sqli) GetByUUID(uuid uuid.UUID) (storage.Bookmark, error) {
	const op = "storage.bookmark.GetByUUID"

	stmt, err := s.db.Prepare("SELECT uuid, title, value, created_at FROM bookmark WHERE uuid = ?")
	if err != nil {
		return storage.Bookmark{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	defer stmt.Close()

	var record storage.Bookmark

	err = stmt.QueryRow(uuid).Scan(
		&record.Uuid,
		&record.Title,
		&record.Value,
		&record.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Bookmark{}, repository.ErrNotFound
		}

		return storage.Bookmark{}, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return record, nil
}

func (s *sqli) GetByValue(val string) (storage.Bookmark, error) {
	const op = "storage.bookmark.GetByValue"

	stmt, err := s.db.Prepare("SELECT uuid, title, value, created_at FROM bookmark WHERE value = ?")
	if err != nil {
		return storage.Bookmark{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	defer stmt.Close()

	var record storage.Bookmark

	err = stmt.QueryRow(val).Scan(
		&record.Uuid,
		&record.Title,
		&record.Value,
		&record.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Bookmark{}, repository.ErrNotFound
		}

		return storage.Bookmark{}, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return record, nil
}

func (s *sqli) Delete(uuid uuid.UUID) error {
	const op = "storage.bookmark.Delete"

	stmt, err := s.db.Prepare(`DELETE FROM bookmark WHERE uuid=?`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(uuid.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.ErrNotFound
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
