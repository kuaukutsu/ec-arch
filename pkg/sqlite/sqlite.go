package sqlite

import (
	"database/sql"
	"fmt"
)

type Sqlite struct {
	driverName, dataSourceName string
	Instance                   *sql.DB
}

func New(options ...Option) (*Sqlite, error) {
	sqlite := &Sqlite{
		driverName: "sqlite3",
	}

	for _, opt := range options {
		opt(sqlite)
	}

	db, err := sql.Open(sqlite.driverName, sqlite.dataSourceName)
	if err != nil {
		return nil, err
	}

	sqlite.Instance = db

	return sqlite, nil
}

func (s *Sqlite) Migrate() error {
	const op = "sqlite.Migrate"

	for _, query := range tableSlice() {
		stmt, err := s.Instance.Prepare(query)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		defer stmt.Close()

		_, err = stmt.Exec()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func tableSlice() []string {
	return []string{
		`
		CREATE TABLE IF NOT EXISTS bookmark(
			uuid TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			value TEXT NOT NULL,
			created_at DATETIME NOT NULL);
		CREATE UNIQUE INDEX IF NOT EXISTS ui_value ON bookmark(value);
		`,
	}
}
