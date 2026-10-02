//go:build integration

package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	DSN       string
	Terminate func(context.Context) error
}

func StartPostgres(t *testing.T, ctx context.Context) *PostgresContainer {
	t.Helper()

	container, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("apna_gate_test"),
		postgres.WithUsername("apna_gate"),
		postgres.WithPassword("apna_gate"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("5432/tcp").WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("postgres connection string: %v", err)
	}

	return &PostgresContainer{
		DSN: dsn,
		Terminate: func(ctx context.Context) error {
			return container.Terminate(ctx)
		},
	}
}
