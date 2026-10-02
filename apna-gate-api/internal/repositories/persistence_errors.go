package repository

import (
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go-server/internal/repositories/contracts"
)

// persistenceError preserves cancellation and unknown failures while translating
// database-specific outcomes into the repository contract.
func persistenceError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return contracts.ErrNotFound
	}
	var pg *pgconn.PgError
	if errors.As(err, &pg) {
		switch pg.Code {
		case "23505":
			return contracts.ErrAlreadyExists
		case "23503", "23514":
			return contracts.ErrConflict
		case "40001", "40P01":
			return contracts.ErrConcurrentModification
		}
	}
	return err
}
