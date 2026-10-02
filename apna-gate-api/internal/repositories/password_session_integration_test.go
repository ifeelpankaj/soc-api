//go:build integration

package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"go-server/pkg/database"
	"sync"
	"testing"
)

func TestPasswordSessionAtomicUpdate(t *testing.T) {
	ctx := context.Background()
	container, err := postgres.Run(ctx, "postgres:16-alpine", postgres.WithDatabase("password_sessions"), postgres.WithUsername("postgres"), postgres.WithPassword("postgres"), postgres.BasicWaitStrategies())
	if err != nil {
		t.Fatalf("PostgreSQL test container unavailable: %v", err)
	}
	t.Cleanup(func() { _ = testcontainers.TerminateContainer(container) })
	connection, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	applyJobIntegrationMigrations(t, connection)
	pool, err := pgxpool.New(ctx, connection)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	database := &database.Database{Pool: pool}
	repo := NewUserRepository(database)
	tx := NewTransactionManager(database)
	id := insertReturningID(t, pool, `INSERT INTO users (full_name,email,password_hash) VALUES ('Password Test','password@example.com','original') RETURNING id`)
	user, err := repo.GetByID(ctx, id)
	if err != nil || user.SessionVersion != 0 {
		t.Fatalf("migration default: %v", err)
	}

	rollback := errors.New("force rollback")
	err = tx.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := repo.UpdatePasswordHash(txCtx, id, "rolled-back", 0); err != nil {
			return err
		}
		return rollback
	})
	if !errors.Is(err, rollback) {
		t.Fatalf("rollback error: %v", err)
	}
	user, err = repo.GetByID(ctx, id)
	if err != nil || user.SessionVersion != 0 || *user.PasswordHash != "original" {
		t.Fatalf("partial transaction persisted: %+v %v", user, err)
	}

	var wg sync.WaitGroup
	results := make(chan error, 2)
	for _, hash := range []string{"device-one", "device-two"} {
		wg.Add(1)
		go func(hash string) { defer wg.Done(); results <- repo.UpdatePasswordHash(ctx, id, hash, 0) }(hash)
	}
	wg.Wait()
	close(results)
	success := 0
	for err := range results {
		if err == nil {
			success++
		}
	}
	user, err = repo.GetByID(ctx, id)
	if err != nil || success != 1 || user.SessionVersion != 1 || *user.PasswordHash == "original" {
		t.Fatalf("concurrent password update: successes=%d user=%+v err=%v", success, user, err)
	}
	if err := repo.UpdatePasswordHash(ctx, id, "next-password", 1); err != nil {
		t.Fatal(err)
	}
	user, err = repo.GetByID(ctx, id)
	if err != nil || user.SessionVersion != 2 {
		t.Fatalf("second revocation failed: %+v %v", user, err)
	}
}
