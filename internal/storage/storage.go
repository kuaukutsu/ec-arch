package storage

import (
	"time"
)

type Bookmark struct {
	Uuid      string
	Title     string
	Value     string
	CreatedAt time.Time
}
