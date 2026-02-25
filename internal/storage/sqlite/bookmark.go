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

type Sqlite struct {
	db *sql.DB
}

func NewBookmark(sqlite *sqlite.Sqlite) (*Sqlite, error) {
	const op = "storage.sqlite.New"

	err := sqlite.Migrate()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Sqlite{db: sqlite.DB}, nil
}

func (s *Sqlite) Create(
	uuid uuid.UUID,
	title, val string,
	time time.Time,
) (storage.Bookmark, error) {
	const op = "storage.bookmark.Create"

	tx, err := s.db.Begin()
	if err != nil {
		return storage.Bookmark{}, fmt.Errorf("%s: %w", op, err)
	}

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

	if err := tx.Commit(); err != nil {
		return storage.Bookmark{}, fmt.Errorf("%s: %w", op, err)
	}

	return storage.Bookmark{
		Uuid:      uuid.String(),
		Title:     title,
		Value:     val,
		CreatedAt: time,
	}, nil
}

func (s *Sqlite) Update(
	uuid uuid.UUID,
	title string,
) (storage.Bookmark, error) {
	return storage.Bookmark{}, nil
}

func (s *Sqlite) GetByUUID(uuid uuid.UUID) (storage.Bookmark, error) {
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

func (s *Sqlite) GetByValue(val string) (storage.Bookmark, error) {
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

func (s *Sqlite) Delete(uuid uuid.UUID) error {
	const op = "storage.bookmark.Delete"

	stmt, err := s.db.Prepare(`DELETE FROM bookmark WHERE uuid=?`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	res, err := stmt.Exec(uuid.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.ErrNotFound
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}
