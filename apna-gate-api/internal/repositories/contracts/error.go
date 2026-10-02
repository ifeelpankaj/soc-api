package contracts

import "errors"

// Persistence errors describe outcomes without exposing the database driver.
var (
	ErrNotFound               = errors.New("repository: not found")
	ErrAlreadyExists          = errors.New("repository: already exists")
	ErrConflict               = errors.New("repository: conflict")
	ErrConcurrentModification = errors.New("repository: concurrent modification")
)
