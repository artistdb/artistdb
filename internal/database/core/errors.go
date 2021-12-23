package core

import "errors"

var (
	ErrNotFound    = errors.New("resource not found")
	ErrInvalidUUID = errors.New("id must be valid UUID")
)
