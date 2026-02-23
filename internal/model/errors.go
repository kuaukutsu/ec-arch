package model

import "errors"

var (
	ErrInvalidTitle = errors.New("invalid bookmark name")
	ErrInvalidValue = errors.New("invalid bookmark value")
)
