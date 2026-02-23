package repository

import "errors"

var (
	ErrNotFound = errors.New("record not found")
	ErrExists   = errors.New("record exists")
)
