package pgsql

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"bookmarks/pkg/postgres"
)

type Pgsql struct {
	pool *pgxpool.Pool
}

func NewBookmark(p *postgres.Pgsql) (*Pgsql, error) {
	return &Pgsql{pool: p.Pool}, nil
}
