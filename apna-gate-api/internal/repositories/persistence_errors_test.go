package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go-server/internal/repositories/contracts"
	"testing"
)

func TestPersistenceErrorClassification(t *testing.T) {
	for _, tc := range []struct{ input, want error }{
		{nil, nil}, {pgx.ErrNoRows, contracts.ErrNotFound},
		{&pgconn.PgError{Code: "23505"}, contracts.ErrAlreadyExists},
		{&pgconn.PgError{Code: "23503"}, contracts.ErrConflict},
		{&pgconn.PgError{Code: "23514"}, contracts.ErrConflict},
		{&pgconn.PgError{Code: "40001"}, contracts.ErrConcurrentModification},
		{&pgconn.PgError{Code: "40P01"}, contracts.ErrConcurrentModification},
		{context.Canceled, context.Canceled}, {context.DeadlineExceeded, context.DeadlineExceeded},
	} {
		input := tc.input
		if input != nil {
			input = fmt.Errorf("operation: %w", input)
		}
		if got := persistenceError(input); !errors.Is(got, tc.want) {
			t.Fatalf("%v: got %v, want %v", input, got, tc.want)
		}
	}
}
