//go:build integration

package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go-server/internal/testutil"
	"go-server/pkg/database"
)

func TestTransactionReuseRollbackCancellationAndCommitFailure(t *testing.T) {
	ctx := context.Background()
	pg := testutil.StartPostgres(t, ctx)
	t.Cleanup(func() { _ = pg.Terminate(context.Background()) })
	pool, err := pgxpool.New(ctx, pg.DSN)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	d := &database.Database{Pool: pool}
	manager := NewTransactionManager(d)
	_, err = pool.Exec(ctx, `CREATE TABLE tx_parent(id bigint PRIMARY KEY);
 CREATE TABLE tx_child(id bigint PRIMARY KEY, parent_id bigint REFERENCES tx_parent(id) DEFERRABLE INITIALLY DEFERRED)`)
	if err != nil {
		t.Fatal(err)
	}
	sentinel := errors.New("abort after nested write")
	err = manager.WithTransaction(ctx, func(outer context.Context) error {
		first := getTxFromContext(outer)
		if _, err := GetExecutor(outer, d).Exec(outer, "INSERT INTO tx_parent VALUES(1)"); err != nil {
			return err
		}
		if err := manager.WithTransaction(outer, func(inner context.Context) error {
			if getTxFromContext(inner) != first {
				t.Fatal("nested scope started a different transaction")
			}
			_, err := GetExecutor(inner, d).Exec(inner, "INSERT INTO tx_child VALUES(1,1)")
			return err
		}); err != nil {
			return err
		}
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatal(err)
	}
	var count int
	if err := pool.QueryRow(ctx, "SELECT (SELECT count(*) FROM tx_parent)+(SELECT count(*) FROM tx_child)").Scan(&count); err != nil || count != 0 {
		t.Fatalf("rollback count=%d err=%v", count, err)
	}
	// This write succeeds; the deferred constraint fails only at COMMIT.
	err = manager.WithTransaction(ctx, func(tx context.Context) error {
		_, err := GetExecutor(tx, d).Exec(tx, "INSERT INTO tx_child VALUES(2,999)")
		return err
	})
	if err == nil {
		t.Fatal("commit failure was swallowed")
	}
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM tx_child").Scan(&count); err != nil || count != 0 {
		t.Fatalf("failed commit retained writes: %d %v", count, err)
	}
	timed, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	err = manager.WithTransaction(timed, func(tx context.Context) error {
		if deadline, ok := tx.Deadline(); !ok || deadline.IsZero() {
			t.Fatal("deadline was lost")
		}
		_, err := GetExecutor(tx, d).Exec(tx, "SELECT pg_sleep(10)")
		return err
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("deadline not propagated: %v", err)
	}
	// A fresh operation proves cancellation cleanup did not retain the connection/transaction.
	if err := manager.WithTransaction(ctx, func(tx context.Context) error {
		_, err := GetExecutor(tx, d).Exec(tx, "INSERT INTO tx_parent VALUES(3)")
		return err
	}); err != nil {
		t.Fatal(err)
	}
}
